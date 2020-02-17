package handlers

import (
	"context"
	"dp-frontend-cookie-controller/mapper"
	"encoding/json"
	"errors"
	"github.com/ONSdigital/dp-cookies/cookies"
	"github.com/ONSdigital/log.go/log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Cookies that will not be removed deleted
var protectedCookies = [5]string{"access_token", "lang", "collection", "timeseriesbasket", "rememberBasket"}

// ClientError is an interface that can be used to retrieve the status code if a client has errored
type ClientError interface {
	Error() string
	Code() int
}

// RenderClient is an interface with methods for require for rendering a template
type RenderClient interface {
	Do(string, []byte) ([]byte, error)
}

// setStatusCode sets the status code of a http response to a relevant error code
func setStatusCode(req *http.Request, w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	if err, ok := err.(ClientError); ok {
		if err.Code() == http.StatusNotFound {
			status = err.Code()
		}
	}
	log.Event(req.Context(), "setting-response-status", log.Error(err))
	w.WriteHeader(status)
}

// getCookiePreferencePage talks to the renderer to get the cookie preference page
func getCookiePreferencePage(ctx context.Context, w http.ResponseWriter, req *http.Request, rendC RenderClient, cp cookies.Policy) error {
	var err error
	m := mapper.CreateCookieSettingPage(cp)
	b, err := json.Marshal(m)
	if err != nil {
		log.Event(ctx, "unable to marshal cookie preferences", log.Error(err))
		setStatusCode(req, w, err)
		return err
	}

	templateHTML, err := rendC.Do("cookies-preferences", b)
	if err != nil {
		log.Event(ctx, "getting template from renderer cookies-preferences failed", log.Error(err))
		setStatusCode(req, w, err)
		return err
	}
	if _, err := w.Write(templateHTML); err != nil {
		log.Event(ctx, "error on write of cookie template", log.Error(err))
		setStatusCode(req, w, err)
	}
	return err
}

// isProtectedCookie is a helper function that checks if a cookie is protected or not
func isProtectedCookie(stringToFind string, slice [5]string) bool {
	for _, element := range slice {
		if element == stringToFind {
			return true
		}
	}
	return false
}

// removeNonProtectedCookies sets cookies which replace any existing ones with an instant expiry; which removes them.
// This will not remove any protected cookies
func removeNonProtectedCookies(w http.ResponseWriter, req *http.Request) {
	for _, cookie := range req.Cookies() {
		if !isProtectedCookie(cookie.Name, protectedCookies) {
			cookie := &http.Cookie{
				Name:    cookie.Name,
				Value:   "",
				Path:    "/",
				Expires: time.Unix(0, 0),

				HttpOnly: false,
			}
			http.SetCookie(w, cookie)
		}
	}
}

// AcceptAll handler for setting all cookies to enabled then refresh the page. when JS has been disabled
// Example usage; JavaScript disabled.
func AcceptAll() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		acceptAll(w, req)
	}
}

// Read Handler
func Read(rendC RenderClient) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		read(w, req, rendC)
	}
}

// Edit Handler
func Edit(rendC RenderClient) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		edit(w, req, rendC)
	}
}

// acceptAll handler for accepting all possible cookies
func acceptAll(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	cp := cookies.Policy{
		Essential: true,
		Usage:     true,
	}
	reqUrl, err := url.Parse(req.URL.Path)
	if err != nil {
		log.Event(ctx, "unable to parse url", log.Error(err))
		setStatusCode(req, w, err)
		return
	}
	cookies.SetPolicy(w, cp, reqUrl.Hostname())
	cookies.SetPreferenceIsSet(w, reqUrl.Hostname())
	referer := req.Header.Get("Referer")
	if referer == "" {
		err := errors.New("cannot redirect due to no referer header")
		log.Event(ctx, "unable to parse url", log.Error(err))
		setStatusCode(req, w, err)
		return
	}

	http.Redirect(w, req, referer, http.StatusFound)
}

// edit handler for changing and setting cookie preferences, returns populated cookie preferences page from the renderer
func edit(w http.ResponseWriter, req *http.Request, rendC RenderClient) {
	ctx := req.Context()
	if err := req.ParseForm(); err != nil {
		log.Event(ctx, "failed to parse form input", log.Error(err))
		setStatusCode(req, w, err)
		return
	}

	essential, err := strconv.ParseBool(req.FormValue("essential"))
	if err != nil {
		log.Event(ctx, "failed to parse cookie value essential", log.Error(err))
		setStatusCode(req, w, err)
		return
	}
	usage, err := strconv.ParseBool(req.FormValue("usage"))
	if err != nil {
		log.Event(ctx, "failed to parse cookie value usage", log.Error(err))
		setStatusCode(req, w, err)
		return
	}
	cp := cookies.Policy{
		Essential: essential,
		Usage:     usage,
	}
	reqUrl, err := url.Parse(req.URL.Path)
	if err != nil {
		log.Event(ctx, "unable to parse url", log.Error(err))
	}
	if !usage {
		removeNonProtectedCookies(w, req)
	}
	cookies.SetPreferenceIsSet(w, reqUrl.Hostname())
	cookies.SetPolicy(w, cp, reqUrl.Hostname())

	err = getCookiePreferencePage(ctx, w, req, rendC, cp)
	if err != nil {
		log.Event(ctx, "getting cookie preference page failed", log.Error(err))
	}
	return

}

// read handler returns a populated cookie preferences page
func read(w http.ResponseWriter, req *http.Request, rendC RenderClient) {
	ctx := req.Context()
	cookiePref := cookies.GetCookiePreferences(req)

	err := getCookiePreferencePage(ctx, w, req, rendC, cookiePref.Policy)
	if err != nil {
		log.Event(ctx, "getting cookie preference page failed", log.Error(err))
	}
	return
}

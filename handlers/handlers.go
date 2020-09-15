package handlers

import (
	"dp-frontend-cookie-controller/mapper"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/ONSdigital/dp-cookies/cookies"
	dphandlers "github.com/ONSdigital/dp-net/handlers"
	"github.com/ONSdigital/log.go/log"
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
func getCookiePreferencePage(w http.ResponseWriter, req *http.Request, rendC RenderClient, cp cookies.Policy, isUpdated bool, lang string) error {
	var err error
	ctx := req.Context()
	m := mapper.CreateCookieSettingPage(cp, isUpdated, lang)
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
func isProtectedCookie(stringToFind string) bool {
	for _, element := range protectedCookies {
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
		if !isProtectedCookie(cookie.Name) {
			cookie := &http.Cookie{
				Name:     cookie.Name,
				Value:    "",
				Path:     "/",
				Expires:  time.Unix(0, 0),
				MaxAge:   0,
				HttpOnly: false,
			}
			http.SetCookie(w, cookie)
		}
	}
}

// Read Handler
func Read(rendC RenderClient) http.HandlerFunc {
	return dphandlers.ControllerHandler(func(w http.ResponseWriter, req *http.Request, lang, collectionID, accessToken string) {
		read(w, req, rendC, lang)
	})
}

// Edit Handler
func Edit(rendC RenderClient, siteDomain string) http.HandlerFunc {
	return dphandlers.ControllerHandler(func(w http.ResponseWriter, req *http.Request, lang, collectionID, accessToken string) {
		edit(w, req, rendC, siteDomain, lang)
	})
}

// edit handler for changing and setting cookie preferences, returns populated cookie preferences page from the renderer
func edit(w http.ResponseWriter, req *http.Request, rendC RenderClient, siteDomain, lang string) {
	ctx := req.Context()
	if err := req.ParseForm(); err != nil {
		log.Event(ctx, "failed to parse form input", log.Error(err))
		setStatusCode(req, w, err)
		return
	}
	cookiePolicyUsage := req.FormValue("cookie-policy-usage")
	if cookiePolicyUsage == "" {
		err := errors.New("request form value cookie-policy-usage not found")
		log.Event(ctx, "failed to get cookie value cookie-policy-usage from form", log.Error(err))
		setStatusCode(req, w, err)
		return
	}
	usage, err := strconv.ParseBool(cookiePolicyUsage)
	if err != nil {
		log.Event(ctx, "failed to parse cookie value usage", log.Error(err))
		setStatusCode(req, w, err)
		return
	}
	cp := cookies.Policy{
		Essential: true,
		Usage:     usage,
	}
	if !usage {
		removeNonProtectedCookies(w, req)
	}
	cookies.SetPreferenceIsSet(w, siteDomain)
	cookies.SetPolicy(w, cp, siteDomain)
	isUpdated := true
	err = getCookiePreferencePage(w, req, rendC, cp, isUpdated, lang)
	if err != nil {
		log.Event(ctx, "getting cookie preference page failed", log.Error(err))
	}
	return
}

// read handler returns a populated cookie preferences page
func read(w http.ResponseWriter, req *http.Request, rendC RenderClient, lang string) {
	ctx := req.Context()
	cookiePref := cookies.GetCookiePreferences(req)

	isUpdated := false
	err := getCookiePreferencePage(w, req, rendC, cookiePref.Policy, isUpdated, lang)
	if err != nil {
		log.Event(ctx, "getting cookie preference page failed", log.Error(err))
	}
	return
}

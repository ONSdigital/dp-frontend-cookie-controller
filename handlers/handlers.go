package handlers

import (
	"dp-frontend-cookie-controller/config"
	"dp-frontend-cookie-controller/mapper"
	"encoding/json"
	"fmt"
	"github.com/ONSdigital/dp-cookies/cookies"
	"github.com/ONSdigital/go-ns/common"
	"github.com/ONSdigital/log.go/log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var protectedCookies = [5]string{"access-token", "lang", "collection", "timeseriesbasket", "rememberBasket"}

// ClientError is an interface that can be used to retrieve the status code if a client has errored
type ClientError interface {
	Error() string
	Code() int
}

// RenderClient is an interface with methods for require for rendering a template
type RenderClient interface {
	Do(string, []byte) ([]byte, error)
}

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

// AcceptAll handler for setting all cookies to enabled then refresh the page. when JS has been disabled
// Example usage; JavaScript disabled.
func AcceptAll(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		acceptAll(w, req, cfg)
	}
}

// Edit Handler
func Read(rendC RenderClient) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		read(w, req, rendC)
	}
}

// Edit Handler
func Edit(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		edit(w, req, cfg)
	}
}

func acceptAll(w http.ResponseWriter, req *http.Request, cfg *config.Config) {
	ctx := req.Context()
	// TODO ensure XSS protection
	lang := common.GetLangFromSubDomain(req)
	uri := req.URL.Query().Get("redirect")
	uri, err := url.PathUnescape(uri)
	if err != nil {
		log.Event(ctx, "error fetching redirect", log.Error(err))
		setStatusCode(req, w, err)
		return
	}
	redirectURL := ""
	if lang != common.LangEN {
		redirectURL = fmt.Sprintf("https://%s.%s/%s", lang, cfg.SiteDomain, uri)
	} else {
		redirectURL = fmt.Sprintf("https://%s/%s", cfg.SiteDomain, uri)
	}

	cp := cookies.Policy{
		Essential: true,
		Usage:     true,
	}
	cookies.SetPolicy(w, cp, cfg.SiteDomain)
	cookies.SetPreferenceIsSet(w, cfg.SiteDomain)

	http.Redirect(w, req, redirectURL, 301)
}

func edit(w http.ResponseWriter, req *http.Request, cfg *config.Config) {
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
	cookies.SetPreferenceIsSet(w, cfg.SiteDomain)
	if !usage {
		removeNonProtectedCookies(w, req)
	}
	cookies.SetPolicy(w, cp, cfg.SiteDomain)
	return
}

func read(w http.ResponseWriter, req *http.Request, rendC RenderClient) {
	ctx := req.Context()
	cookiePref := cookies.GetCookiePreferences(req)

	m := mapper.CreateCookieSettingPage(cookiePref.Policy)

	b, err := json.Marshal(m)
	if err != nil {
		setStatusCode(req, w, err)
		return
	}

	templateHTML, err := rendC.Do("cookies-preferences", b)
	if err != nil {
		setStatusCode(req, w, err)
		return
	}

	if _, err := w.Write(templateHTML); err != nil {
		log.Event(ctx, "error on write of cookie template", log.Error(err))
		setStatusCode(req, w, err)
	}
	return
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

				HttpOnly: true,
			}
			http.SetCookie(w, cookie)
		}
	}
}

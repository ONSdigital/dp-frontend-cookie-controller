package handlers

import (
	"dp-frontend-cookie-controller/mapper"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/ONSdigital/dp-cookies/cookies"
	dphandlers "github.com/ONSdigital/dp-net/handlers"
	"github.com/ONSdigital/dp-renderer/model"
	"github.com/ONSdigital/log.go/v2/log"
)

// Cookies that will not be removed deleted
var protectedCookies = []string{"access_token", "refresh_token", "id_token", "lang", "collection", "timeseriesbasket", "rememberBasket"}

// RenderClient is an interface with required methods for building a template from a page model
type RenderClient interface {
	BuildPage(w io.Writer, pageModel interface{}, templateName string)
	NewBasePageModel() model.Page
}

// ClientError is an interface that can be used to retrieve the status code if a client has errored
type ClientError interface {
	Error() string
	Code() int
}

// setStatusCode sets the status code of a http response to a relevant error code
func setStatusCode(req *http.Request, w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	if err, ok := err.(ClientError); ok {
		if err.Code() == http.StatusNotFound {
			status = err.Code()
		}
	}
	log.Error(req.Context(), "setting-response-status", err)
	w.WriteHeader(status)
}

// getCookiePreferencePage talks to the renderer to get the cookie preference page
func getCookiePreferencePage(w http.ResponseWriter, rendC RenderClient, cp cookies.Policy, isUpdated bool, lang string) {
	basePage := rendC.NewBasePageModel()
	m := mapper.CreateCookieSettingPage(basePage, cp, isUpdated, lang)
	rendC.BuildPage(w, m, "cookies-preferences")
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
		log.Error(ctx, "failed to parse form input", err)
		setStatusCode(req, w, err)
		return
	}
	cookiePolicyUsage := req.FormValue("cookie-policy-usage")
	if cookiePolicyUsage == "" {
		err := errors.New("request form value cookie-policy-usage not found")
		log.Error(ctx, "failed to get cookie value cookie-policy-usage from form", err)
		setStatusCode(req, w, err)
		return
	}
	usage, err := strconv.ParseBool(cookiePolicyUsage)
	if err != nil {
		log.Error(ctx, "failed to parse cookie value usage", err)
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
	getCookiePreferencePage(w, rendC, cp, isUpdated, lang)
}

// read handler returns a populated cookie preferences page
func read(w http.ResponseWriter, req *http.Request, rendC RenderClient, lang string) {
	cookiePref := cookies.GetCookiePreferences(req)

	isUpdated := false
	getCookiePreferencePage(w, rendC, cookiePref.Policy, isUpdated, lang)
}

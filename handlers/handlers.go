package handlers

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/ONSdigital/dp-cookies/cookies"
	"github.com/ONSdigital/dp-frontend-cookie-controller/mapper"
	dphandlers "github.com/ONSdigital/dp-net/v3/handlers"
	"github.com/ONSdigital/dp-renderer/v2/model"
	"github.com/ONSdigital/log.go/v2/log"
)

// Cookies that will not be removed deleted
var protectedCookies = []string{"access_token", "refresh_token", "id_token", "lang", "collection", "timeseriesbasket", "rememberBasket"}

// To mock interfaces in this file
//go:generate mockgen -source=handlers.go -destination=mock_handlers.go -package=handlers github.com/ONSdigital/dp-frontend-cookie-controller/handlers RenderClient

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

func setStatusCode(req *http.Request, w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	if err, ok := err.(ClientError); ok {
		status = err.Code()
	}
	log.Error(req.Context(), "setting-response-status", err)
	w.WriteHeader(status)
}

// getCookiePreferencePage talks to the renderer to get the cookie preference page
func getCookiePreferencePage(w http.ResponseWriter, rendC RenderClient, cp cookies.ONSPolicy, isUpdated bool, lang string) {
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
			setCookie := &http.Cookie{
				Name:     cookie.Name,
				Value:    "",
				Path:     "/",
				Expires:  time.Unix(0, 0),
				MaxAge:   0,
				HttpOnly: false,
			}
			http.SetCookie(w, setCookie)
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

	// get form values
	cookiePolicyUsage := req.FormValue("cookie-policy-usage")
	if cookiePolicyUsage == "" {
		err := clientErr{errors.New("request form value cookie-policy-usage not found")}
		log.Info(ctx, "failed to get cookie value cookie-policy-usage from form", log.Data{"client_error": err})
		setStatusCode(req, w, err)
		return
	}
	cookiePolicyComms := req.FormValue("cookie-policy-comms")
	if cookiePolicyComms == "" {
		err := clientErr{errors.New("request form value cookie-policy-comms not found")}
		log.Info(ctx, "failed to get cookie value cookie-policy-comms from form", log.Data{"client_error": err})
		setStatusCode(req, w, err)
		return
	}
	cookiePolicySiteSettings := req.FormValue("cookie-policy-site-settings")
	if cookiePolicySiteSettings == "" {
		err := clientErr{errors.New("request form value cookie-policy-site-settings not found")}
		log.Info(ctx, "failed to get cookie value cookie-policy-site-settings from form", log.Data{"client_error": err})
		setStatusCode(req, w, err)
		return
	}

	// parse form values and make type safe
	usage, err := strconv.ParseBool(cookiePolicyUsage)
	if err != nil {
		err := clientErr{errors.New("request form value cookie-policy-usage not valid")}
		log.Info(ctx, "failed to parse cookie value usage", log.Data{"client_error": err})
		setStatusCode(req, w, err)
		return
	}
	comms, err := strconv.ParseBool(cookiePolicyComms)
	if err != nil {
		err := clientErr{errors.New("request form value cookie-policy-comms not valid")}
		log.Info(ctx, "failed to parse cookie value comms", log.Data{"client_error": err})
		setStatusCode(req, w, err)
		return
	}
	siteSettings, err := strconv.ParseBool(cookiePolicySiteSettings)
	if err != nil {
		err := clientErr{errors.New("request form value cookie-policy-site-usage not valid")}
		log.Info(ctx, "failed to parse cookie value site settings", log.Data{"client_error": err})
		setStatusCode(req, w, err)
		return
	}

	cp := cookies.ONSPolicy{
		Campaigns: comms,
		Essential: true, // always set to true
		Settings:  siteSettings,
		Usage:     usage,
	}

	// always remove non-protected cookies
	removeNonProtectedCookies(w, req)
	cookies.SetONSPreferenceIsSet(w, siteDomain)
	cookies.SetONSPolicy(w, cp, siteDomain)
	isUpdated := true
	getCookiePreferencePage(w, rendC, cp, isUpdated, lang)
}

// read handler returns a populated cookie preferences page
func read(w http.ResponseWriter, req *http.Request, rendC RenderClient, lang string) {
	cookiePref := cookies.GetONSCookiePreferences(req)

	isUpdated := false
	getCookiePreferencePage(w, rendC, cookiePref.Policy, isUpdated, lang)
}

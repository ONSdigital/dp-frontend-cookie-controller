package handlers

import (
	"context"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/ONSdigital/dis-design-system-go/model"
	"github.com/ONSdigital/dp-cookies/cookies"
	"github.com/ONSdigital/dp-frontend-cookie-controller/mapper"
	dphandlers "github.com/ONSdigital/dp-net/v3/handlers"
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
	if clientErr, ok := err.(ClientError); ok {
		status = clientErr.Code()
		log.Info(req.Context(), "setting client error response status")
	} else {
		log.Error(req.Context(), "setting internal error response status", err)
	}
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
func Edit(rendC RenderClient) http.HandlerFunc {
	return dphandlers.ControllerHandler(func(w http.ResponseWriter, req *http.Request, lang, collectionID, accessToken string) {
		edit(w, req, rendC, lang)
	})
}

// edit handler for changing and setting cookie preferences, returns populated cookie preferences page from the renderer
func edit(w http.ResponseWriter, req *http.Request, rendC RenderClient, lang string) {
	ctx := req.Context()
	if err := req.ParseForm(); err != nil {
		log.Error(ctx, "failed to parse form input", err)
		setStatusCode(req, w, err)
		return
	}

	// get and parse form values
	usage, err := getParsedBool(ctx, req, "cookie-policy-usage")
	if err != nil {
		setStatusCode(req, w, err)
		return
	}
	comms, err := getParsedBool(ctx, req, "cookie-policy-comms")
	if err != nil {
		setStatusCode(req, w, err)
		return
	}
	siteSettings, err := getParsedBool(ctx, req, "cookie-policy-site-settings")
	if err != nil {
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
	domain := req.Header.Get("X-Forwarded-Host")
	cookies.SetONSPreferenceIsSet(w, domain)
	cookies.SetONSPolicy(w, cp, domain)
	isUpdated := true
	getCookiePreferencePage(w, rendC, cp, isUpdated, lang)
}

// read handler returns a populated cookie preferences page
func read(w http.ResponseWriter, req *http.Request, rendC RenderClient, lang string) {
	cookiePref := cookies.GetONSCookiePreferences(req)

	isUpdated := false
	getCookiePreferencePage(w, rendC, cookiePref.Policy, isUpdated, lang)
}

// getFormValue is a helper function that retrieves the value of a form field from the request or returns an error
func getFormValue(ctx context.Context, req *http.Request, key string) (string, error) {
	value := req.FormValue(key)
	if value == "" {
		log.Info(ctx, "failed to get form value", log.Data{"key": key})
		return "", &clientErr{}
	}
	return value, nil
}

// parseFormValue is a helper function that parses a form value into a type safe boolean or returns an error
func parseFormValue(ctx context.Context, value string) (bool, error) {
	parsedValue, err := strconv.ParseBool(value)
	if err != nil {
		log.Info(ctx, "failed to parse form value", log.Data{"value": value})
		return false, &clientErr{}
	}

	return parsedValue, nil
}

// getParsedBool is a helper function to retrieve and parse form values
func getParsedBool(ctx context.Context, req *http.Request, key string) (bool, error) {
	value, err := getFormValue(ctx, req, key)
	if err != nil {
		return false, err
	}
	return parseFormValue(ctx, value)
}

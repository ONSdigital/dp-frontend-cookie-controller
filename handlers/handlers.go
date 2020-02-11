package handlers

import (
	"dp-frontend-cookie-controller/mapper"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/ONSdigital/log.go/log"
	"net/http"
	"net/url"
)

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
func AcceptAll() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		acceptAll(w, req)
	}
}

// Edit Handler
func Read(rendC RenderClient) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		read(w, req, rendC)
	}
}

// Edit Handler
func Edit() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		edit(w, req)
	}
}

func acceptAll(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	redirectURL := url.QueryEscape(req.URL.Query().Get("redirect"))
	if redirectURL == "" {
		err := errors.New("missing redirect URL")
		log.Event(ctx, "setting-response-status", log.Error(err))
		setStatusCode(req, w, err)
		return
	}

	http.Redirect(w, req, redirectURL, 301)
}

func edit(w http.ResponseWriter, req *http.Request) {
	//ctx := req.Context()

	return
}

func read(w http.ResponseWriter, req *http.Request, rendC RenderClient) {
	ctx := req.Context()

	cookie, err := req.Cookie("cookies_policy")
	if err != nil {
		// break out here, router should be defaulting the cookies if not set
		log.Event(ctx, "no cookie_policy found on request")
		return
	}

	// Use base64 to avoid issues with any illegal characters
	data, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		log.Event(ctx, "error decoding cookie", log.Error(err))
	}

	m := mapper.CreateCookieSettingPage(ctx, data)

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
	}
	return
}

package handlers

import (
	"dp-frontend-cookie-controller/mapper"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/ONSdigital/log.go/log"
	"net/http"
)

// ClientError is an interface that can be used to retrieve the status code if a client has errored
type ClientError interface {
	Error() string
	Code() int
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

func AcceptAll() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		acceptAll(w, req)
	}
}

// Edit Handler
func Read() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		read(w, req)
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
	greetingsModel := mapper.HelloModel{Greeting: "Hello", Who: "World"}
	m := mapper.HelloWorld(ctx, greetingsModel)

	b, err := json.Marshal(m)
	if err != nil {
		setStatusCode(req, w, err)
		return
	}

	_, err = w.Write(b)
	if err != nil {
		log.Event(ctx, "failed to write bytes for http response", log.Error(err))
		setStatusCode(req, w, err)
		return
	}
	return
}


func edit(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	greetingsModel := mapper.HelloModel{Greeting: "Hello", Who: "World"}
	m := mapper.HelloWorld(ctx, greetingsModel)

	b, err := json.Marshal(m)
	if err != nil {
		setStatusCode(req, w, err)
		return
	}

	_, err = w.Write(b)
    	if err != nil {
    		log.Event(ctx, "failed to write bytes for http response", log.Error(err))
    		setStatusCode(req, w, err)
    		return
    	}
    	return
}

func read(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	cookie, err := req.Cookie("cookies_policy")
	if err != nil {
		// break out here, router should be defaulting the cookies if not set
		log.Event(ctx, "no cookie_policy found on request")
		return
	}

	fmt.Printf("%s=%s\r\n", cookie.Name, cookie.Value)

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

	_, err = w.Write(b)
	if err != nil {
		log.Event(ctx, "failed to write bytes for http response", log.Error(err))
		setStatusCode(req, w, err)
		return
	}
	return
}

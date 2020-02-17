package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/ONSdigital/dp-cookies/cookies"
	"github.com/ONSdigital/log.go/log"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

type testCliError struct{}

func (e *testCliError) Error() string { return "client error" }
func (e *testCliError) Code() int     { return http.StatusNotFound }

// doTestRequest helper function that creates a router and mocks requests
func doTestRequest(target string, req *http.Request, handlerFunc http.HandlerFunc, w *httptest.ResponseRecorder) *httptest.ResponseRecorder {
	if w == nil {
		w = httptest.NewRecorder()
	}
	router := mux.NewRouter()
	router.HandleFunc(target, handlerFunc)
	router.ServeHTTP(w, req)
	return w
}

// cookiePolicyTest helper function that compares cookies on a httptest.ResponseRecorder with a given cookies.Policy
func cookiePolicyTest(w *httptest.ResponseRecorder, correctPolicy cookies.Policy) {
	allCookies := w.Result().Cookies()

	for _, c := range allCookies {
		if c.Name == "cookies_preferences_set" {
			So(c.Value, ShouldEqual, "true")
		}
		if c.Name == "cookies_policy" {
			cookiesPolicyUnescaped, err := url.QueryUnescape(c.Value)
			if err != nil {
				log.Event(nil, "unable to parse cookie", log.Error(err))
				return
			}
			var cpp cookies.Policy
			s, _ := strconv.Unquote(cookiesPolicyUnescaped)
			err = json.Unmarshal([]byte(s), &cpp)
			if err != nil {
				log.Event(nil, "unable to parse cookie", log.Error(err))
				return
			}
			So(cpp, ShouldResemble, correctPolicy)
		}
	}
}

func protectedCookiesTest(w *httptest.ResponseRecorder) bool {
	allCookies := w.Result().Cookies()
	allProtectedCookiesSafe := true

	// create map of protected cookies
	protectedCookieMap := make(map[string]bool)
	for i := 0; i < len(protectedCookies); i++ {
		protectedCookieMap[protectedCookies[i]] = false
	}

	// Check existing cookies for any with protected names, ensure they are still present
	for _, c := range allCookies {
		for key := range protectedCookieMap {
			if c.Name == key {
				protectedCookieMap[key] = true
			}
		}
	}
	for _, value := range protectedCookieMap {
		if value == false {
			allProtectedCookiesSafe = false
		}
	}

	return allProtectedCookiesSafe
}

// TestUnitHandlers unit tests for all handlers
func TestUnitHandlers(t *testing.T) {
	t.Parallel()
	Convey("test setStatusCode", t, func() {

		Convey("test status code handles 404 response from client", func() {
			req := httptest.NewRequest("GET", "http://localhost:23800", nil)
			w := httptest.NewRecorder()
			err := &testCliError{}
			setStatusCode(req, w, err)

			So(w.Code, ShouldEqual, http.StatusNotFound)
		})

		Convey("test status code handles internal server error", func() {
			req := httptest.NewRequest("GET", "/cookies/accept-all", nil)
			w := httptest.NewRecorder()
			err := errors.New("internal server error")
			setStatusCode(req, w, err)

			So(w.Code, ShouldEqual, http.StatusInternalServerError)
		})
	})

	Convey("test acceptAll", t, func() {

		Convey("is success", func() {
			cookiesPol := cookies.Policy{
				Essential: true,
				Usage:     true,
			}
			referer := "https://www.ons.gov.uk"
			req := httptest.NewRequest("GET", "/cookies/accept-all", nil)
			req.Header.Set("Referer", referer)
			w := doTestRequest("/cookies/accept-all", req, AcceptAll(), nil)

			So(w.Header().Get("Location"), ShouldEqual, referer)
			So(w.Code, ShouldEqual, http.StatusFound)
			cookiePolicyTest(w, cookiesPol)
			// TODO once library update check cookies have been set
		})

		Convey("is failure no referer header", func() {
			req := httptest.NewRequest("GET", "/cookies/accept-all", nil)
			w := doTestRequest("/cookies/accept-all", req, AcceptAll(), nil)

			So(w.Code, ShouldEqual, http.StatusInternalServerError)
		})
	})

	Convey("test read", t, func() {
		Convey("with no cookies set", func() {
			mockCtrl := gomock.NewController(t)
			mockRend := NewMockRenderClient(mockCtrl)
			mockRend.EXPECT().Do("cookies-preferences", gomock.Any()).Return([]byte(`<html><body><h1>Some HTML from renderer!</h1></body></html>`), nil)
			req := httptest.NewRequest("GET", "/cookies/edit", nil)
			w := doTestRequest("/cookies/edit", req, Read(mockRend), nil)
			So(w.Code, ShouldEqual, http.StatusOK)
			So(w.Body.String(), ShouldEqual, "<html><body><h1>Some HTML from renderer!</h1></body></html>")
		})

		Convey("with cookies already set", func() {
			cookiesPol := cookies.Policy{
				Essential: true,
				Usage:     false,
			}

			w := httptest.NewRecorder()
			cookies.SetPreferenceIsSet(w, "domain")
			cookies.SetPolicy(w, cookiesPol, "domain")

			mockCtrl := gomock.NewController(t)
			mockRend := NewMockRenderClient(mockCtrl)
			mockRend.EXPECT().Do("cookies-preferences", gomock.Any()).Return([]byte(`<html><body><h1>Some HTML from renderer!</h1></body></html>`), nil)
			req := httptest.NewRequest("GET", "/cookies/edit", nil)
			w = doTestRequest("/cookies/edit", req, Read(mockRend), w)
			So(w.Code, ShouldEqual, http.StatusOK)
			So(w.Body.String(), ShouldEqual, "<html><body><h1>Some HTML from renderer!</h1></body></html>")

		})

		Convey("with renderer failing", func() {
			mockCtrl := gomock.NewController(t)
			mockRend := NewMockRenderClient(mockCtrl)
			mockRend.EXPECT().Do("cookies-preferences", gomock.Any()).Return(nil, errors.New("error from renderer"))
			req := httptest.NewRequest("GET", "/cookies/edit", nil)
			w := doTestRequest("/cookies/edit", req, Read(mockRend), nil)
			So(w.Code, ShouldEqual, http.StatusInternalServerError)
		})
	})

	Convey("test edit", t, func() {
		cookieTimeSeriesBasket := &http.Cookie{
			Name:     "timeseriesbasket",
			Value:    url.QueryEscape("timeseriesbasketData"),
			Path:     "/",
			Domain:   "domain",
			HttpOnly: false,
			SameSite: http.SameSiteLaxMode,
		}
		cookieRememberBasket := &http.Cookie{
			Name:     "rememberBasket",
			Value:    url.QueryEscape("rememberBasketData"),
			Path:     "/",
			Domain:   "domain",
			HttpOnly: false,
			SameSite: http.SameSiteLaxMode,
		}

		mockCtrl := gomock.NewController(t)
		mockRend := NewMockRenderClient(mockCtrl)
		Convey("success with good form no prior cookies set", func() {
			cookiesPol := cookies.Policy{
				Essential: true,
				Usage:     true,
			}
			mockRend.EXPECT().Do("cookies-preferences", gomock.Any()).Return([]byte(`<html><body><h1>Some HTML from renderer!</h1></body></html>`), nil)
			b := `essential=true&usage=true`
			req := httptest.NewRequest("POST", "/cookies/edit", bytes.NewBufferString(b))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := doTestRequest("/cookies/edit", req, Edit(mockRend), nil)
			So(w.Code, ShouldEqual, http.StatusOK)
			So(w.Body.String(), ShouldEqual, "<html><body><h1>Some HTML from renderer!</h1></body></html>")
			cookiePolicyTest(w, cookiesPol)
		})
		Convey("success with good form and prior cookies set", func() {
			authToken := "token"
			collection := "collection"
			lang := "cy"
			mockRend.EXPECT().Do("cookies-preferences", gomock.Any()).Return([]byte(`<html><body><h1>Some HTML from renderer!</h1></body></html>`), nil)
			b := `essential=true&usage=false`
			req := httptest.NewRequest("POST", "/cookies/edit", bytes.NewBufferString(b))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			cookiesPol := cookies.Policy{
				Essential: true,
				Usage:     false,
			}

			w := httptest.NewRecorder()
			cookies.SetPreferenceIsSet(w, "domain")
			cookies.SetPolicy(w, cookiesPol, "domain")
			cookies.SetUserAuthToken(w, authToken, "domain")
			cookies.SetCollection(w, collection, "domain")
			cookies.SetLang(w, lang, "domain")
			http.SetCookie(w, cookieRememberBasket)
			http.SetCookie(w, cookieTimeSeriesBasket)
			w = doTestRequest("/cookies/edit", req, Edit(mockRend), w)
			So(w.Code, ShouldEqual, http.StatusOK)
			So(w.Body.String(), ShouldEqual, "<html><body><h1>Some HTML from renderer!</h1></body></html>")
			cookiePolicyTest(w, cookiesPol)
			allProtectedCookiesFound := protectedCookiesTest(w)
			So(allProtectedCookiesFound, ShouldEqual, true)
		})

		Convey("fail with bad form names", func() {
			mockRend.EXPECT().Do("cookies-preferences", gomock.Any()).Return([]byte(`<html><body><h1>Some HTML from renderer!</h1></body></html>`), nil)
			b := `waffles=true&usage=true`
			req := httptest.NewRequest("POST", "/cookies/edit", bytes.NewBufferString(b))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := doTestRequest("/cookies/edit", req, Edit(mockRend), nil)
			So(w.Code, ShouldEqual, http.StatusInternalServerError)
		})
		Convey("fail with bad form values", func() {
			mockRend.EXPECT().Do("cookies-preferences", gomock.Any()).Return([]byte(`<html><body><h1>Some HTML from renderer!</h1></body></html>`), nil)
			b := `essential=nonboolvalue&usage=true`
			req := httptest.NewRequest("POST", "/cookies/edit", bytes.NewBufferString(b))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := doTestRequest("/cookies/edit", req, Edit(mockRend), nil)
			So(w.Code, ShouldEqual, http.StatusInternalServerError)
		})
		Convey("fail with renderer error", func() {
			mockRend.EXPECT().Do("cookies-preferences", gomock.Any()).Return(nil, errors.New("error from renderer"))
			b := `essential=true&usage=true`
			req := httptest.NewRequest("POST", "/cookies/edit", bytes.NewBufferString(b))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := doTestRequest("/cookies/edit", req, Edit(mockRend), nil)
			So(w.Code, ShouldEqual, http.StatusInternalServerError)
		})
	})
}

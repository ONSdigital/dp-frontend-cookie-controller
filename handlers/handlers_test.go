package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	coreModel "github.com/ONSdigital/dis-design-system-go/model"
	"github.com/ONSdigital/dp-cookies/cookies"
	"github.com/ONSdigital/dp-frontend-cookie-controller/config"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

type testCliError struct{}

func (e *testCliError) Error() string { return "client error" }
func (e *testCliError) Code() int     { return http.StatusNotFound }

func TestReadHandler(t *testing.T) {
	cfg := initialiseMockConfig()
	Convey("test read", t, func() {
		Convey("with no cookies set", func() {
			mockCtrl := gomock.NewController(t)
			mockRend := NewMockRenderClient(mockCtrl)
			mockRend.EXPECT().NewBasePageModel().Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain))
			mockRend.EXPECT().BuildPage(gomock.Any(), gomock.Any(), gomock.Eq("cookies-preferences"))
			req := httptest.NewRequest("GET", "/cookies", http.NoBody)
			w := doTestRequest("/cookies", req, Read(mockRend), nil)
			So(w.Code, ShouldEqual, http.StatusOK)
		})

		Convey("with cookies already set", func() {
			cookiesSetPolicy := cookies.ONSPolicy{
				Essential: true,
				Usage:     true,
				Campaigns: false,
				Settings:  false,
			}

			w := httptest.NewRecorder()
			cookies.SetONSPreferenceIsSet(w, "domain")
			cookies.SetONSPolicy(w, cookiesSetPolicy, "domain")

			mockCtrl := gomock.NewController(t)
			mockRend := NewMockRenderClient(mockCtrl)
			mockRend.EXPECT().NewBasePageModel().Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain))
			mockRend.EXPECT().BuildPage(gomock.Any(), gomock.Any(), gomock.Eq("cookies-preferences"))

			req := httptest.NewRequest("GET", "/cookies", http.NoBody)
			w = doTestRequest("/cookies", req, Read(mockRend), w)
			So(w.Code, ShouldEqual, http.StatusOK)
		})
	})
}

func TestEditHandler(t *testing.T) {
	cfg := initialiseMockConfig()
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
			cookiesPol := cookies.ONSPolicy{
				Campaigns: true,
				Essential: true,
				Settings:  true,
				Usage:     true,
			}

			mockRend.EXPECT().NewBasePageModel().Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain))
			mockRend.EXPECT().BuildPage(gomock.Any(), gomock.Any(), gomock.Eq("cookies-preferences"))

			b := `cookie-policy-usage=true&cookie-policy-comms=true&cookie-policy-site-settings=true`
			req := httptest.NewRequest("POST", "/cookies", bytes.NewBufferString(b))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := doTestRequest("/cookies", req, Edit(mockRend), nil)
			So(w.Code, ShouldEqual, http.StatusOK)
			cookiePolicyTest(w, cookiesPol)
		})
		Convey("success with good form and prior cookies set", func() {
			essentialSetCookiesPolicy := cookies.ONSPolicy{
				Campaigns: false,
				Essential: true,
				Usage:     false,
				Settings:  false,
			}

			authToken := "token"
			refreshToken := "refresh"
			idToken := "id"
			collection := "collection"
			lang := "cy"

			mockRend.EXPECT().NewBasePageModel().Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain))
			mockRend.EXPECT().BuildPage(gomock.Any(), gomock.Any(), gomock.Eq("cookies-preferences"))

			b := `cookie-policy-usage=false&cookie-policy-comms=false&cookie-policy-site-settings=false`
			req := httptest.NewRequest("POST", "/cookies", bytes.NewBufferString(b))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			w := httptest.NewRecorder()

			cookies.SetONSPreferenceIsSet(w, "domain")
			cookies.SetONSPolicy(w, essentialSetCookiesPolicy, "domain")
			cookies.SetUserAuthToken(w, authToken, "domain")
			cookies.SetRefreshToken(w, refreshToken, "domain")
			cookies.SetIDToken(w, idToken, "domain")
			cookies.SetCollection(w, collection, "domain")
			cookies.SetLang(w, lang, "domain")

			http.SetCookie(w, cookieRememberBasket)
			http.SetCookie(w, cookieTimeSeriesBasket)

			w = doTestRequest("/cookies", req, Edit(mockRend), w)

			So(w.Code, ShouldEqual, http.StatusOK)
			cookiePolicyTest(w, essentialSetCookiesPolicy)
			allProtectedCookiesFound := protectedCookiesTest(w)
			So(allProtectedCookiesFound, ShouldEqual, true)
		})

		Convey("400 with bad form names", func() {
			b := `cookie-policy-waffles=true`
			req := httptest.NewRequest("POST", "/cookies", bytes.NewBufferString(b))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := doTestRequest("/cookies", req, Edit(mockRend), nil)
			So(w.Code, ShouldEqual, http.StatusBadRequest)
		})

		Convey("400 with omitted form values", func() {
			Convey("cookie-policy-usage", func() {
				b := `cookie-policy-usage=&cookie-policy-comms=false&cookie-policy-site-settings=false`
				req := httptest.NewRequest("POST", "/cookies", bytes.NewBufferString(b))
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				w := doTestRequest("/cookies", req, Edit(mockRend), nil)
				So(w.Code, ShouldEqual, http.StatusBadRequest)
			})
			Convey("cookie-policy-comms", func() {
				b := `cookie-policy-usage=false&cookie-policy-comms=`
				req := httptest.NewRequest("POST", "/cookies", bytes.NewBufferString(b))
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				w := doTestRequest("/cookies", req, Edit(mockRend), nil)
				So(w.Code, ShouldEqual, http.StatusBadRequest)
			})
			Convey("cookie-policy-settings", func() {
				b := `cookie-policy-usage=false&cookie-policy-comms=false&cookie-policy-site-settings=`
				req := httptest.NewRequest("POST", "/cookies", bytes.NewBufferString(b))
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				w := doTestRequest("/cookies", req, Edit(mockRend), nil)
				So(w.Code, ShouldEqual, http.StatusBadRequest)
			})
		})

		Convey("400 with bad form values", func() {
			Convey("cookie-policy-usage", func() {
				b := `cookie-policy-usage=nonbool&cookie-policy-comms=false&cookie-policy-site-settings=false`
				req := httptest.NewRequest("POST", "/cookies", bytes.NewBufferString(b))
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				w := doTestRequest("/cookies", req, Edit(mockRend), nil)
				So(w.Code, ShouldEqual, http.StatusBadRequest)
			})
			Convey("cookie-policy-comms", func() {
				b := `cookie-policy-usage=false&cookie-policy-comms=blah&cookie-policy-site-settings=false`
				req := httptest.NewRequest("POST", "/cookies", bytes.NewBufferString(b))
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				w := doTestRequest("/cookies", req, Edit(mockRend), nil)
				So(w.Code, ShouldEqual, http.StatusBadRequest)
			})
			Convey("cookie-policy-settings", func() {
				b := `cookie-policy-usage=false&cookie-policy-comms=false&cookie-policy-site-settings=notbool`
				req := httptest.NewRequest("POST", "/cookies", bytes.NewBufferString(b))
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				w := doTestRequest("/cookies", req, Edit(mockRend), nil)
				So(w.Code, ShouldEqual, http.StatusBadRequest)
			})
		})
	})
}

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
func cookiePolicyTest(w *httptest.ResponseRecorder, correctPolicy cookies.ONSPolicy) {
	allCookies := w.Result().Cookies()
	defer w.Result().Body.Close()

	for _, c := range allCookies {
		if c.Name == "ons_cookies_preferences_set" {
			So(c.Value, ShouldEqual, "true")
		}
		if c.Name == "ons_cookie_policy" {
			cookiesPolicyUnescaped, err := url.QueryUnescape(c.Value)
			if err != nil {
				log.Error(context.Background(), "unable to parse cookie", err)
				return
			}
			var cpp cookies.ONSPolicy
			s, _ := strconv.Unquote(cookiesPolicyUnescaped)
			err = json.Unmarshal([]byte(s), &cpp)
			if err != nil {
				log.Error(context.Background(), "unable to parse cookie", err)
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

func initialiseMockConfig() config.Config {
	return config.Config{
		PatternLibraryAssetsPath: "http://localhost:9002/dist",
		SiteDomain:               "ons",
		SupportedLanguages:       [2]string{"en", "cy"},
	}
}

// TestUnitHandlers unit tests for all handlers
func TestUnitHandlers(t *testing.T) {
	Convey("test setStatusCode", t, func() {
		Convey("test status code handles 404 response from client", func() {
			req := httptest.NewRequest("GET", "http://localhost:24100", http.NoBody)
			w := httptest.NewRecorder()
			err := &testCliError{}
			setStatusCode(req, w, err)

			So(w.Code, ShouldEqual, http.StatusNotFound)
		})
	})
}

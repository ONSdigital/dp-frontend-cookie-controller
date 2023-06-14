package handlers

import (
	"bytes"
	"context"
	"dp-frontend-cookie-controller/config"
	"dp-frontend-cookie-controller/model"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	coreModel "github.com/ONSdigital/dp-renderer/v2/model"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"

	"github.com/ONSdigital/dp-cookies/cookies"
	"github.com/golang/mock/gomock"
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
			mockRend.EXPECT().BuildPage(gomock.Any(), gomock.Eq(initialiseMockCookiesPageModel(&cfg, cookies.Policy{Essential: true}, false, false, "en")), gomock.Eq("cookies-preferences"))
			req := httptest.NewRequest("GET", "/cookies", nil)
			w := doTestRequest("/cookies", req, Read(mockRend), nil)
			So(w.Code, ShouldEqual, http.StatusOK)
		})

		Convey("with cookies already set", func() {
			cookiesSetPolicy := cookies.Policy{
				Essential: true,
				Usage:     false,
			}

			w := httptest.NewRecorder()
			cookies.SetPreferenceIsSet(w, "domain")
			cookies.SetPolicy(w, cookiesSetPolicy, "domain")

			mockCtrl := gomock.NewController(t)
			mockRend := NewMockRenderClient(mockCtrl)
			mockRend.EXPECT().NewBasePageModel().Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain))
			mockRend.EXPECT().BuildPage(gomock.Any(), gomock.Eq(initialiseMockCookiesPageModel(&cfg, cookiesSetPolicy, false, false, "en")), gomock.Eq("cookies-preferences")).Times(1)

			req := httptest.NewRequest("GET", "/cookies", nil)
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
			cookiesPol := cookies.Policy{
				Essential: true,
				Usage:     true,
			}

			mockRend.EXPECT().NewBasePageModel().Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain))
			mockRend.EXPECT().BuildPage(gomock.Any(), gomock.Any(), gomock.Eq("cookies-preferences"))

			b := `cookie-policy-usage=true`
			req := httptest.NewRequest("POST", "/cookies", bytes.NewBufferString(b))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := doTestRequest("/cookies", req, Edit(mockRend, cfg.SiteDomain), nil)
			So(w.Code, ShouldEqual, http.StatusOK)
			cookiePolicyTest(w, cookiesPol)
		})
		Convey("success with good form and prior cookies set", func() {
			essentialSetCookiesPolicy := cookies.Policy{
				Essential: true,
				Usage:     false,
			}

			authToken := "token"
			refreshToken := "refresh"
			idToken := "id"
			collection := "collection"
			lang := "cy"
			hasBeenUpdated := true
			cookiesPreferenceIsSet := true

			mockRend.EXPECT().NewBasePageModel().Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain))
			mockRend.EXPECT().BuildPage(gomock.Any(), initialiseMockCookiesPageModel(&cfg, essentialSetCookiesPolicy, hasBeenUpdated, cookiesPreferenceIsSet, "en"), gomock.Eq("cookies-preferences"))

			b := `cookie-policy-usage=false`
			req := httptest.NewRequest("POST", "/cookies", bytes.NewBufferString(b))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			w := httptest.NewRecorder()

			cookies.SetPreferenceIsSet(w, "domain")
			cookies.SetPolicy(w, essentialSetCookiesPolicy, "domain")
			cookies.SetUserAuthToken(w, authToken, "domain")
			cookies.SetRefreshToken(w, refreshToken, "domain")
			cookies.SetIDToken(w, idToken, "domain")
			cookies.SetCollection(w, collection, "domain")
			cookies.SetLang(w, lang, "domain")

			http.SetCookie(w, cookieRememberBasket)
			http.SetCookie(w, cookieTimeSeriesBasket)

			w = doTestRequest("/cookies", req, Edit(mockRend, cfg.SiteDomain), w)

			So(w.Code, ShouldEqual, http.StatusOK)
			cookiePolicyTest(w, essentialSetCookiesPolicy)
			allProtectedCookiesFound := protectedCookiesTest(w)
			So(allProtectedCookiesFound, ShouldEqual, true)
		})

		Convey("fail with bad form names", func() {
			mockRend.EXPECT().NewBasePageModel().Return(coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain))
			mockRend.EXPECT().BuildPage(gomock.Any(), gomock.Any(), gomock.Eq("cookies-preferences"))
			b := `cookie-policy-waffles=true`
			req := httptest.NewRequest("POST", "/cookies", bytes.NewBufferString(b))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := doTestRequest("/cookies", req, Edit(mockRend, cfg.SiteDomain), nil)
			So(w.Code, ShouldEqual, http.StatusInternalServerError)
		})

		Convey("fail with bad form values", func() {
			b := `cookie-policy-usage=nonbool`
			req := httptest.NewRequest("POST", "/cookies", bytes.NewBufferString(b))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := doTestRequest("/cookies", req, Edit(mockRend, cfg.SiteDomain), nil)
			So(w.Code, ShouldEqual, http.StatusInternalServerError)
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
func cookiePolicyTest(w *httptest.ResponseRecorder, correctPolicy cookies.Policy) {
	allCookies := w.Result().Cookies()
	defer w.Result().Body.Close()

	for _, c := range allCookies {
		if c.Name == "cookies_preferences_set" {
			So(c.Value, ShouldEqual, "true")
		}
		if c.Name == "cookies_policy" {
			cookiesPolicyUnescaped, err := url.QueryUnescape(c.Value)
			if err != nil {
				log.Error(context.Background(), "unable to parse cookie", err)
				return
			}
			var cpp cookies.Policy
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
			req := httptest.NewRequest("GET", "http://localhost:24100", nil)
			w := httptest.NewRecorder()
			err := &testCliError{}
			setStatusCode(req, w, err)

			So(w.Code, ShouldEqual, http.StatusNotFound)
		})
	})
}

func initialiseMockCookiesPageModel(cfg *config.Config, policy cookies.Policy, isUpdated, hasSetPreference bool, lang string) model.CookiesPreference {
	page := model.CookiesPreference{
		Page: coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain),
	}
	page.Breadcrumb = []coreModel.TaxonomyNode{
		{
			Title: "Home",
			URI:   "/",
		},
		{
			Title: "Cookies",
		},
	}
	page.Metadata.Title = "Cookies"
	page.Language = lang
	page.CookiesPreferencesSet = true
	page.CookiesPolicy.Essential = policy.Essential
	page.CookiesPolicy.Usage = policy.Usage
	page.FeatureFlags.HideCookieBanner = true
	page.SiteDomain = cfg.SiteDomain
	page.PatternLibraryAssetsPath = cfg.PatternLibraryAssetsPath
	page.PreferencesUpdated = hasSetPreference

	page.TypeRadios = coreModel.RadioFieldset{
		Radios: []coreModel.Radio{
			{
				Input: coreModel.Input{
					ID:        "usage-on",
					IsChecked: page.CookiesPolicy.Usage,
					Label: coreModel.Localisation{
						LocaleKey: "On",
						Plural:    1,
					},
					Name:  "cookie-policy-usage",
					Value: "true",
				},
			},
			{
				Input: coreModel.Input{
					ID:        "usage-off",
					IsChecked: !page.CookiesPolicy.Usage,
					Label: coreModel.Localisation{
						LocaleKey: "Off",
						Plural:    1,
					},
					Name:  "cookie-policy-usage",
					Value: "false",
				},
			},
		},
	}
	return page
}

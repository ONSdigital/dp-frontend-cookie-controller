package handlers

import (
	"errors"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testCliError struct{}

func (e *testCliError) Error() string { return "client error" }
func (e *testCliError) Code() int     { return http.StatusNotFound }

func TestUnitHandlers(t *testing.T) {
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

	Convey("test acceptAll ", t, func() {

		Convey("is success", func() {
			referer := "https://www.ons.gov.uk"
			req := httptest.NewRequest("GET", "/cookies/accept-all", nil)
			req.Header.Set("Referer", referer)
			w := doTestRequest("/cookies/accept-all", req, AcceptAll())

			So(w.Header().Get("Location"), ShouldEqual, referer)
			So(w.Code, ShouldEqual, http.StatusFound)
			// TODO once library update check cookies have been set
		})

		Convey("is failure no referer header", func() {
			req := httptest.NewRequest("GET", "/cookies/accept-all", nil)
			w := doTestRequest("/cookies/accept-all", req, AcceptAll())

			So(w.Code, ShouldEqual, http.StatusInternalServerError)
		})
	})
}

func doTestRequest(target string, req *http.Request, handlerFunc http.HandlerFunc) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc(target, handlerFunc)
	router.ServeHTTP(w, req)
	return w
}

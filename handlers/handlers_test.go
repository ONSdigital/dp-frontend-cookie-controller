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
			req := httptest.NewRequest("GET", "http://localhost:23800", nil)
			w := httptest.NewRecorder()
			err := errors.New("internal server error")

			setStatusCode(req, w, err)

			So(w.Code, ShouldEqual, http.StatusInternalServerError)
		})
	})

	Convey("test acceptAll", t, func() {
		//ref := "https://www.ons.gov.uk"
		//req := httptest.NewRequest("GET", "http://localhost:23800/cookies/accept-all", nil)
		//w := httptest.NewRecorder()
		//r := mux.NewRouter()
		//req.Header.Set("Referer", ref)
		//
		//r.HandleFunc(AcceptAll())
		//r.ServeHTTP(w, req)
		//So(w.Header().Get("Location"), ShouldEqual, ref)
	})
}


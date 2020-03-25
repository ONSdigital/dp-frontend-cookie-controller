package mapper

import (
	"fmt"
	"testing"

	"github.com/ONSdigital/dp-cookies/cookies"
	"github.com/ONSdigital/dp-frontend-models/model"
	"github.com/ONSdigital/dp-frontend-models/model/cookiespreferences"
	. "github.com/smartystreets/goconvey/convey"
)

// TestUnitMapper tests mapper functions
func TestUnitMapper(t *testing.T) {
	t.Parallel()
	cookiesPolicy := cookies.Policy{
		Essential: true,
		Usage:     false,
	}
	expectedModel := cookiespreferences.Page{}
	expectedModel.Breadcrumb = []model.TaxonomyNode{
		{
			Title: "Home",
			URI:   "/",
		},
		{
			Title: "Cookies",
		},
	}
	expectedModel.Metadata.Title = "Cookies"
	expectedModel.CookiesPreferencesSet = true
	expectedModel.CookiesPolicy.Essential = true
	expectedModel.CookiesPolicy.Usage = false
	expectedModel.PreferencesUpdated = false
	Convey("test CreateCookieSettingPage", t, func() {
		mcp := CreateCookieSettingPage(cookiesPolicy, false)
		fmt.Printf("%+v\n", mcp)
		So(expectedModel, ShouldResemble, mcp)
	})
}

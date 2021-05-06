package mapper

import (
	"dp-frontend-cookie-controller/config"
	"dp-frontend-cookie-controller/models"
	"fmt"
	"testing"

	coreModel "github.com/rav-pradhan/test-modules/render/models"

	"github.com/ONSdigital/dp-cookies/cookies"
	request "github.com/ONSdigital/dp-net/request"
	. "github.com/smartystreets/goconvey/convey"
)

// TestUnitMapper tests mapper functions
func TestUnitMapper(t *testing.T) {
	t.Parallel()
	cookiesPolicy := cookies.Policy{
		Essential: true,
		Usage:     false,
	}
	expectedModel := models.CookiesPreference{}
	expectedModel.Breadcrumb = []coreModel.TaxonomyNode{
		{
			Title: "Home",
			URI:   "/",
		},
		{
			Title: "Cookies",
		},
	}
	expectedModel.Language = "en"
	expectedModel.Metadata.Title = "Cookies"
	expectedModel.CookiesPreferencesSet = true
	expectedModel.CookiesPolicy.Essential = true
	expectedModel.CookiesPolicy.Usage = false
	expectedModel.PreferencesUpdated = false
	expectedModel.FeatureFlags.HideCookieBanner = true
	expectedModel.SiteDomain = "abcd"
	expectedModel.PatternLibraryAssetsPath = "1234"

	Convey("test CreateCookieSettingPage", t, func() {
		mcp := CreateCookieSettingPage(&config.Config{
			PatternLibraryAssetsPath: "1234",
			SiteDomain:               "abcd",
		}, cookiesPolicy, false, request.DefaultLang)
		fmt.Printf("%+v\n", mcp)
		So(expectedModel, ShouldResemble, mcp)
	})
}

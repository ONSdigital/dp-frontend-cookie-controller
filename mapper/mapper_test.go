package mapper

import (
	"dp-frontend-cookie-controller/model"
	"fmt"
	"testing"

	"github.com/ONSdigital/dp-cookies/cookies"
	request "github.com/ONSdigital/dp-net/request"
	coreModel "github.com/ONSdigital/dp-renderer/model"
	. "github.com/smartystreets/goconvey/convey"
)

// TestUnitMapper tests mapper functions
func TestUnitMapper(t *testing.T) {
	t.Parallel()
	cookiesPolicy := cookies.Policy{
		Essential: true,
		Usage:     false,
	}
	expectedModel := model.CookiesPreference{}
	expectedModel.Breadcrumb = []coreModel.TaxonomyNode{
		{
			Title: "Home",
			URI:   "/",
		},
		{
			Title: "Cookies",
		},
	}
	expectedModel.PatternLibraryAssetsPath = "path/to/assets"
	expectedModel.SiteDomain = "site-domain"
	expectedModel.Language = "en"
	expectedModel.Metadata.Title = "Cookies"
	expectedModel.CookiesPreferencesSet = true
	expectedModel.CookiesPolicy.Essential = true
	expectedModel.CookiesPolicy.Usage = false
	expectedModel.PreferencesUpdated = false
	expectedModel.FeatureFlags.HideCookieBanner = true

	basePage := coreModel.NewPage("path/to/assets", "site-domain")
	Convey("test CreateCookieSettingPage", t, func() {
		mcp := CreateCookieSettingPage(basePage, cookiesPolicy, false, request.DefaultLang)
		fmt.Printf("%+v\n", mcp)
		So(expectedModel, ShouldResemble, mcp)
	})
}

package mapper

import (
	"fmt"
	"testing"

	"github.com/ONSdigital/dp-cookies/cookies"
	"github.com/ONSdigital/dp-frontend-cookie-controller/model"
	"github.com/ONSdigital/dp-net/v3/request"
	coreModel "github.com/ONSdigital/dp-renderer/v2/model"
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
	}
	expectedModel.PatternLibraryAssetsPath = "path/to/assets"
	expectedModel.SiteDomain = "site-domain"
	expectedModel.Language = "en"
	expectedModel.Metadata.Title = CookiesStr
	expectedModel.CookiesPreferencesSet = true
	expectedModel.CookiesPolicy.Essential = true
	expectedModel.CookiesPolicy.Usage = false
	expectedModel.PreferencesUpdated = false
	expectedModel.FeatureFlags.HideCookieBanner = true
	expectedModel.UsageRadios = coreModel.RadioFieldset{
		Radios: []coreModel.Radio{
			{
				Input: coreModel.Input{
					ID:        "usage-on",
					IsChecked: false,
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
					IsChecked: true,
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

	basePage := coreModel.NewPage("path/to/assets", "site-domain")
	Convey("test CreateCookieSettingPage", t, func() {
		mcp := CreateCookieSettingPage(basePage, cookiesPolicy, false, request.DefaultLang)
		fmt.Printf("%+v\n", mcp)
		So(expectedModel, ShouldResemble, mcp)
	})
}

package mapper

import (
	"fmt"
	"testing"

	coreModel "github.com/ONSdigital/dis-design-system-go/v2/model"
	"github.com/ONSdigital/dp-cookies/cookies"
	"github.com/ONSdigital/dp-frontend-cookie-controller/model"
	"github.com/ONSdigital/dp-net/v3/request"
	. "github.com/smartystreets/goconvey/convey"
)

// TestUnitMapper tests mapper functions
func TestUnitMapper(t *testing.T) {
	t.Parallel()
	cookiesPolicy := cookies.ONSPolicy{
		Campaigns: false,
		Essential: true,
		Usage:     false,
		Settings:  false,
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
	expectedModel.CookiesPolicy = coreModel.CookiesPolicy{
		Communications: cookiesPolicy.Campaigns,
		Essential:      cookiesPolicy.Essential,
		Settings:       cookiesPolicy.Settings,
		Usage:          cookiesPolicy.Usage,
	}
	expectedModel.PreferencesUpdated = false
	expectedModel.FeatureFlags.HideCookieBanner = true
	expectedModel.UsageRadios = coreModel.RadioFieldset{
		HasBorder: true,
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
	expectedModel.CommsRadios = coreModel.RadioFieldset{
		HasBorder: true,
		Radios: []coreModel.Radio{
			{
				Input: coreModel.Input{
					ID:        "comms-on",
					IsChecked: false,
					Label: coreModel.Localisation{
						LocaleKey: "On",
						Plural:    1,
					},
					Name:  "cookie-policy-comms",
					Value: "true",
				},
			},
			{
				Input: coreModel.Input{
					ID:        "comms-off",
					IsChecked: true,
					Label: coreModel.Localisation{
						LocaleKey: "Off",
						Plural:    1,
					},
					Name:  "cookie-policy-comms",
					Value: "false",
				},
			},
		},
	}
	expectedModel.SiteSettingsRadios = coreModel.RadioFieldset{
		HasBorder: true,
		Radios: []coreModel.Radio{
			{
				Input: coreModel.Input{
					ID:        "site-settings-on",
					IsChecked: false,
					Label: coreModel.Localisation{
						LocaleKey: "On",
						Plural:    1,
					},
					Name:  "cookie-policy-site-settings",
					Value: "true",
				},
			},
			{
				Input: coreModel.Input{
					ID:        "site-settings-off",
					IsChecked: true,
					Label: coreModel.Localisation{
						LocaleKey: "Off",
						Plural:    1,
					},
					Name:  "cookie-policy-site-settings",
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

package mapper

import (
	"dp-frontend-cookie-controller/model"

	"github.com/ONSdigital/dp-cookies/cookies"
	coreModel "github.com/ONSdigital/dp-renderer/v2/model"
)

const (
	CookiesStr = "Cookies"
)

// CreateCookieSettingPage maps type cookies.Policy to model.Page
func CreateCookieSettingPage(basePage coreModel.Page, policy cookies.Policy, isUpdated bool, lang string) model.CookiesPreference {
	page := model.CookiesPreference{
		Page: basePage,
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
	page.Metadata.Title = CookiesStr
	page.Language = lang
	page.CookiesPreferencesSet = true
	page.CookiesPolicy.Essential = policy.Essential
	page.CookiesPolicy.Usage = policy.Usage
	page.FeatureFlags.HideCookieBanner = true

	// Determine whether or not to show success message. Currently this will
	// be shown when cookies preferences have been updated by the user.
	page.PreferencesUpdated = isUpdated

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
					IsChecked: page.CookiesPolicy.Usage,
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

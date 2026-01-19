package mapper

import (
	coreModel "github.com/ONSdigital/dis-design-system-go/v2/model"
	"github.com/ONSdigital/dp-cookies/cookies"
	"github.com/ONSdigital/dp-frontend-cookie-controller/model"
)

const (
	CookiesStr = "Cookies"
)

// CreateCookieSettingPage maps type cookies.Policy to model.Page
func CreateCookieSettingPage(basePage coreModel.Page, policy cookies.ONSPolicy, isUpdated bool, lang string) model.CookiesPreference {
	page := model.CookiesPreference{
		Page: basePage,
	}
	page.Breadcrumb = []coreModel.TaxonomyNode{
		{
			Title: "Home",
			URI:   "/",
		},
	}
	page.Metadata.Title = CookiesStr
	page.Language = lang
	page.CookiesPreferencesSet = true
	page.CookiesPolicy = coreModel.CookiesPolicy{
		Communications: policy.Campaigns,
		Essential:      policy.Essential,
		Settings:       policy.Settings,
		Usage:          policy.Usage,
	}

	page.FeatureFlags.HideCookieBanner = true

	// Determine whether or not to show success message. Currently this will
	// be shown when cookies preferences have been updated by the user.
	page.PreferencesUpdated = isUpdated

	page.UsageRadios = coreModel.RadioFieldset{
		HasBorder: true,
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

	page.CommsRadios = coreModel.RadioFieldset{
		HasBorder: true,
		Radios: []coreModel.Radio{
			{
				Input: coreModel.Input{
					ID:        "comms-on",
					IsChecked: page.CookiesPolicy.Communications,
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
					IsChecked: !page.CookiesPolicy.Communications,
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

	page.SiteSettingsRadios = coreModel.RadioFieldset{
		HasBorder: true,
		Radios: []coreModel.Radio{
			{
				Input: coreModel.Input{
					ID:        "site-settings-on",
					IsChecked: page.CookiesPolicy.Settings,
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
					IsChecked: !page.CookiesPolicy.Settings,
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

	return page
}

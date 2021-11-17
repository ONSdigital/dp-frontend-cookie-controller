package mapper

import (
	"dp-frontend-cookie-controller/model"

	"github.com/ONSdigital/dp-cookies/cookies"
	coreModel "github.com/ONSdigital/dp-renderer/model"
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
	page.Metadata.Title = "Cookies"
	page.Language = lang
	page.CookiesPreferencesSet = true
	page.CookiesPolicy.Essential = policy.Essential
	page.CookiesPolicy.Usage = policy.Usage
	page.FeatureFlags.HideCookieBanner = true
	page.FeatureFlags.SixteensVersion = "67f6982"

	// Determine whether or not to show success message. Currently this will
	// be shown when cookies preferences have been updated by the user.
	page.PreferencesUpdated = isUpdated

	return page
}

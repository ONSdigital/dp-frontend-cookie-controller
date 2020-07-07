package mapper

import (
	"github.com/ONSdigital/dp-cookies/cookies"
	"github.com/ONSdigital/dp-frontend-models/model"
	"github.com/ONSdigital/dp-frontend-models/model/cookiespreferences"
)

// CreateCookieSettingPage maps type cookies.Policy to model.Page
func CreateCookieSettingPage(policy cookies.Policy, isUpdated bool) cookiespreferences.Page {
	var page cookiespreferences.Page
	page.Breadcrumb = []model.TaxonomyNode{
		{
			Title: "Home",
			URI:   "/",
		},
		{
			Title: "Cookies",
		},
	}
	page.Metadata.Title = "Cookies"
	page.CookiesPreferencesSet = true
	page.CookiesPolicy.Essential = policy.Essential
	page.CookiesPolicy.Usage = policy.Usage
	page.FeatureFlags.HideCookieBanner = true

	// Determine whether or not to show success message. Currently this will be shown when cookies preferences have been updated by the user.
	page.PreferencesUpdated = isUpdated

	return page
}

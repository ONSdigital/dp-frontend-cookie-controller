package mapper

import (
	"dp-frontend-cookie-controller/config"
	"dp-frontend-cookie-controller/models"

	"github.com/ONSdigital/dp-cookies/cookies"
	coreModel "github.com/rav-pradhan/test-modules/render/models"
)

// CreateCookieSettingPage maps type cookies.Policy to model.Page
func CreateCookieSettingPage(cfg *config.Config, policy cookies.Policy, isUpdated bool, lang string) models.CookiesPreference {
	page := models.CookiesPreference{
		Page: *coreModel.NewPage(cfg.PatternLibraryAssetsPath, cfg.SiteDomain),
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
	page.SiteDomain = cfg.SiteDomain
	page.PatternLibraryAssetsPath = cfg.PatternLibraryAssetsPath

	// Determine whether or not to show success message. Currently this will be shown when cookies preferences have been updated by the user.
	page.PreferencesUpdated = isUpdated

	return page
}

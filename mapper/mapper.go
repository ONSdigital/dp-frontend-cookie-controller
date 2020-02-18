package mapper

import (
	"github.com/ONSdigital/dp-cookies/cookies"
	"github.com/ONSdigital/dp-frontend-models/model"
)

// CreateCookieSettingPage maps type cookies.Policy to model.Page
func CreateCookieSettingPage(policy cookies.Policy) model.Page {
	var page model.Page
	page.Breadcrumb = []model.TaxonomyNode{
		{
			Title: "Home",
			URI:   "/",
		},
		{
			Title: "Cookies",
		},
	}
	page.CookiesPreferenceSet = true
	page.CookiesPolicy.Essential = policy.Essential
	page.CookiesPolicy.Usage = policy.Usage
	return page
}

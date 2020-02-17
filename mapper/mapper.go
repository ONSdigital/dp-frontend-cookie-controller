package mapper

import (
	"github.com/ONSdigital/dp-cookies/cookies"
	"github.com/ONSdigital/dp-frontend-models/model"
)

// CreateCookieSettingPage maps type cookies.Policy to model.CookiesPolicy
func CreateCookieSettingPage(policy cookies.Policy) model.CookiesPolicy {
	return model.CookiesPolicy(policy)
}

package model

import "github.com/ONSdigital/dis-design-system-go/model"

// CookiesPreference is the model struct for the cookies preferences form
type CookiesPreference struct {
	model.Page
	PreferencesUpdated bool                `json:"preferences_updated"`
	UsageRadios        model.RadioFieldset `json:"usage_radios"`
}

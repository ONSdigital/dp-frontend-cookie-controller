package model

import "github.com/ONSdigital/dp-renderer/v2/model"

// CookiesPreference is the model struct for the cookies preferences form
type CookiesPreference struct {
	model.Page
	CommsRadios        model.RadioFieldset `json:"comms_radios"`
	PreferencesUpdated bool                `json:"preferences_updated"`
	SiteSettingsRadios model.RadioFieldset `json:"site_settings_radios"`
	UsageRadios        model.RadioFieldset `json:"usage_radios"`
}

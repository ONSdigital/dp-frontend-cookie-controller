package models

import "github.com/rav-pradhan/test-modules/render/models"

// Page model data for the cookies preferences form
type CookiesPreference struct {
	models.Page
	PreferencesUpdated bool `json:"preferences_updated"`
}

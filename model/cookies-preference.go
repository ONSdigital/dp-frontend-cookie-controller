package model

import "github.com/ONSdigital/dp-renderer/model"

// Page model data for the cookies preferences form
type CookiesPreference struct {
	model.Page
	PreferencesUpdated bool `json:"preferences_updated"`
}

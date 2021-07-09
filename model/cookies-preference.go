package model

import "github.com/ONSdigital/dp-renderer/model"

// CookiesPreference is the model struct for the cookies preferences form
type CookiesPreference struct {
	model.Page
	PreferencesUpdated bool `json:"preferences_updated"`
}

package model

import "github.com/ONSdigital/dp-renderer/v2/model"

// CookiesPreference is the model struct for the cookies preferences form
type CookiesPreference struct {
	model.Page
	PreferencesUpdated bool                `json:"preferences_updated"`
	TypeRadios         model.RadioFieldset `json:"type_radios"`
}

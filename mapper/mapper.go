package mapper

import (
	"context"
	"encoding/json"
	"github.com/ONSdigital/dp-frontend-models/model"
	"github.com/ONSdigital/log.go/log"
)

func CreateCookieSettingPage(ctx context.Context, b []byte) model.CookiesPolicy {
	var cp model.CookiesPolicy
	if err := json.Unmarshal(b, &cp); err != nil {
		log.Event(ctx, "unable to unmarshal cookie", log.Error(err))
	}
	return cp
}

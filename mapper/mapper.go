package mapper

import (
	"context"
)

type CookiesPolicy struct {
	Essential    bool    `json:essential`
	Usage        bool    `json:usage`
}
type HelloModel struct {
	Greeting string `json:"greeting"`
	Who      string `json:"who"`
}

type HelloWorldModel struct {
	HelloWho string `json:"hello-who"`
}

func CreateCookieSettingPage(ctx context.Context, b []byte) CookiesPolicy {
	cp := CookiesPolicy{
		Essential: true,
		Usage:     false,
	}
	return cp
}

package mapper

import (
	"context"
	"dp-frontend-cookie-controller/config"
	"fmt"
)

type HelloModel struct {
	Greeting string `json:"greeting"`
	Who      string `json:"who"`
}

type HelloWorldModel struct {
	HelloWho string `json:"hello-who"`
}

func CreateCookieSettingPage(ctx context.Context, []byte) HelloWorldModel {
	var hwm HelloWorldModel
	hwm.HelloWho = fmt.Sprintf("%s %s", hm.Greeting, hm.Who)
	if cfg.Emphasise {
		hwm.HelloWho += "!"
	}
	return hwm
}

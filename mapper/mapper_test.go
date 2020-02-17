package mapper

import (
	"github.com/ONSdigital/dp-cookies/cookies"
	"github.com/ONSdigital/dp-frontend-models/model"

	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// TestUnitMapper tests mapper functions
func TestUnitMapper(t *testing.T) {
	t.Parallel()
	cookiesPolicy := cookies.Policy{
		Essential: true,
		Usage:     false,
	}
	idealModelCookiesPolicy := model.CookiesPolicy{
		Essential: true,
		Usage:     false,
	}
	Convey("test CreateCookieSettingPage", t, func() {
		mcp := CreateCookieSettingPage(cookiesPolicy)
		So(idealModelCookiesPolicy, ShouldResemble, mcp)
	})
}

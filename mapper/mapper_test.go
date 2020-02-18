package mapper

import (
	"fmt"

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
	expectedModel := model.Page{
		Breadcrumb: []model.TaxonomyNode{
			{
				Title: "Home",
				URI:   "/",
			},
			{
				Title: "Cookies",
			},
		},
		CookiesPolicy: model.CookiesPolicy{
			Essential: true,
			Usage:     false,
		},
		CookiesPreferenceSet: true,
	}
	Convey("test CreateCookieSettingPage", t, func() {
		mcp := CreateCookieSettingPage(cookiesPolicy)
		fmt.Printf("%+v\n", mcp)
		So(expectedModel, ShouldResemble, mcp)
	})
}

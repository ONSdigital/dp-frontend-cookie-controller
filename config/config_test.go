package config

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

// TestConfig tests config options correctly default if not set
func TestConfig(t *testing.T) {
	t.Parallel()
	Convey("Given an environment with no environment variables set", t, func() {
		cfg, err := Get()

		Convey("When the config values are retrieved", func() {

			Convey("Then there should be no error returned", func() {
				So(err, ShouldBeNil)
			})

			Convey("That the values should be set to the expected defaults", func() {
				So(cfg.GracefulShutdownTimeout, ShouldEqual, 5*time.Second)
				So(cfg.HealthCheckInterval, ShouldEqual, 10*time.Second)
				So(cfg.HealthCheckCriticalTimeout, ShouldEqual, time.Minute)
				So(cfg.BindAddr, ShouldEqual, ":24100")
				So(cfg.RendererURL, ShouldEqual, "http://localhost:20010")
			})
		})
	})
}

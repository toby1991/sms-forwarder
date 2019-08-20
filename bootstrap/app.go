package bootstrap

import (
	"github.com/totoval/framework/helpers/zone"
	"github.com/totoval/framework/http/middleware"
	"github.com/totoval/framework/logs"
	"github.com/totoval/framework/request"
	"github.com/totoval/framework/sentry"
	"github.com/totoval/framework/validator"

	"totoval/config"
	"totoval/resources/lang"

	c "github.com/totoval/framework/config"
)

func Initialize() {
	config.Initialize()
	sentry.Initialize()
	logs.Initialize()
	zone.Initialize()
	lang.Initialize() // an translation must contains resources/lang/xx.json file (then a resources/lang/validation_translator/xx.go)
	// cache.Initialize()
	// database.Initialize()
	// m.Initialize()
	// queue.Initialize()
	// jobs.Initialize()
	// events.Initialize()
	// listeners.Initialize()

	validator.UpgradeValidatorV8toV9()
}

func Middleware(r *request.Engine) {
	r.Use(middleware.RequestLogger())

	if c.GetString("app.env") == "production" {
		r.Use(middleware.Logger())
		r.Use(middleware.Recovery())
	}

	r.Use(middleware.Locale())
}

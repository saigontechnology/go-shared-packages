package cors

import (
	gincors "github.com/gin-contrib/cors"
)

// HealthCorsConfig allow to have x-locale in header
// This is similar to config in gincors.Default() but allows x-locale.
func HealthCorsConfig() gincors.Config {
	corsCfg := gincors.DefaultConfig()
	corsCfg.AllowAllOrigins = true
	corsCfg.AllowHeaders = append(corsCfg.AllowHeaders, "x-locale")

	return corsCfg
}

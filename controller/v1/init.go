package v1

import "github.com/JDSchmitz/kolide-archive/config"

var (
	configuration *config.Config
)

// Init v1 contollers
func Init(c *config.Config) {
	configuration = c
}

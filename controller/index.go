package controller

import (
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/JDSchmitz/kolide-archive/model"
	"github.com/JDSchmitz/kolide-archive/router/middleware/session"
	"github.com/JDSchmitz/kolide-archive/shared/osquery"
	"github.com/JDSchmitz/kolide-archive/shared/token"
)

// Index route
func Index(c *gin.Context) {
	var csrf string
	var template = "index"

	user := session.User(c)

	if user == nil {
		template = "login"
	} else {
		csrf, _ = token.New(
			token.CsrfToken,
			user.Email,
			user,
		).Sign(user.Hash)

		queries := make(chan osquery.Query)
		c.Set("Queries", &queries)
	}

	nodes, err := model.AllNodes(nil)

	if err != nil {
		log.Error(err.Error())
	}

	c.HTML(200, "layout.html", gin.H{
		"Template":  template,
		"User":      user,
		"Csrf":      csrf,
		"Nodes":     nodes,
		"timestamp": time.Now().Unix(),
	})

}

// Login route
func Login(c *gin.Context) {
	Index(c)
}

package v1

import (
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/JDSchmitz/kolide-archive/controller/helpers"
	"github.com/JDSchmitz/kolide-archive/model"
)

// DeleteSavedQuery route
func DeleteSavedQuery(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		helpers.JsonError(c, 500, err)
		return
	}

	q, err := model.FindSavedQueryById(id)

	if err != nil {
		helpers.JsonError(c, 500, err)
		return
	}

	if err := q.Delete(); err != nil {
		helpers.JsonError(c, 500, err)
		return
	}

	helpers.JsonResp(c, 200, gin.H{
		"error": nil,
	})
}

// SavedQueries route
func SavedQueries(c *gin.Context) {
	data, err := model.AllSavedQueries()

	if err != nil {
		helpers.JsonError(c, 500, err)
		return
	}

	helpers.JsonResp(c, 200, gin.H{
		"queries": data,
		"error":   nil,
	})
}

// CreateSavedQuery route
func CreateSavedQuery(c *gin.Context) {
	query := model.SavedQuery{}

	if err := json.Unmarshal(helpers.GetBody(c), &query); err != nil {
		helpers.JsonError(c, 500, err)
		return
	}

	err := model.NewSavedQuery(query)

	if err != nil {
		helpers.JsonError(c, 500, err)
		return
	}

	helpers.JsonResp(c, 200, gin.H{
		"query": query,
		"error": nil,
	})
}

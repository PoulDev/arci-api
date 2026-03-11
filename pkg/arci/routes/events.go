package routes

import (
	"time"

	"arci.it/pkg/arci/db"
	"github.com/gin-gonic/gin"
)

type EventData struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Date        time.Time      `json:"date"`
	Roles       []db.RoleEvent `json:"roles"`
}

func GetRoles(c *gin.Context) {
	roles, err := db.GetRoles()
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"roles": roles,
	})
}

func GetEvents(c *gin.Context) {
	events, err := db.GetEvents()
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"events": events,
	})
}

func NewEvent(c *gin.Context) {
	var eventData EventData
	if err := c.ShouldBindJSON(&eventData); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := db.AddEvent(eventData.Name, eventData.Description, eventData.Date, eventData.Roles)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"ok": true,
	})
}

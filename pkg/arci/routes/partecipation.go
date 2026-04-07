package routes

import (
	"strconv"

	"arci.it/pkg/arci/db"
	"github.com/gin-gonic/gin"
)

type PartecipationData struct {
	Role string `json:"role" binding:"required"`
}

func Partecipate(c *gin.Context) {
	eventID, err := strconv.Atoi(c.Param("event_id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid event ID"})
		return
	}

	var data PartecipationData
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	memberID := c.GetInt("member_id")

	err = db.Partecipate(eventID, memberID, data.Role)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"ok": true})
}

func CancelPartecipation(c *gin.Context) {
	eventID, err := strconv.Atoi(c.Param("event_id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid event ID"})
		return
	}
 
	memberID := c.GetInt("member_id")
 
	err = db.CancelPartecipation(eventID, memberID)
	if err != nil {
		if err.Error() == "partecipation not found" {
			c.JSON(404, gin.H{"error": err.Error()})
			return
		}
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
 
	c.JSON(200, gin.H{"ok": true})
}
 

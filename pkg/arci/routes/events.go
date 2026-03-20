package routes

import (
	"strconv"
	"time"

	"arci.it/pkg/arci/db"
	"github.com/gin-gonic/gin"
)

type EventData struct {
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Date        time.Time           `json:"date"`
	Roles       []db.RoleEventInput `json:"roles"`
}

type RoleData struct {
	Name string `json:"name" binding:"required"`
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

func NewRole(c *gin.Context) {
	var roleData RoleData
	if err := c.ShouldBindJSON(&roleData); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	role, err := db.AddRole(roleData.Name)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(201, gin.H{
		"ok":   true,
		"role": role,
	})
}

func DeleteRole(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid role ID",
		})
		return
	}

	err = db.DeleteRole(id)
	if err != nil {
		if err.Error() == "role not found" {
			c.JSON(404, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "role is assigned to one or more events" {
			c.JSON(409, gin.H{"error": err.Error()})
			return
		}
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"ok": true})
}

func GetEvents(c *gin.Context) {
	memberID := c.GetInt("member_id")

	events, err := db.GetEvents(memberID)
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

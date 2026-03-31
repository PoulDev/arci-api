package main

import (
	"fmt"
	"time"

	"arci.it/pkg/arci/db"
	"arci.it/pkg/arci/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	err := db.ConnectDatabase()
	if err != nil {
		fmt.Println(err)
		return
	}

	router := gin.Default()

	router.POST("/register", routes.Register)
	router.POST("/login", routes.Login)

	protected := router.Group("/")
	protected.Use(routes.AuthMiddleware())
	{
		protected.GET("/events", routes.GetEvents)
		protected.POST("/events", routes.NewEvent)

		protected.GET("/roles", routes.GetRoles)
	}

	admin := router.Group("/")
	admin.Use(routes.AuthMiddleware(), routes.AdminMiddleware())
	{
		admin.POST("/roles", routes.NewRole)
		admin.DELETE("/roles/:id", routes.DeleteRole)
	}

	fmt.Println(time.Now().String())

	err = router.Run(":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
}

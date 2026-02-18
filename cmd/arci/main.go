package main

import (
	"fmt"

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
	}

	err = router.Run(":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
}

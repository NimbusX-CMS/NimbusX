package main

import (
	"fmt"

	"github.com/NimbusX-CMS/NimbusX/api/internal/api"
	"github.com/NimbusX-CMS/NimbusX/api/internal/db/multi_db"
	"github.com/NimbusX-CMS/NimbusX/api/internal/resources"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	db, err := multi_db.ConnectToSQLite("db.sqlite")
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	err = db.EnsureTablesCreation()
	if err != nil {
		fmt.Println("Error creating tables:", err)
		return
	}
	server := &resources.Server{
		DB: db,
	}

	api.RegisterHandlers(router, server)

	err = router.Run(":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
	fmt.Println("Server started")
}

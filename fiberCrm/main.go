package main

import (
	"fmt"

	"github.com/gofiber/fiber"
	"github.com/robinmuhia/GolangProjects/fiberCrm/database"
	"github.com/robinmuhia/GolangProjects/fiberCrm/lead"
)

func setupRoutes(app *fiber.App) {
	app.Get("/api/v1/lead",lead.GetLeads)
	app.Get("/api/v1/lead:id",lead.GetLead)
	app.Post("/api/v1/lead",lead.NewLeads)
	app.Delete("/api/v1/lead:id",lead.DeleteLead)
}

func initDatabase(){
	database.Connect()
	db := database.GetDB()
	db.AutoMigrate(&lead.Lead{})
	fmt.Println("Database Migrated")
}

func main() {
	app := fiber.New()
	initDatabase()
	setupRoutes(app)
	app.Listen(3000)
	defer database.DBConn.Close()
}
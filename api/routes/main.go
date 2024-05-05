package routes

import (
	"CatsSocial/api/handlers"
	"CatsSocial/db/functions"

	"github.com/gofiber/fiber/v2"
)

func RouteRegister(app *fiber.App, deps handlers.Dependencies) {
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	userHandler := handlers.User{
		Database: functions.NewUser(deps.DbPool, deps.Cfg),
	}

	UserRoutes(app, userHandler)

	catHandler := handlers.Cat{
		Database:     functions.NewCatFn(deps.DbPool),
		UserDatabase: functions.NewUser(deps.DbPool, deps.Cfg),
	}

	CatRoutes(app, catHandler)

	matchHandler := handlers.MatchHandler{
		Match:        *functions.NewMatch(deps.DbPool),
		CatDatabase:  functions.NewCatFn(deps.DbPool),
		UserDatabase: functions.NewUser(deps.DbPool, deps.Cfg),
	}

	MatchRoutes(app, matchHandler)
}

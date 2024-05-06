package routes

import (
	"CatsSocial/api/handlers"
	"CatsSocial/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func MatchRoutes(app *fiber.App, h handlers.MatchHandler) {
	g := app.Group("/v1/cat/match").Use(middleware.JWTAuth())

	g.Post("", h.Create)
	g.Get("", h.Get)
	g.Post("/approve", h.Approve)
	g.Post("/reject", h.Reject)
	g.Delete("/:id", h.Delete)
}

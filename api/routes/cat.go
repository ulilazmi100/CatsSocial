package routes

import (
	"CatsSocial/api/handlers"
	"CatsSocial/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func CatRoutes(app *fiber.App, h handlers.Cat) {
	g := app.Group("/v1/cat").Use(middleware.JWTAuth())
	g.Get("", h.GetCats)
	g.Post("", h.AddCat)
	g.Put("/:id", h.UpdateCat)
	g.Delete("/:id", h.DeleteCat)
}

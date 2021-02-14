package main

import (
	"github.com/gofiber/fiber"
)

func accessTokenMiddleware(ctx *fiber.Ctx) {
	if ctx.Get("Access-Token") == "" {
		errorHandler(ctx, fiber.StatusUnauthorized, "Invalid access token provided!")
		return
	}

	if ctx.Get("Access-Token") != cdnConfig.AccessToken {
		errorHandler(ctx, fiber.StatusUnauthorized, "Invalid access token provided!")
		return
	}

	ctx.Next()
}

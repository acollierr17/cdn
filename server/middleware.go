package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
)

func accessTokenMiddleware(ctx *fiber.Ctx) error {
	accessToken := ctx.Get("Access-Token")

	if ctx.Path() == "/api/token" {
		return ctx.Next()
	}

	if accessToken == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "No access token provided!")
	}

	docSnap, err := firebaseFirestore.Collection("tokens").Where("token", "==", accessToken).Documents(context.Background()).Next()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if docSnap == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid access token provided!")
	}

	return ctx.Next()
}

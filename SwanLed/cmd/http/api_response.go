package http

import "github.com/gofiber/fiber/v2"

type ApiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func SendErrorResponse(ctx *fiber.Ctx, status int, message string) error {
	return ctx.Status(status).JSON(ApiResponse{
		Success: false,
		Error:   message,
	})
}

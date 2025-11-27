package utils

import "github.com/gofiber/fiber/v2"

type ErrorResponse struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

func Error(c *fiber.Ctx, status int, code, message string, fields map[string]string) error {
	return c.Status(status).JSON(fiber.Map{"error": ErrorResponse{Code: code, Message: message, Fields: fields}})
}

func Success(c *fiber.Ctx, status int, data interface{}) error {
	return c.Status(status).JSON(fiber.Map{"data": data})
}

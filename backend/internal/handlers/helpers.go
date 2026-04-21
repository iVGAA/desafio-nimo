package handlers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func badRequest(c *fiber.Ctx, msg string) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": msg})
}

func notFound(c *fiber.Ctx, msg string) error {
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": msg})
}

func conflict(c *fiber.Ctx, msg string) error {
	return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": msg})
}

func internal(c *fiber.Ctx, msg string) error {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": msg})
}

func parseDate(val string) (time.Time, error) {
	if val == "" {
		return time.Time{}, fmt.Errorf("data vazia")
	}
	return time.Parse("2006-01-02", val)
}

func queryInt64(c *fiber.Ctx, key string) (*int64, error) {
	val := c.Query(key)
	if val == "" {
		return nil, nil
	}
	parsed, err := strconv.ParseInt(val, 10, 64)
	if err != nil || parsed <= 0 {
		return nil, fmt.Errorf("%s inválido", key)
	}
	return &parsed, nil
}

func queryInt(c *fiber.Ctx, key string) (*int, error) {
	val := c.Query(key)
	if val == "" {
		return nil, nil
	}
	parsed, err := strconv.Atoi(val)
	if err != nil || parsed < 1 {
		return nil, fmt.Errorf("%s inválido", key)
	}
	return &parsed, nil
}

func queryDate(c *fiber.Ctx, key string) (*time.Time, error) {
	val := c.Query(key)
	if val == "" {
		return nil, nil
	}
	parsed, err := parseDate(val)
	if err != nil {
		return nil, fmt.Errorf("%s inválida; use YYYY-MM-DD", key)
	}
	return &parsed, nil
}

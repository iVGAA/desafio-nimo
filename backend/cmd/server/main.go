package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/desafionimo/backend/internal/config"
	"github.com/desafionimo/backend/internal/database"
	"github.com/desafionimo/backend/internal/handlers"
	"github.com/desafionimo/backend/internal/repositories"
	"github.com/desafionimo/backend/internal/routes"
	"github.com/desafionimo/backend/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}
	port, err := cfg.MustPort()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	pool, err := database.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("banco: %v", err)
	}
	defer pool.Close()

	produtoRepo := repositories.NewProdutoRepository(pool)
	fornecedorRepo := repositories.NewFornecedorRepository(pool)
	precoRepo := repositories.NewPrecoRepository(pool)

	produtoSvc := services.NewProdutoService(produtoRepo)
	fornecedorSvc := services.NewFornecedorService(fornecedorRepo)
	precoSvc := services.NewPrecoService(precoRepo, produtoRepo, fornecedorRepo)

	produtoH := handlers.NewProdutoHandler(produtoSvc)
	fornecedorH := handlers.NewFornecedorHandler(fornecedorSvc)
	precoH := handlers.NewPrecoHandler(precoSvc)

	app := fiber.New(fiber.Config{
		AppName:      "Combustível API",
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.AllowOrigins,
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: false,
	}))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	routes.Register(app, routes.Deps{
		Produtos:     produtoH,
		Fornecedores: fornecedorH,
		Precos:       precoH,
	})

	go func() {
		addr := fmt.Sprintf(":%d", port)
		log.Printf("servidor escutando em http://localhost%s", addr)
		if err := app.Listen(addr); err != nil {
			log.Printf("fiber: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}

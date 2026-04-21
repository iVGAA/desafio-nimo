package routes

import (
	"github.com/desafionimo/backend/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

// Deps agrupa handlers registrados nas rotas.
type Deps struct {
	Produtos     *handlers.ProdutoHandler
	Fornecedores *handlers.FornecedorHandler
	Precos       *handlers.PrecoHandler
}

// Register monta todas as rotas com prefixo /api/v1.
func Register(app *fiber.App, d Deps) {
	v1 := app.Group("/api/v1")

	v1.Get("/produtos", d.Produtos.List)
	v1.Post("/produtos", d.Produtos.Create)

	v1.Get("/fornecedores", d.Fornecedores.List)
	v1.Post("/fornecedores", d.Fornecedores.Create)

	v1.Post("/precos", d.Precos.Create)
	v1.Get("/precos", d.Precos.List)
	v1.Get("/precos/historico/:produto_id", d.Precos.Historico)
	v1.Get("/precos/comparativo", d.Precos.Comparativo)
}

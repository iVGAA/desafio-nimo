package handlers

import (
	"errors"

	"github.com/desafionimo/backend/internal/repositories"
	"github.com/desafionimo/backend/internal/services"
	"github.com/gofiber/fiber/v2"
)

// ProdutoHandler expõe endpoints de produtos.
type ProdutoHandler struct {
	svc *services.ProdutoService
}

func NewProdutoHandler(svc *services.ProdutoService) *ProdutoHandler {
	return &ProdutoHandler{svc: svc}
}

type createProdutoRequest struct {
	Nome    string `json:"nome"`
	Unidade string `json:"unidade"`
}

// List retorna todos os produtos cadastrados.
func (handler *ProdutoHandler) List(c *fiber.Ctx) error {
	items, err := handler.svc.List(c.UserContext())
	if err != nil {
		return internal(c, "não foi possível listar produtos")
	}
	return c.JSON(items)
}

// Create cadastra um novo produto, validando duplicações de nome.
func (handler *ProdutoHandler) Create(c *fiber.Ctx) error {
	var body createProdutoRequest
	if err := c.BodyParser(&body); err != nil {
		return badRequest(c, "corpo JSON inválido")
	}
	
	produto, err := handler.svc.Create(c.UserContext(), body.Nome, body.Unidade)
	if err != nil {
		if errors.Is(err, repositories.ErrDuplicateProdutoNome) {
			return conflict(c, repositories.ErrDuplicateProdutoNome.Error())
		}
		return badRequest(c, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(produto)
}

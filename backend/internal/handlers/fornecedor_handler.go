package handlers

import (
	"errors"

	"github.com/desafionimo/backend/internal/repositories"
	"github.com/desafionimo/backend/internal/services"
	"github.com/gofiber/fiber/v2"
)

// FornecedorHandler expõe endpoints de fornecedores.
type FornecedorHandler struct {
	svc *services.FornecedorService
}

func NewFornecedorHandler(svc *services.FornecedorService) *FornecedorHandler {
	return &FornecedorHandler{svc: svc}
}

type createFornecedorRequest struct {
	Nome string `json:"nome"`
	CNPJ string `json:"cnpj"`
}

// List retorna todos os fornecedores cadastrados.
func (handler *FornecedorHandler) List(c *fiber.Ctx) error {
	items, err := handler.svc.List(c.UserContext())
	if err != nil {
		return internal(c, "não foi possível listar fornecedores")
	}
	return c.JSON(items)
}

// Create cadastra um novo fornecedor, validando duplicações de nome e CNPJ.
func (handler *FornecedorHandler) Create(c *fiber.Ctx) error {
	var body createFornecedorRequest
	if err := c.BodyParser(&body); err != nil {
		return badRequest(c, "corpo JSON inválido")
	}
	
	fornecedor, err := handler.svc.Create(c.UserContext(), body.Nome, body.CNPJ)
	if err != nil {
		if errors.Is(err, repositories.ErrDuplicateFornecedorNome) {
			return conflict(c, repositories.ErrDuplicateFornecedorNome.Error())
		}
		if errors.Is(err, repositories.ErrDuplicateCNPJ) {
			return conflict(c, repositories.ErrDuplicateCNPJ.Error())
		}
		return badRequest(c, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(fornecedor)
}

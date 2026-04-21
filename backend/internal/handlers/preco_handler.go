package handlers

import (
	"strconv"

	"github.com/desafionimo/backend/internal/repositories"
	"github.com/desafionimo/backend/internal/services"
	"github.com/gofiber/fiber/v2"
)

// PrecoHandler expõe endpoints de preços.
type PrecoHandler struct {
	svc *services.PrecoService
}

func NewPrecoHandler(svc *services.PrecoService) *PrecoHandler {
	return &PrecoHandler{svc: svc}
}

type createPrecoRequest struct {
	ProdutoID      int64   `json:"produto_id"`
	FornecedorID   int64   `json:"fornecedor_id"`
	PrecoPorLitro  float64 `json:"preco_por_litro"`
	DataReferencia string  `json:"data_referencia"`
}

func (h *PrecoHandler) Create(c *fiber.Ctx) error {
	var body createPrecoRequest
	if err := c.BodyParser(&body); err != nil {
		return badRequest(c, "corpo JSON inválido")
	}
	if body.DataReferencia == "" {
		return badRequest(c, "data_referencia é obrigatória (formato YYYY-MM-DD)")
	}
	
	d, err := parseDate(body.DataReferencia)
	if err != nil {
		return badRequest(c, "data_referencia inválida; use YYYY-MM-DD")
	}

	p, err := h.svc.Register(c.UserContext(), body.ProdutoID, body.FornecedorID, body.PrecoPorLitro, d)
	if err != nil {
		msg := err.Error()
		switch msg {
		case "produto não encontrado", "fornecedor não encontrado":
			return notFound(c, msg)
		default:
			return badRequest(c, msg)
		}
	}
	return c.Status(fiber.StatusCreated).JSON(p)
}

func (h *PrecoHandler) List(c *fiber.Ctx) error {
	f := repositories.PrecoFilters{
		Page:  1,
		Limit: 20,
	}

	fornecedorID, err := queryInt64(c, "fornecedor_id")
	if err != nil {
		return badRequest(c, err.Error())
	}
	f.FornecedorID = fornecedorID

	produtoID, err := queryInt64(c, "produto_id")
	if err != nil {
		return badRequest(c, err.Error())
	}
	f.ProdutoID = produtoID

	dataInicio, err := queryDate(c, "data_inicio")
	if err != nil {
		return badRequest(c, err.Error())
	}
	f.DataInicio = dataInicio

	dataFim, err := queryDate(c, "data_fim")
	if err != nil {
		return badRequest(c, err.Error())
	}
	f.DataFim = dataFim

	page, err := queryInt(c, "page")
	if err != nil {
		return badRequest(c, err.Error())
	}
	if page != nil {
		f.Page = *page
	}

	limit, err := queryInt(c, "limit")
	if err != nil {
		return badRequest(c, err.Error())
	}
	if limit != nil {
		f.Limit = *limit
	}

	res, err := h.svc.List(c.UserContext(), f)
	if err != nil {
		return internal(c, "não foi possível listar preços")
	}
	return c.JSON(res)
}

func (h *PrecoHandler) Historico(c *fiber.Ctx) error {
	id, err := c.ParamsInt("produto_id")
	if err != nil || id < 1 {
		return badRequest(c, "produto_id inválido na URL")
	}
	out, err := h.svc.Historico(c.UserContext(), int64(id))
	if err != nil {
		if err.Error() == "produto não encontrado" {
			return notFound(c, err.Error())
		}
		return internal(c, "não foi possível carregar histórico")
	}
	return c.JSON(out)
}

func (h *PrecoHandler) Comparativo(c *fiber.Ctx) error {
	v := c.Query("produto_id")
	if v == "" {
		return badRequest(c, "query produto_id é obrigatória")
	}
	id, err := strconv.ParseInt(v, 10, 64)
	if err != nil || id < 1 {
		return badRequest(c, "produto_id inválida")
	}
	out, err := h.svc.Comparativo(c.UserContext(), id)
	if err != nil {
		if err.Error() == "produto não encontrado" {
			return notFound(c, err.Error())
		}
		return internal(c, "não foi possível carregar comparativo")
	}
	return c.JSON(out)
}

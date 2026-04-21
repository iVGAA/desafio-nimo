package services

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/desafionimo/backend/internal/models"
	"github.com/desafionimo/backend/internal/repositories"
)

// PrecoService contém regras de negócio para registros de preço.
type PrecoService struct {
	precoRepo      *repositories.PrecoRepository
	produtoRepo    *repositories.ProdutoRepository
	fornecedorRepo *repositories.FornecedorRepository
}

// NewPrecoService inicializa o serviço de preços.
func NewPrecoService(
	precoRepo *repositories.PrecoRepository,
	produtoRepo *repositories.ProdutoRepository,
	fornecedorRepo *repositories.FornecedorRepository,
) *PrecoService {
	return &PrecoService{
		precoRepo:      precoRepo,
		produtoRepo:    produtoRepo,
		fornecedorRepo: fornecedorRepo,
	}
}

// Register cadastra ou atualiza um preço, garantindo que as entidades relacionadas existam.
func (service *PrecoService) Register(ctx context.Context, produtoID, fornecedorID int64, preco float64, dataRef time.Time) (*models.Preco, error) {
	if produtoID <= 0 || fornecedorID <= 0 {
		return nil, fmt.Errorf("produto_id e fornecedor_id devem ser positivos")
	}
	if preco <= 0 || math.IsNaN(preco) || math.IsInf(preco, 0) {
		return nil, fmt.Errorf("preco_por_litro deve ser um número positivo")
	}

	produto, err := service.produtoRepo.GetByID(ctx, produtoID)
	if err != nil {
		return nil, fmt.Errorf("PrecoService.Register verificar produto: %w", err)
	}
	if produto == nil {
		return nil, fmt.Errorf("produto não encontrado")
	}

	fornecedor, err := service.fornecedorRepo.GetByID(ctx, fornecedorID)
	if err != nil {
		return nil, fmt.Errorf("PrecoService.Register verificar fornecedor: %w", err)
	}
	if fornecedor == nil {
		return nil, fmt.Errorf("fornecedor não encontrado")
	}

	precoRegistrado, err := service.precoRepo.Upsert(ctx, produtoID, fornecedorID, preco, dataRef)
	if err != nil {
		switch repositories.PostgresCode(err) {
		case "23503":
			return nil, fmt.Errorf("produto ou fornecedor inválido")
		case "23514":
			return nil, fmt.Errorf("preco_por_litro inválido para o banco (deve ser > 0)")
		default:
			return nil, fmt.Errorf("PrecoService.Register upsert falhou: %w", err)
		}
	}
	return precoRegistrado, nil
}

// List busca preços baseados em filtros com paginação sanitizada.
func (service *PrecoService) List(ctx context.Context, filters repositories.PrecoFilters) (models.PrecoListResult, error) {
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.Limit < 1 {
		filters.Limit = 20
	}
	if filters.Limit > 100 {
		filters.Limit = 100
	}
	return service.precoRepo.List(ctx, filters)
}

// Historico obtém todos os registros de preço de um produto e estrutura para gráficos.
func (service *PrecoService) Historico(ctx context.Context, produtoID int64) (*models.HistoricoProdutoResponse, error) {
	if produtoID <= 0 {
		return nil, fmt.Errorf("produto_id inválido")
	}
	
	produto, err := service.produtoRepo.GetByID(ctx, produtoID)
	if err != nil {
		return nil, fmt.Errorf("PrecoService.Historico verificar produto: %w", err)
	}
	if produto == nil {
		return nil, fmt.Errorf("produto não encontrado")
	}

	pontos, err := service.precoRepo.HistoricoPorProduto(ctx, produtoID)
	if err != nil {
		return nil, fmt.Errorf("PrecoService.Historico extrair pontos: %w", err)
	}
	if pontos == nil {
		pontos = []models.HistoricoPonto{}
	}
	
	return &models.HistoricoProdutoResponse{
		ProdutoID: produtoID,
		Pontos:    pontos,
	}, nil
}

// Comparativo lista o último preço de cada fornecedor para um produto e destaca o menor.
func (service *PrecoService) Comparativo(ctx context.Context, produtoID int64) (*models.ComparativoResponse, error) {
	if produtoID <= 0 {
		return nil, fmt.Errorf("produto_id inválido")
	}
	
	produto, err := service.produtoRepo.GetByID(ctx, produtoID)
	if err != nil {
		return nil, fmt.Errorf("PrecoService.Comparativo verificar produto: %w", err)
	}
	if produto == nil {
		return nil, fmt.Errorf("produto não encontrado")
	}

	itens, err := service.precoRepo.ComparativoAtual(ctx, produtoID)
	if err != nil {
		return nil, fmt.Errorf("PrecoService.Comparativo extrair itens: %w", err)
	}
	if itens == nil {
		itens = []models.ComparativoItem{}
	}

	resp := &models.ComparativoResponse{
		ProdutoID: produtoID,
		Itens:     itens,
	}

	if len(itens) == 0 {
		return resp, nil
	}

	minVal := itens[0].PrecoPorLitro
	minID := itens[0].FornecedorID
	for _, it := range itens[1:] {
		if it.PrecoPorLitro < minVal {
			minVal = it.PrecoPorLitro
			minID = it.FornecedorID
		}
	}
	resp.MenorPreco = &minVal
	resp.MenorPrecoFornecedorID = &minID
	
	return resp, nil
}

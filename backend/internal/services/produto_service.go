package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/desafionimo/backend/internal/models"
	"github.com/desafionimo/backend/internal/repositories"
)

// ProdutoService contém regras de negócio para produtos.
type ProdutoService struct {
	repo *repositories.ProdutoRepository
}

// NewProdutoService inicializa o serviço de produtos.
func NewProdutoService(repo *repositories.ProdutoRepository) *ProdutoService {
	return &ProdutoService{repo: repo}
}

// List retorna todos os produtos disponíveis.
func (service *ProdutoService) List(ctx context.Context) ([]models.Produto, error) {
	return service.repo.List(ctx)
}

// Create cadastra um novo produto preenchendo valor padrão "litro" caso a unidade falte.
func (service *ProdutoService) Create(ctx context.Context, nome, unidade string) (*models.Produto, error) {
	nome = strings.TrimSpace(nome)
	if nome == "" {
		return nil, fmt.Errorf("nome é obrigatório")
	}
	nome = strings.ToLower(nome)
	
	unidade = strings.TrimSpace(unidade)
	if unidade == "" {
		unidade = "litro"
	}
	
	produtoCriado, err := service.repo.Create(ctx, nome, unidade)
	if err != nil {
		return nil, fmt.Errorf("ProdutoService.Create: %w", err)
	}
	return produtoCriado, nil
}

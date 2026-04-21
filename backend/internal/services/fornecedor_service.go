package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/desafionimo/backend/internal/models"
	"github.com/desafionimo/backend/internal/repositories"
)

// FornecedorService contém regras de negócio para fornecedores.
type FornecedorService struct {
	repo *repositories.FornecedorRepository
}

// NewFornecedorService inicializa o serviço de fornecedores.
func NewFornecedorService(repo *repositories.FornecedorRepository) *FornecedorService {
	return &FornecedorService{repo: repo}
}

// List retorna a lista de todos os fornecedores.
func (service *FornecedorService) List(ctx context.Context) ([]models.Fornecedor, error) {
	return service.repo.List(ctx)
}

// Create cadastra um novo fornecedor, garantindo que nome e CNPJ não sejam vazios.
func (service *FornecedorService) Create(ctx context.Context, nome, cnpj string) (*models.Fornecedor, error) {
	nome = strings.TrimSpace(nome)
	if nome == "" {
		return nil, fmt.Errorf("nome é obrigatório")
	}
	nome = strings.ToLower(nome)
	cnpj = strings.TrimSpace(cnpj)
	if cnpj == "" {
		return nil, fmt.Errorf("cnpj é obrigatório")
	}
	
	fornecedorCriado, err := service.repo.Create(ctx, nome, cnpj)
	if err != nil {
		return nil, fmt.Errorf("FornecedorService.Create: %w", err)
	}
	return fornecedorCriado, nil
}

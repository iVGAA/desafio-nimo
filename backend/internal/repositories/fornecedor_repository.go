package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/desafionimo/backend/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// FornecedorRepository persiste e consulta fornecedores.
type FornecedorRepository struct {
	pool *pgxpool.Pool
}

func NewFornecedorRepository(pool *pgxpool.Pool) *FornecedorRepository {
	return &FornecedorRepository{pool: pool}
}

func (r *FornecedorRepository) List(ctx context.Context) ([]models.Fornecedor, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, nome, cnpj, criado_em
		FROM fornecedores
		ORDER BY nome ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("listar fornecedores: %w", err)
	}
	defer rows.Close()

	var out []models.Fornecedor
	for rows.Next() {
		var f models.Fornecedor
		if err := rows.Scan(&f.ID, &f.Nome, &f.CNPJ, &f.CriadoEm); err != nil {
			return nil, fmt.Errorf("scan fornecedor: %w", err)
		}
		out = append(out, f)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iter fornecedores: %w", err)
	}
	return out, nil
}

func (r *FornecedorRepository) Create(ctx context.Context, nome, cnpj string) (*models.Fornecedor, error) {
	row := r.pool.QueryRow(ctx, `
		INSERT INTO fornecedores (nome, cnpj)
		VALUES ($1, $2)
		RETURNING id, nome, cnpj, criado_em
	`, nome, cnpj)

	var f models.Fornecedor
	if err := row.Scan(&f.ID, &f.Nome, &f.CNPJ, &f.CriadoEm); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "uq_fornecedores_nome":
				return nil, ErrDuplicateFornecedorNome
			case "uq_fornecedores_nome_ci":
				return nil, ErrDuplicateFornecedorNome
			case "uq_fornecedores_cnpj":
				return nil, ErrDuplicateCNPJ
			}
		}
		return nil, fmt.Errorf("inserir fornecedor: %w", err)
	}
	return &f, nil
}

func (r *FornecedorRepository) GetByID(ctx context.Context, id int64) (*models.Fornecedor, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, nome, cnpj, criado_em
		FROM fornecedores
		WHERE id = $1
	`, id)

	var f models.Fornecedor
	if err := row.Scan(&f.ID, &f.Nome, &f.CNPJ, &f.CriadoEm); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("buscar fornecedor: %w", err)
	}
	return &f, nil
}

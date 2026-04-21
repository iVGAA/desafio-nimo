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

// ProdutoRepository persiste e consulta produtos.
type ProdutoRepository struct {
	pool *pgxpool.Pool
}

func NewProdutoRepository(pool *pgxpool.Pool) *ProdutoRepository {
	return &ProdutoRepository{pool: pool}
}

func (r *ProdutoRepository) List(ctx context.Context) ([]models.Produto, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, nome, unidade, criado_em
		FROM produtos
		ORDER BY nome ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("listar produtos: %w", err)
	}
	defer rows.Close()

	var out []models.Produto
	for rows.Next() {
		var p models.Produto
		if err := rows.Scan(&p.ID, &p.Nome, &p.Unidade, &p.CriadoEm); err != nil {
			return nil, fmt.Errorf("scan produto: %w", err)
		}
		out = append(out, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iter produtos: %w", err)
	}
	return out, nil
}

func (r *ProdutoRepository) Create(ctx context.Context, nome, unidade string) (*models.Produto, error) {
	row := r.pool.QueryRow(ctx, `
		INSERT INTO produtos (nome, unidade)
		VALUES ($1, $2)
		RETURNING id, nome, unidade, criado_em
	`, nome, unidade)

	var p models.Produto
	if err := row.Scan(&p.ID, &p.Nome, &p.Unidade, &p.CriadoEm); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			if pgErr.ConstraintName == "uq_produtos_nome" || pgErr.ConstraintName == "uq_produtos_nome_ci" {
				return nil, ErrDuplicateProdutoNome
			}
		}
		return nil, fmt.Errorf("inserir produto: %w", err)
	}
	return &p, nil
}

func (r *ProdutoRepository) GetByID(ctx context.Context, id int64) (*models.Produto, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, nome, unidade, criado_em
		FROM produtos
		WHERE id = $1
	`, id)

	var p models.Produto
	if err := row.Scan(&p.ID, &p.Nome, &p.Unidade, &p.CriadoEm); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("buscar produto: %w", err)
	}
	return &p, nil
}

package repositories

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/desafionimo/backend/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PrecoFilters define filtros opcionais para listagem de preços.
type PrecoFilters struct {
	FornecedorID *int64
	ProdutoID    *int64
	DataInicio   *time.Time
	DataFim      *time.Time
	Page         int
	Limit        int
}

// PrecoRepository persiste e consulta registros de preço.
type PrecoRepository struct {
	pool *pgxpool.Pool
}

func NewPrecoRepository(pool *pgxpool.Pool) *PrecoRepository {
	return &PrecoRepository{pool: pool}
}

func (r *PrecoRepository) Upsert(ctx context.Context, produtoID, fornecedorID int64, preco float64, dataRef time.Time) (*models.Preco, error) {
	row := r.pool.QueryRow(ctx, `
		INSERT INTO precos (produto_id, fornecedor_id, preco_por_litro, data_referencia)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (fornecedor_id, produto_id, data_referencia)
		DO UPDATE SET preco_por_litro = EXCLUDED.preco_por_litro
		RETURNING id, produto_id, fornecedor_id, preco_por_litro, data_referencia, criado_em
	`, produtoID, fornecedorID, preco, dataRef)

	var p models.Preco
	if err := row.Scan(&p.ID, &p.ProdutoID, &p.FornecedorID, &p.PrecoPorLitro, &p.DataReferencia, &p.CriadoEm); err != nil {
		return nil, fmt.Errorf("upsert preço: %w", err)
	}
	return &p, nil
}

func (r *PrecoRepository) List(ctx context.Context, f PrecoFilters) (models.PrecoListResult, error) {
	where := []string{"1=1"}
	args := []any{}
	argPos := 1

	if f.FornecedorID != nil {
		where = append(where, fmt.Sprintf("p.fornecedor_id = $%d", argPos))
		args = append(args, *f.FornecedorID)
		argPos++
	}
	if f.ProdutoID != nil {
		where = append(where, fmt.Sprintf("p.produto_id = $%d", argPos))
		args = append(args, *f.ProdutoID)
		argPos++
	}
	if f.DataInicio != nil {
		where = append(where, fmt.Sprintf("p.data_referencia >= $%d", argPos))
		args = append(args, *f.DataInicio)
		argPos++
	}
	if f.DataFim != nil {
		where = append(where, fmt.Sprintf("p.data_referencia <= $%d", argPos))
		args = append(args, *f.DataFim)
		argPos++
	}

	whereSQL := strings.Join(where, " AND ")

	countSQL := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM precos p
		WHERE %s
	`, whereSQL)

	var total int64
	if err := r.pool.QueryRow(ctx, countSQL, args...).Scan(&total); err != nil {
		return models.PrecoListResult{}, fmt.Errorf("contar preços: %w", err)
	}

	offset := (f.Page - 1) * f.Limit
	listArgs := append([]any{}, args...)
	listArgs = append(listArgs, f.Limit, offset)

	listSQL := fmt.Sprintf(`
		SELECT
			p.id,
			p.produto_id,
			p.fornecedor_id,
			p.preco_por_litro::float8,
			p.data_referencia,
			p.criado_em,
			pr.nome,
			fo.nome
		FROM precos p
		INNER JOIN produtos pr ON pr.id = p.produto_id
		INNER JOIN fornecedores fo ON fo.id = p.fornecedor_id
		WHERE %s
		ORDER BY p.data_referencia DESC, p.id DESC
		LIMIT $%d OFFSET $%d
	`, whereSQL, argPos, argPos+1)

	rows, err := r.pool.Query(ctx, listSQL, listArgs...)
	if err != nil {
		return models.PrecoListResult{}, fmt.Errorf("listar preços: %w", err)
	}
	defer rows.Close()

	var data []models.Preco
	for rows.Next() {
		var p models.Preco
		if err := rows.Scan(
			&p.ID,
			&p.ProdutoID,
			&p.FornecedorID,
			&p.PrecoPorLitro,
			&p.DataReferencia,
			&p.CriadoEm,
			&p.ProdutoNome,
			&p.FornecedorNome,
		); err != nil {
			return models.PrecoListResult{}, fmt.Errorf("scan preço: %w", err)
		}
		data = append(data, p)
	}
	if err := rows.Err(); err != nil {
		return models.PrecoListResult{}, fmt.Errorf("iter preços: %w", err)
	}

	return models.PrecoListResult{
		Data:  data,
		Page:  f.Page,
		Limit: f.Limit,
		Total: total,
	}, nil
}

func (r *PrecoRepository) HistoricoPorProduto(ctx context.Context, produtoID int64) ([]models.HistoricoPonto, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT
			p.data_referencia,
			p.fornecedor_id,
			fo.nome,
			p.preco_por_litro::float8
		FROM precos p
		INNER JOIN fornecedores fo ON fo.id = p.fornecedor_id
		WHERE p.produto_id = $1
		ORDER BY p.data_referencia ASC, fo.nome ASC
	`, produtoID)
	if err != nil {
		return nil, fmt.Errorf("histórico por produto: %w", err)
	}
	defer rows.Close()

	var out []models.HistoricoPonto
	for rows.Next() {
		var h models.HistoricoPonto
		if err := rows.Scan(&h.DataReferencia, &h.FornecedorID, &h.FornecedorNome, &h.PrecoPorLitro); err != nil {
			return nil, fmt.Errorf("scan histórico: %w", err)
		}
		out = append(out, h)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iter histórico: %w", err)
	}
	return out, nil
}

func (r *PrecoRepository) ComparativoAtual(ctx context.Context, produtoID int64) ([]models.ComparativoItem, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT DISTINCT ON (p.fornecedor_id)
			p.fornecedor_id,
			fo.nome,
			p.preco_por_litro::float8,
			p.data_referencia
		FROM precos p
		INNER JOIN fornecedores fo ON fo.id = p.fornecedor_id
		WHERE p.produto_id = $1
		ORDER BY p.fornecedor_id, p.data_referencia DESC
	`, produtoID)
	if err != nil {
		return nil, fmt.Errorf("comparativo atual: %w", err)
	}
	defer rows.Close()

	var out []models.ComparativoItem
	for rows.Next() {
		var it models.ComparativoItem
		if err := rows.Scan(&it.FornecedorID, &it.FornecedorNome, &it.PrecoPorLitro, &it.DataReferencia); err != nil {
			return nil, fmt.Errorf("scan comparativo: %w", err)
		}
		out = append(out, it)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iter comparativo: %w", err)
	}
	return out, nil
}


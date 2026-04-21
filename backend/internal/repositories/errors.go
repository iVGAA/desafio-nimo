package repositories

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	// ErrDuplicateProdutoNome indica violação da constraint única de nome do produto.
	ErrDuplicateProdutoNome = errors.New("Produto com este nome já existe.")
	// ErrDuplicateFornecedorNome indica violação da constraint única de nome do fornecedor.
	ErrDuplicateFornecedorNome = errors.New("Fornecedor com este nome já existe.")
	// ErrDuplicateCNPJ indica violação da constraint única de CNPJ.
	ErrDuplicateCNPJ = errors.New("Fornecedor com este CNPJ já existe.")
)

// PostgresCode retorna o código SQLSTATE do erro Postgres, ou string vazia.
func PostgresCode(err error) string {
	var pe *pgconn.PgError
	if errors.As(err, &pe) {
		return pe.Code
	}
	return ""
}

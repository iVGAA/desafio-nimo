package models

import "time"

// Fornecedor representa uma empresa que pratica preços.
type Fornecedor struct {
	ID       int64     `json:"id"`
	Nome     string    `json:"nome"`
	CNPJ     string    `json:"cnpj"`
	CriadoEm time.Time `json:"criado_em"`
}

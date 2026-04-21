package models

import "time"

// Produto representa um tipo de combustível cadastrado.
type Produto struct {
	ID       int64     `json:"id"`
	Nome     string    `json:"nome"`
	Unidade  string    `json:"unidade"`
	CriadoEm time.Time `json:"criado_em"`
}

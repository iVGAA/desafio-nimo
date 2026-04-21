package models

import "time"

// Preco é um registro de preço por litro em uma data de referência.
type Preco struct {
	ID             int64     `json:"id"`
	ProdutoID      int64     `json:"produto_id"`
	FornecedorID   int64     `json:"fornecedor_id"`
	PrecoPorLitro  float64   `json:"preco_por_litro"`
	DataReferencia time.Time `json:"data_referencia"`
	CriadoEm       time.Time `json:"criado_em"`
	ProdutoNome    string    `json:"produto_nome,omitempty"`
	FornecedorNome string    `json:"fornecedor_nome,omitempty"`
}

// HistoricoPonto é um ponto da série histórica de preços de um produto.
type HistoricoPonto struct {
	DataReferencia time.Time `json:"data_referencia"`
	FornecedorID   int64     `json:"fornecedor_id"`
	FornecedorNome string    `json:"fornecedor_nome"`
	PrecoPorLitro  float64   `json:"preco_por_litro"`
}

// HistoricoProdutoResponse agrupa o histórico por produto.
type HistoricoProdutoResponse struct {
	ProdutoID int64            `json:"produto_id"`
	Pontos    []HistoricoPonto `json:"pontos"`
}

// ComparativoItem é o preço mais recente de um fornecedor para o produto.
type ComparativoItem struct {
	FornecedorID   int64     `json:"fornecedor_id"`
	FornecedorNome string    `json:"fornecedor_nome"`
	PrecoPorLitro  float64   `json:"preco_por_litro"`
	DataReferencia time.Time `json:"data_referencia"`
}

// ComparativoResponse compara fornecedores para um mesmo produto.
type ComparativoResponse struct {
	ProdutoID              int64             `json:"produto_id"`
	Itens                  []ComparativoItem `json:"itens"`
	MenorPreco             *float64          `json:"menor_preco,omitempty"`
	MenorPrecoFornecedorID *int64            `json:"menor_preco_fornecedor_id,omitempty"`
}

// PrecoListResult devolve página de preços e total para paginação.
type PrecoListResult struct {
	Data  []Preco `json:"data"`
	Page  int     `json:"page"`
	Limit int     `json:"limit"`
	Total int64   `json:"total"`
}

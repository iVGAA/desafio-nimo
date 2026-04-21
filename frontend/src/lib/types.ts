export interface Produto {
  id: number;
  nome: string;
  unidade: string;
  criado_em: string;
}

export interface Fornecedor {
  id: number;
  nome: string;
  cnpj: string;
  criado_em: string;
}

export interface CreateProdutoPayload {
  nome: string;
  unidade: string;
}

export interface CreateFornecedorPayload {
  nome: string;
  cnpj: string;
}

export interface HistoricoPonto {
  data_referencia: string;
  fornecedor_id: number;
  fornecedor_nome: string;
  preco_por_litro: number;
}

export interface HistoricoResponse {
  produto_id: number;
  pontos: HistoricoPonto[];
}

export interface ComparativoItem {
  fornecedor_id: number;
  fornecedor_nome: string;
  preco_por_litro: number;
  data_referencia: string;
}

export interface ComparativoResponse {
  produto_id: number;
  itens: ComparativoItem[];
  menor_preco?: number;
  menor_preco_fornecedor_id?: number;
}

export interface CreatePrecoPayload {
  produto_id: number;
  fornecedor_id: number;
  preco_por_litro: number;
  data_referencia: string;
}

export interface Preco {
  id: number;
  produto_id: number;
  fornecedor_id: number;
  preco_por_litro: number;
  data_referencia: string;
  criado_em: string;
}

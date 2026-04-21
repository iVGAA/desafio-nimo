import type {
  ComparativoResponse,
  CreateFornecedorPayload,
  CreatePrecoPayload,
  CreateProdutoPayload,
  Fornecedor,
  HistoricoResponse,
  Preco,
  Produto
} from "@/lib/types";

const isServer = typeof window === "undefined";

const API_URL = isServer
  ? process.env.NEXT_PUBLIC_API_URL || "http://backend:8080"
  : "http://localhost:8080";

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`${API_URL}${path}`, {
    ...init,
    headers: {
      "Content-Type": "application/json",
      ...init?.headers
    },
    cache: init?.cache ?? "no-store"
  });

  if (!response.ok) {
    let message = `Erro ${response.status}`;
    try {
      const body = (await response.json()) as { error?: string };
      if (body.error) {
        message = body.error;
      }
    } catch {
      // mantém mensagem padrão quando corpo não é JSON.
    }
    throw new Error(message);
  }

  return (await response.json()) as T;
}

export function listProdutos(): Promise<Produto[]> {
  return request<Produto[]>("/api/v1/produtos");
}

export function listFornecedores(): Promise<Fornecedor[]> {
  return request<Fornecedor[]>("/api/v1/fornecedores");
}

export function createProduto(payload: CreateProdutoPayload): Promise<Produto> {
  return request<Produto>("/api/v1/produtos", {
    method: "POST",
    body: JSON.stringify(payload)
  });
}

export function createFornecedor(payload: CreateFornecedorPayload): Promise<Fornecedor> {
  return request<Fornecedor>("/api/v1/fornecedores", {
    method: "POST",
    body: JSON.stringify(payload)
  });
}

export function getHistoricoProduto(produtoId: number): Promise<HistoricoResponse> {
  return request<HistoricoResponse>(`/api/v1/precos/historico/${produtoId}`);
}

export function getComparativoProduto(
  produtoId: number
): Promise<ComparativoResponse> {
  return request<ComparativoResponse>(
    `/api/v1/precos/comparativo?produto_id=${produtoId}`
  );
}

export function createPreco(payload: CreatePrecoPayload): Promise<Preco> {
  return request<Preco>("/api/v1/precos", {
    method: "POST",
    body: JSON.stringify(payload)
  });
}

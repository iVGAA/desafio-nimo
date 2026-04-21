"use client";

import { FormEvent, useEffect, useState, CSSProperties } from "react";
import {
  createFornecedor,
  createProduto,
  listFornecedores,
  listProdutos
} from "@/lib/api";
import { formatCNPJ } from "@/lib/formatters";
import type { Fornecedor, Produto } from "@/lib/types";

type ProdutoForm = {
  nome: string;
  unidade: string;
};

type FornecedorForm = {
  nome: string;
  cnpj: string;
};

const initialProdutoForm: ProdutoForm = {
  nome: "",
  unidade: "litro"
};

const initialFornecedorForm: FornecedorForm = {
  nome: "",
  cnpj: ""
};

const inputStyle: CSSProperties = {
  width: "100%",
  padding: 10,
  borderRadius: 8,
  border: "1px solid #d1d5db"
};

const getBtnStyle = (disabled: boolean): CSSProperties => ({
  padding: "10px 14px",
  borderRadius: 8,
  border: "1px solid #2563eb",
  background: "#2563eb",
  color: "#fff",
  fontWeight: 600,
  cursor: disabled ? "not-allowed" : "pointer",
  opacity: disabled ? 0.6 : 1
});

export default function GestaoPage() {
  const [produtos, setProdutos] = useState<Produto[]>([]);
  const [fornecedores, setFornecedores] = useState<Fornecedor[]>([]);

  const [produtoForm, setProdutoForm] = useState<ProdutoForm>(initialProdutoForm);
  const [fornecedorForm, setFornecedorForm] = useState<FornecedorForm>(initialFornecedorForm);

  const [loading, setLoading] = useState<boolean>(true);
  const [loadingProduto, setLoadingProduto] = useState<boolean>(false);
  const [loadingFornecedor, setLoadingFornecedor] = useState<boolean>(false);

  const [produtoError, setProdutoError] = useState<string | null>(null);
  const [produtoSuccess, setProdutoSuccess] = useState<string | null>(null);
  const [fornecedorError, setFornecedorError] = useState<string | null>(null);
  const [fornecedorSuccess, setFornecedorSuccess] = useState<string | null>(null);
  const [pageError, setPageError] = useState<string | null>(null);

  useEffect(() => {
    let active = true;

    async function loadData() {
      setLoading(true);
      setPageError(null);
      try {
        const [produtosData, fornecedoresData] = await Promise.all([
          listProdutos(),
          listFornecedores()
        ]);
        if (!active) return;
        setProdutos(produtosData);
        setFornecedores(fornecedoresData);
      } catch (err) {
        if (!active) return;
        setPageError(err instanceof Error ? err.message : "Erro ao carregar dados da gestão");
      } finally {
        if (active) setLoading(false);
      }
    }

    void loadData();
    return () => {
      active = false;
    };
  }, []);

  async function refreshProdutos() {
    const data = await listProdutos();
    setProdutos(data);
  }

  async function refreshFornecedores() {
    const data = await listFornecedores();
    setFornecedores(data);
  }

  function validateProdutoForm(): string | null {
    if (!produtoForm.nome.trim()) return "Nome do produto é obrigatório.";
    if (!produtoForm.unidade.trim()) return "Unidade é obrigatória.";
    return null;
  }

  function validateFornecedorForm(): string | null {
    if (!fornecedorForm.nome.trim()) return "Nome do fornecedor é obrigatório.";
    if (!fornecedorForm.cnpj.trim()) return "CNPJ é obrigatório.";
    return null;
  }

  async function handleSubmitProduto(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setProdutoError(null);
    setProdutoSuccess(null);

    const validationError = validateProdutoForm();
    if (validationError) {
      setProdutoError(validationError);
      return;
    }

    setLoadingProduto(true);
    try {
      await createProduto({
        nome: produtoForm.nome.trim(),
        unidade: produtoForm.unidade.trim()
      });
      await refreshProdutos();
      setProdutoSuccess("Produto cadastrado com sucesso.");
      setProdutoForm(initialProdutoForm);
    } catch (err) {
      setProdutoError(err instanceof Error ? err.message : "Erro ao cadastrar produto");
    } finally {
      setLoadingProduto(false);
    }
  }

  async function handleSubmitFornecedor(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setFornecedorError(null);
    setFornecedorSuccess(null);

    const validationError = validateFornecedorForm();
    if (validationError) {
      setFornecedorError(validationError);
      return;
    }

    setLoadingFornecedor(true);
    try {
      await createFornecedor({
        nome: fornecedorForm.nome.trim(),
        cnpj: fornecedorForm.cnpj.trim()
      });
      await refreshFornecedores();
      setFornecedorSuccess("Fornecedor cadastrado com sucesso.");
      setFornecedorForm(initialFornecedorForm);
    } catch (err) {
      setFornecedorError(err instanceof Error ? err.message : "Erro ao cadastrar fornecedor");
    } finally {
      setLoadingFornecedor(false);
    }
  }

  return (
    <main>
      <h1>Gestão</h1>
      <p className="helper">Cadastre produtos e fornecedores usados no monitoramento.</p>

      {pageError && <p className="error" style={{ marginTop: 16 }}>{pageError}</p>}
      {loading && <p className="helper" style={{ marginTop: 16 }}>Carregando...</p>}

      <section className="grid grid-2" style={{ marginTop: 16 }}>
        <article className="card">
          <h2 className="title">Cadastrar Produto</h2>
          <form className="grid" onSubmit={handleSubmitProduto}>
            <div>
              <label htmlFor="produtoNome">Nome</label>
              <input
                id="produtoNome"
                type="text"
                placeholder="Ex.: Diesel S10"
                value={produtoForm.nome}
                onChange={(event) =>
                  setProdutoForm((current) => ({ ...current, nome: event.target.value }))
                }
                disabled={loadingProduto}
                style={inputStyle}
              />
            </div>
            <div>
              <label htmlFor="produtoUnidade">Unidade</label>
              <input
                id="produtoUnidade"
                type="text"
                placeholder="litro"
                value={produtoForm.unidade}
                onChange={(event) =>
                  setProdutoForm((current) => ({ ...current, unidade: event.target.value }))
                }
                disabled={loadingProduto}
                style={inputStyle}
              />
            </div>
            <button
              type="submit"
              disabled={loadingProduto}
              style={getBtnStyle(loadingProduto)}
            >
              {loadingProduto ? "Salvando..." : "Salvar produto"}
            </button>
          </form>
          {produtoError && <p className="error" style={{ marginTop: 12 }}>{produtoError}</p>}
          {produtoSuccess && (
            <p style={{ marginTop: 12, color: "#166534", fontWeight: 600 }}>{produtoSuccess}</p>
          )}

          <div style={{ marginTop: 12 }}>
            <p className="helper">Produtos cadastrados: {produtos.length}</p>
          </div>
        </article>

        <article className="card">
          <h2 className="title">Cadastrar Fornecedor</h2>
          <form className="grid" onSubmit={handleSubmitFornecedor}>
            <div>
              <label htmlFor="fornecedorNome">Nome</label>
              <input
                id="fornecedorNome"
                type="text"
                placeholder="Ex.: Distribuidora XPTO"
                value={fornecedorForm.nome}
                onChange={(event) =>
                  setFornecedorForm((current) => ({ ...current, nome: event.target.value }))
                }
                disabled={loadingFornecedor}
                style={inputStyle}
              />
            </div>
            <div>
              <label htmlFor="fornecedorCnpj">CNPJ</label>
              <input
                id="fornecedorCnpj"
                type="text"
                placeholder="Ex.: 12.345.678/0001-90"
                value={formatCNPJ(fornecedorForm.cnpj)}
                onChange={(event) => {
                  const rawValue = event.target.value.replace(/\D/g, "").slice(0, 14);
                  setFornecedorForm((current) => ({ ...current, cnpj: rawValue }));
                }}
                disabled={loadingFornecedor}
                style={inputStyle}
              />
            </div>
            <button
              type="submit"
              disabled={loadingFornecedor}
              style={getBtnStyle(loadingFornecedor)}
            >
              {loadingFornecedor ? "Salvando..." : "Salvar fornecedor"}
            </button>
          </form>
          {fornecedorError && (
            <p className="error" style={{ marginTop: 12 }}>{fornecedorError}</p>
          )}
          {fornecedorSuccess && (
            <p style={{ marginTop: 12, color: "#166534", fontWeight: 600 }}>
              {fornecedorSuccess}
            </p>
          )}

          <div style={{ marginTop: 12 }}>
            <p className="helper">Fornecedores cadastrados: {fornecedores.length}</p>
          </div>
        </article>
      </section>
    </main>
  );
}

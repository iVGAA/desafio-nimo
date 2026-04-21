"use client";

import { FormEvent, useEffect, useMemo, useState } from "react";
import { createPreco, listFornecedores, listProdutos } from "@/lib/api";
import { getTodayISODateLocal } from "@/lib/date";
import type { Fornecedor, Produto } from "@/lib/types";

type FormState = {
  produtoId: string;
  fornecedorId: string;
  dataReferencia: string;
  precoPorLitro: string;
};

const initialFormState: FormState = {
  produtoId: "",
  fornecedorId: "",
  dataReferencia: "",
  precoPorLitro: ""
};

function isDateValid(value: string): boolean {
  if (!/^\d{4}-\d{2}-\d{2}$/.test(value)) {
    return false;
  }
  const parsed = new Date(`${value}T00:00:00Z`);
  return !Number.isNaN(parsed.getTime());
}

export default function NovoPrecoPage() {
  const [produtos, setProdutos] = useState<Produto[]>([]);
  const [fornecedores, setFornecedores] = useState<Fornecedor[]>([]);
  const [form, setForm] = useState<FormState>(initialFormState);
  const [loadingInitial, setLoadingInitial] = useState<boolean>(true);
  const [submitting, setSubmitting] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const today = useMemo(() => getTodayISODateLocal(), []);

  useEffect(() => {
    let active = true;

    async function loadData() {
      setLoadingInitial(true);
      setError(null);
      try {
        const [produtosData, fornecedoresData] = await Promise.all([
          listProdutos(),
          listFornecedores()
        ]);
        if (!active) return;

        setProdutos(produtosData);
        setFornecedores(fornecedoresData);
        setForm((current) => ({
          ...current,
          produtoId: produtosData[0]?.id ? String(produtosData[0].id) : "",
          fornecedorId: fornecedoresData[0]?.id ? String(fornecedoresData[0].id) : "",
          dataReferencia: getTodayISODateLocal()
        }));
      } catch (err) {
        if (!active) return;
        setError(err instanceof Error ? err.message : "Erro ao carregar dados do formulário");
      } finally {
        if (active) {
          setLoadingInitial(false);
        }
      }
    }

    void loadData();
    return () => {
      active = false;
    };
  }, []);

  const isFormDisabled = useMemo(
    () => loadingInitial || submitting || produtos.length === 0 || fornecedores.length === 0,
    [loadingInitial, submitting, produtos.length, fornecedores.length]
  );

  function validate(): string | null {
    if (!form.produtoId) return "Selecione um produto.";
    if (!form.fornecedorId) return "Selecione um fornecedor.";
    if (!form.dataReferencia) return "Informe a data de referência.";
    if (!isDateValid(form.dataReferencia)) return "Data inválida. Use o formato YYYY-MM-DD.";
    if (form.dataReferencia > today) return "A data de referência não pode ser futura.";
    if (!form.precoPorLitro) return "Informe o valor por litro.";

    const preco = Number(form.precoPorLitro);
    if (Number.isNaN(preco) || preco <= 0) {
      return "O valor deve ser um número maior que zero.";
    }
    return null;
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setError(null);
    setSuccess(null);

    const validationError = validate();
    if (validationError) {
      setError(validationError);
      return;
    }

    setSubmitting(true);
    try {
      await createPreco({
        produto_id: Number(form.produtoId),
        fornecedor_id: Number(form.fornecedorId),
        preco_por_litro: Number(form.precoPorLitro),
        data_referencia: form.dataReferencia
      });
      setSuccess("Preço cadastrado com sucesso. Se já existia para o mesmo dia, foi atualizado.");
      setForm((current) => ({
        ...current,
        precoPorLitro: ""
      }));
    } catch (err) {
      setError(err instanceof Error ? err.message : "Erro ao cadastrar preço");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <main>
      <h1>Cadastro de Preço</h1>
      <p className="helper">Registre o preço por fornecedor, produto e data de referência.</p>

      <section className="card" style={{ marginTop: 16, maxWidth: 560 }}>
        <form onSubmit={handleSubmit} className="grid">
          <div>
            <label htmlFor="produto">Produto</label>
            <select
              id="produto"
              value={form.produtoId}
              onChange={(event) => setForm((current) => ({ ...current, produtoId: event.target.value }))}
              disabled={isFormDisabled}
            >
              {produtos.length === 0 && <option value="">Nenhum produto cadastrado</option>}
              {produtos.map((produto) => (
                <option key={produto.id} value={produto.id}>
                  {produto.nome}
                </option>
              ))}
            </select>
          </div>

          <div>
            <label htmlFor="fornecedor">Fornecedor</label>
            <select
              id="fornecedor"
              value={form.fornecedorId}
              onChange={(event) => setForm((current) => ({ ...current, fornecedorId: event.target.value }))}
              disabled={isFormDisabled}
            >
              {fornecedores.length === 0 && <option value="">Nenhum fornecedor cadastrado</option>}
              {fornecedores.map((fornecedor) => (
                <option key={fornecedor.id} value={fornecedor.id}>
                  {fornecedor.nome}
                </option>
              ))}
            </select>
          </div>

          <div>
            <label htmlFor="dataReferencia">Data de referência</label>
            <input
              id="dataReferencia"
              type="date"
              value={form.dataReferencia}
              max={today}
              onChange={(event) =>
                setForm((current) => ({ ...current, dataReferencia: event.target.value }))
              }
              disabled={isFormDisabled}
              style={{ width: "100%", padding: 10, borderRadius: 8, border: "1px solid #d1d5db" }}
            />
          </div>

          <div>
            <label htmlFor="precoPorLitro">Valor por litro (R$)</label>
            <input
              id="precoPorLitro"
              type="number"
              min="0"
              step="0.001"
              inputMode="decimal"
              placeholder="Ex.: 5.799"
              value={form.precoPorLitro}
              onChange={(event) =>
                setForm((current) => ({ ...current, precoPorLitro: event.target.value }))
              }
              disabled={isFormDisabled}
              style={{ width: "100%", padding: 10, borderRadius: 8, border: "1px solid #d1d5db" }}
            />
          </div>

          <button
            type="submit"
            disabled={isFormDisabled}
            style={{
              padding: "10px 14px",
              borderRadius: 8,
              border: "1px solid #2563eb",
              background: "#2563eb",
              color: "#fff",
              fontWeight: 600,
              cursor: isFormDisabled ? "not-allowed" : "pointer",
              opacity: isFormDisabled ? 0.6 : 1
            }}
          >
            {submitting ? "Salvando..." : "Salvar preço"}
          </button>
        </form>
      </section>

      {error && <p className="error" style={{ marginTop: 16 }}>{error}</p>}
      {success && <p style={{ marginTop: 16, color: "#166534", fontWeight: 600 }}>{success}</p>}
    </main>
  );
}

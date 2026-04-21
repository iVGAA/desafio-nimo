"use client";

import { useEffect, useMemo, useState } from "react";
import CurrentPricesList from "@/components/CurrentPricesList";
import PriceHistoryChart from "@/components/PriceHistoryChart";
import {
  getComparativoProduto,
  getHistoricoProduto,
  listProdutos
} from "@/lib/api";
import type {
  ComparativoResponse,
  HistoricoResponse,
  Produto
} from "@/lib/types";

type DashboardData = {
  historico: HistoricoResponse;
  comparativo: ComparativoResponse;
};

export default function DashboardPage() {
  const [produtos, setProdutos] = useState<Produto[]>([]);
  const [selectedProdutoId, setSelectedProdutoId] = useState<number | null>(null);
  const [dashboardData, setDashboardData] = useState<DashboardData | null>(null);
  
  const [isLoadingProdutos, setIsLoadingProdutos] = useState<boolean>(true);
  const [isFetchingDashboard, setIsFetchingDashboard] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let isMounted = true;

    async function fetchProdutos() {
      setIsLoadingProdutos(true);
      setError(null);
      try {
        const data = await listProdutos();
        if (!isMounted) return;
        setProdutos(data);
        if (data.length > 0) {
          setSelectedProdutoId(data[0].id);
        }
      } catch (err) {
        if (!isMounted) return;
        setError(err instanceof Error ? err.message : "Erro ao carregar produtos");
      } finally {
        if (isMounted) setIsLoadingProdutos(false);
      }
    }

    void fetchProdutos();
    return () => {
      isMounted = false;
    };
  }, []);

  useEffect(() => {
    if (!selectedProdutoId) {
      setDashboardData(null);
      return;
    }

    let isMounted = true;
    setIsFetchingDashboard(true);
    setError(null);

    async function fetchDashboardData() {
      try {
        const [historico, comparativo] = await Promise.all([
          getHistoricoProduto(selectedProdutoId!),
          getComparativoProduto(selectedProdutoId!)
        ]);
        if (!isMounted) return;
        setDashboardData({ historico, comparativo });
      } catch (err) {
        if (!isMounted) return;
        setError(err instanceof Error ? err.message : "Erro ao carregar dados do dashboard");
      } finally {
        if (isMounted) setIsFetchingDashboard(false);
      }
    }

    void fetchDashboardData();
    return () => {
      isMounted = false;
    };
  }, [selectedProdutoId]);

  const selectedProduto = useMemo(
    () => produtos.find((item) => item.id === selectedProdutoId),
    [produtos, selectedProdutoId]
  );

  return (
    <main>
      <h1>Dashboard de Preços de Combustível</h1>
      <p className="helper">Selecione um produto para visualizar histórico e comparativo atual.</p>

      <section className="card" style={{ marginTop: 16 }}>
        <label htmlFor="produto">Produto</label>
        <select
          id="produto"
          value={selectedProdutoId ?? ""}
          onChange={(event) => setSelectedProdutoId(Number(event.target.value))}
          disabled={produtos.length === 0 || isLoadingProdutos}
        >
          {produtos.length === 0 && <option value="">Nenhum produto encontrado</option>}
          {produtos.map((produto) => (
            <option key={produto.id} value={produto.id}>
              {produto.nome}
            </option>
          ))}
        </select>
      </section>

      {isLoadingProdutos && <p className="helper" style={{ marginTop: 16 }}>Carregando produtos...</p>}
      {error && <p className="error" style={{ marginTop: 16 }}>{error}</p>}

      {!isLoadingProdutos && !error && selectedProduto && (
        <section 
          className="grid grid-2" 
          style={{ 
            marginTop: 16, 
            opacity: isFetchingDashboard ? 0.5 : 1, 
            transition: "opacity 0.2s" 
          }}
        >
          <article className="card">
            <h2 className="title">Histórico de preços - {selectedProduto.nome}</h2>
            <PriceHistoryChart 
              points={dashboardData?.historico.pontos ?? []} 
              isLoading={!dashboardData && isFetchingDashboard} 
            />
          </article>

          <article className="card">
            <h2 className="title">Preço atual por fornecedor</h2>
            <CurrentPricesList
              items={dashboardData?.comparativo.itens ?? []}
              cheapestSupplierId={dashboardData?.comparativo.menor_preco_fornecedor_id}
            />
          </article>
        </section>
      )}
    </main>
  );
}

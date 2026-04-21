"use client";

import {
  CartesianGrid,
  Legend,
  Line,
  LineChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis
} from "recharts";
import type { CSSProperties } from "react";
import { extractDateOnly, formatDateBRFromISODate } from "@/lib/date";
import type { HistoricoPonto } from "@/lib/types";

type PriceHistoryChartProps = {
  points: HistoricoPonto[];
  isLoading?: boolean;
};

type ChartRow = {
  data: string;
  [fornecedor: string]: number | string;
};

const containerStyle: CSSProperties = { width: "100%", height: 360 };
const centerStyle: CSSProperties = {
  ...containerStyle,
  display: "flex",
  alignItems: "center",
  justifyContent: "center"
};

function formatDate(value: string): string {
  return formatDateBRFromISODate(value);
}

function groupPointsByDate(points: HistoricoPonto[]): ChartRow[] {
  const byDate = new Map<string, ChartRow>();

  for (const point of points) {
    const key = extractDateOnly(point.data_referencia);
    const current = byDate.get(key) ?? { data: key };
    current[point.fornecedor_nome] = point.preco_por_litro;
    byDate.set(key, current);
  }

  return Array.from(byDate.values()).sort((a, b) => a.data.localeCompare(b.data));
}

function getUniqueSuppliers(points: HistoricoPonto[]): string[] {
  return Array.from(new Set(points.map((item) => item.fornecedor_nome)));
}

const COLORS = ["#2563eb", "#16a34a", "#ea580c", "#9333ea", "#0891b2", "#dc2626"];

export default function PriceHistoryChart({ points, isLoading }: PriceHistoryChartProps) {
  if (isLoading) {
    return (
      <div style={centerStyle}>
        <span>Carregando...</span>
      </div>
    );
  }

  if (!points || points.length === 0) {
    return (
      <div style={centerStyle}>
        <span>Nenhum dado disponível</span>
      </div>
    );
  }

  const rows = groupPointsByDate(points);
  const supplierNames = getUniqueSuppliers(points);

  return (
    <div style={containerStyle}>
      <ResponsiveContainer width="100%" height="100%">
        <LineChart data={rows}>
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis dataKey="data" tickFormatter={formatDate} />
          <YAxis />
          <Tooltip
            formatter={(value) => `R$ ${Number(value).toFixed(3)}`}
            labelFormatter={(value) => `Data: ${formatDate(String(value))}`}
          />
          <Legend />
          {supplierNames.map((name, index) => (
            <Line
              key={name}
              type="monotone"
              dataKey={name}
              stroke={COLORS[index % COLORS.length]}
              strokeWidth={2}
              dot={{ r: 3 }}
              activeDot={{ r: 5 }}
            />
          ))}
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
}

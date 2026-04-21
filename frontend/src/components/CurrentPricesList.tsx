"use client";

import type { ComparativoItem } from "@/lib/types";
import DateText from "@/components/DateText";

type CurrentPricesListProps = {
  items: ComparativoItem[];
  cheapestSupplierId?: number;
};

export default function CurrentPricesList({
  items,
  cheapestSupplierId
}: CurrentPricesListProps) {
  if (items.length === 0) {
    return <p className="helper">Sem preços atuais para este produto.</p>;
  }

  return (
    <div className="list">
      {items.map((item) => {
        const isBest = item.fornecedor_id === cheapestSupplierId;
        return (
          <div key={item.fornecedor_id} className={`list-item ${isBest ? "best" : ""}`}>
            <div>
              <strong>{item.fornecedor_nome}</strong>
              <div className="helper">
                Ref.: <DateText value={item.data_referencia} />
              </div>
            </div>
            <strong>R$ {item.preco_por_litro.toFixed(3)}</strong>
          </div>
        );
      })}
    </div>
  );
}

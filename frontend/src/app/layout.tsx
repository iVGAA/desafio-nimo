import "./globals.css";
import type { Metadata } from "next";
import type { ReactNode } from "react";
import TopNav from "@/components/TopNav";

export const metadata: Metadata = {
  title: "Monitor de Preços de Combustível",
  description: "Dashboard de histórico e comparativo de preços"
};

type RootLayoutProps = {
  children: ReactNode;
};

export default function RootLayout({ children }: RootLayoutProps) {
  return (
    <html lang="pt-BR">
      <body>
        <TopNav />
        {children}
      </body>
    </html>
  );
}

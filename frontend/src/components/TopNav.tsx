"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

const NAV_ITEMS = [
  { href: "/", label: "Dashboard" },
  { href: "/precos/novo", label: "Cadastro de Preço" },
  { href: "/gestao", label: "Gestão" }
];

export default function TopNav() {
  const pathname = usePathname();

  return (
    <header className="topnav-wrap">
      <nav className="topnav" aria-label="Navegação principal">
        {NAV_ITEMS.map((item) => {
          const isActive = pathname === item.href;
          return (
            <Link
              key={item.href}
              href={item.href}
              className={`topnav-link ${isActive ? "topnav-link-active" : ""}`}
            >
              {item.label}
            </Link>
          );
        })}
      </nav>
    </header>
  );
}

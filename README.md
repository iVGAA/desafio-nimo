# Painel de Monitoramento de Preços de Combustível

Sistema web para registrar e acompanhar a variação de preços de combustíveis por fornecedor, construído com Go (Fiber), Next.js, TypeScript e PostgreSQL.

---

## Pré-requisitos

- [Docker](https://www.docker.com/) e [Docker Compose](https://docs.docker.com/compose/) instalados

---

## Como rodar o projeto

```bash
# 1. Clone o repositório
git clone https://github.com/iVGAA/desafio-nimo.git
cd desafio-nimo

# 2. Suba todos os serviços
docker compose up --build
```

Aguarde os três contêineres inicializarem. O banco sobe primeiro (com healthcheck), depois o backend, depois o frontend.

| Serviço  | URL                        |
|----------|----------------------------|
| Frontend | http://localhost:3000       |
| Backend  | http://localhost:8080       |
| Banco    | localhost:5432              |

> O banco já é inicializado com dados de exemplo (3 produtos, 3 fornecedores e registros de preço dos últimos 6 dias) para que o dashboard fique imediatamente funcional.

Para encerrar:

```bash
docker compose down
```

Para encerrar e apagar o volume do banco (reset completo):

```bash
docker compose down -v
```

---

## Estrutura do projeto

```
desafio-nimo/
├── backend/
│   ├── cmd/server/         # Ponto de entrada (main.go)
│   └── internal/
│       ├── config/         # Carregamento de variáveis de ambiente
│       ├── database/       # Pool de conexões pgx
│       ├── handlers/       # Camada HTTP (Fiber) — parse, validação, resposta
│       ├── models/         # Structs de domínio e resposta
│       ├── repositories/   # Queries SQL diretas ao PostgreSQL
│       ├── routes/         # Registro das rotas em /api/v1
│       └── services/       # Regras de negócio
├── db/
│   ├── schema.sql          # DDL: tabelas, constraints, índices
│   └── init.sql            # Dados iniciais de exemplo
├── frontend/
│   └── src/
│       ├── app/            # Pages (App Router)
│       │   ├── page.tsx            # Dashboard de preços
│       │   ├── precos/novo/        # Cadastro de preço
│       │   └── gestao/             # Cadastro de produtos e fornecedores
│       ├── components/     # Componentes reutilizáveis
│       └── lib/
│           ├── api.ts      # Funções de acesso à API REST
│           ├── types.ts    # Interfaces TypeScript
│           ├── date.ts     # Helpers de data sem timezone
│           └── formatters.ts # Máscara de CNPJ
└── docker-compose.yml
```

---

## Endpoints da API

Base: `http://localhost:8080/api/v1`

| Método | Rota                              | Descrição                                              |
|--------|-----------------------------------|--------------------------------------------------------|
| GET    | `/produtos`                       | Lista todos os produtos                                |
| POST   | `/produtos`                       | Cadastra um novo produto                               |
| GET    | `/fornecedores`                   | Lista todos os fornecedores                            |
| POST   | `/fornecedores`                   | Cadastra um novo fornecedor                            |
| POST   | `/precos`                         | Registra um novo preço                                 |
| GET    | `/precos`                         | Lista preços com filtros e paginação                   |
| GET    | `/precos/historico/:produto_id`   | Série histórica de preços de um produto                |
| GET    | `/precos/comparativo`             | Preço mais recente por fornecedor para um produto      |
| GET    | `/health`                         | Health check do servidor                               |

### Filtros disponíveis em `GET /precos`

| Parâmetro      | Tipo   | Descrição                                     |
|----------------|--------|-----------------------------------------------|
| `produto_id`   | int    | Filtra pelo produto                           |
| `fornecedor_id`| int    | Filtra pelo fornecedor                        |
| `data_inicio`  | string | Data inicial no formato `YYYY-MM-DD`          |
| `data_fim`     | string | Data final no formato `YYYY-MM-DD`            |
| `page`         | int    | Página (padrão: 1)                            |
| `limit`        | int    | Registros por página (padrão: 20, máximo: 100)|

**Exemplo:**
```
GET /api/v1/precos?produto_id=1&data_inicio=2026-04-15&data_fim=2026-04-20&page=1&limit=10
```

---

## Variáveis de ambiente

O backend lê as seguintes variáveis (já configuradas no `docker-compose.yml`):

| Variável             | Descrição                              | Padrão (Docker)                                              |
|----------------------|----------------------------------------|--------------------------------------------------------------|
| `PORT`               | Porta do servidor Go                   | `8080`                                                       |
| `DATABASE_URL`       | String de conexão PostgreSQL           | `postgres://combustivel:combustivel_dev@postgres:5432/...`   |
| `CORS_ALLOW_ORIGINS` | Origens permitidas para CORS           | `http://localhost:3000`                                      |

Consulte `.env.example` no diretório `backend/` para referência ao rodar fora do Docker.

---

## Decisões técnicas

### Filtros do `GET /precos`

O endpoint aceita filtragem por `produto_id`, `fornecedor_id`, `data_inicio` e `data_fim`, além de paginação via `page` e `limit`.

Essa combinação foi escolhida porque reflete os casos de uso reais de uma distribuidora: consultar o histórico de um produto específico, auditar os registros de um fornecedor em um período, ou exportar dados de uma janela de tempo para análise. Todos os filtros são opcionais e combinam entre si via `AND`, o que permite consultas simples ou compostas sem necessidade de endpoints adicionais.

### Preço duplicado no mesmo dia (mesmo fornecedor + produto)

A decisão foi **sobrescrever** o registro existente via `UPSERT` (`ON CONFLICT ... DO UPDATE`).

A constraint única `(fornecedor_id, produto_id, data_referencia)` garante que exista no máximo um preço por combinação fornecedor-produto-dia. Quando um novo preço é registrado para a mesma combinação, o valor é atualizado para o mais recente.

Essa abordagem faz sentido no contexto do sistema: o objetivo é monitorar o preço praticado em cada data, e um segundo registro no mesmo dia representa uma correção ou atualização, não um evento distinto. O frontend exibe uma mensagem informando o comportamento ao usuário após o envio do formulário.

### Unicidade de nomes (case-insensitive)

Produtos e fornecedores têm unicidade garantida via índice funcional `LOWER(nome)` no banco, o que impede cadastros duplicados independentemente da capitalização (ex.: "Diesel S10" e "diesel s10" são tratados como o mesmo produto).

### Tratamento de datas sem timezone

As datas de referência são armazenadas como `DATE` no PostgreSQL e tratadas no frontend como strings `YYYY-MM-DD` puras, sem conversão de timezone. Isso evita que a exibição de uma data mude conforme o fuso horário do navegador do usuário — um problema comum em sistemas que convertem `DATE` para `timestamp` e depois para local time.

### Backend Go — arquitetura em camadas

O backend segue separação explícita em três camadas: `handlers` (HTTP), `services` (regras de negócio) e `repositories` (acesso ao banco). Essa divisão facilita testes unitários por camada, já que cada serviço recebe seus repositórios por injeção de dependência no `main.go`. Erros de banco com semântica de negócio (duplicata, FK inválida) são mapeados para erros sentinela tipados em `repositories/errors.go` e convertidos nos handlers para os status HTTP corretos (409 para conflito, 404 para não encontrado, 400 para validação).

### Banco de dados — `NUMERIC` para preços

O tipo `NUMERIC(14, 4)` foi escolhido para `preco_por_litro` em vez de `FLOAT`, pois aritmética de ponto flutuante introduz erros de precisão em valores monetários. O cast para `float8` acontece apenas no momento da leitura pelo Go, para serialização em JSON.

---

## Bibliotecas utilizadas

**Backend**
- [gofiber/fiber v2](https://github.com/gofiber/fiber) — framework HTTP
- [jackc/pgx v5](https://github.com/jackc/pgx) — driver PostgreSQL

**Frontend**
- [Recharts](https://recharts.org/) — gráfico de linha do histórico de preços
- Next.js App Router com TypeScript

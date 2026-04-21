-- Schema inicial: produtos, fornecedores e preços.
-- Extensão para UUIDs opcionais (não usada nas PKs abaixo; mantida para evolução futura).
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Tipos de combustível (ex.: Diesel S10, Gasolina Comum).
CREATE TABLE produtos (
    id          BIGSERIAL PRIMARY KEY,
    nome        TEXT NOT NULL,
    unidade     TEXT NOT NULL DEFAULT 'litro',
    criado_em   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Empresas que praticam preços.
CREATE TABLE fornecedores (
    id          BIGSERIAL PRIMARY KEY,
    nome        TEXT NOT NULL,
    cnpj        TEXT NOT NULL,
    criado_em   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_fornecedores_cnpj UNIQUE (cnpj)
);

-- Registro de preço por litro em uma data de referência.
-- Um par (fornecedor, produto, data) só pode existir uma vez → UPSERT no backend.
CREATE TABLE precos (
    id                BIGSERIAL PRIMARY KEY,
    produto_id        BIGINT NOT NULL REFERENCES produtos (id) ON DELETE CASCADE,
    fornecedor_id     BIGINT NOT NULL REFERENCES fornecedores (id) ON DELETE CASCADE,
    preco_por_litro   NUMERIC(14, 4) NOT NULL CHECK (preco_por_litro > 0),
    data_referencia   DATE NOT NULL,
    criado_em         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_preco_fornecedor_produto_data UNIQUE (fornecedor_id, produto_id, data_referencia)
);

CREATE INDEX idx_precos_produto_data ON precos (produto_id, data_referencia DESC);
CREATE INDEX idx_precos_fornecedor_data ON precos (fornecedor_id, data_referencia DESC);
CREATE INDEX idx_precos_data ON precos (data_referencia);

-- Unicidade case-insensitive para nomes.
CREATE UNIQUE INDEX uq_produtos_nome_ci ON produtos (LOWER(nome));
CREATE UNIQUE INDEX uq_fornecedores_nome_ci ON fornecedores (LOWER(nome));

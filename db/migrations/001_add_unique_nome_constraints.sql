-- Migration para ambientes já existentes.
-- Garante unicidade de nome em produtos e fornecedores.

ALTER TABLE produtos
    ADD CONSTRAINT uq_produtos_nome UNIQUE (nome);

ALTER TABLE fornecedores
    ADD CONSTRAINT uq_fornecedores_nome UNIQUE (nome);

-- Torna unicidade de nome case-insensitive.
-- Mantém validação de CNPJ já existente.

ALTER TABLE produtos DROP CONSTRAINT IF EXISTS uq_produtos_nome;
ALTER TABLE fornecedores DROP CONSTRAINT IF EXISTS uq_fornecedores_nome;

DROP INDEX IF EXISTS uq_produtos_nome_ci;
DROP INDEX IF EXISTS uq_fornecedores_nome_ci;

CREATE UNIQUE INDEX uq_produtos_nome_ci ON produtos (LOWER(nome));
CREATE UNIQUE INDEX uq_fornecedores_nome_ci ON fornecedores (LOWER(nome));

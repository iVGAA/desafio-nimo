FROM node:20-alpine

# Define o diretório de trabalho dentro do container
WORKDIR /app

# Copia apenas os arquivos de dependência primeiro
# Isso aproveita o cache das camadas do Docker se o package.json não mudar
COPY package*.json ./

# Instala as dependências
RUN npm install

# Copia o restante do código-fonte
COPY . .

# Expõe a porta 3000 que o Next.js usa por padrão
EXPOSE 3000

# Comando para rodar a aplicação em modo de desenvolvimento (com hot reload)
CMD ["npm", "run", "dev"]

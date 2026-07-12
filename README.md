# Sistema de Leilões - Go Expert

Projeto em Go para gerenciamento de leilões, usuários e lances, com fechamento automático de leilões usando goroutines e MongoDB como banco de dados.

## Pré-requisitos

Para executar o projeto e os testes, tenha instalado:

- Go 1.25 ou superior
- Docker
- Docker Compose

> Os testes usam Testcontainers para subir um MongoDB temporário. Por isso, o Docker precisa estar em execução antes de rodar os testes.

## Variáveis de ambiente

As configurações da aplicação ficam no arquivo `cmd/auction/.env`.

Principais variáveis:

| Variável | Descrição | Exemplo |
| --- | --- | --- |
| `BATCH_INSERT_INTERVAL` | Intervalo para processamento em lote dos lances | `20s` |
| `MAX_BATCH_SIZE` | Quantidade máxima de itens no lote | `4` |
| `AUCTION_INTERVAL` | Intervalo usado no fluxo de leilões | `20s` |
| `AUCTION_DURATION` | Tempo até o fechamento automático do leilão | `5m` |
| `MONGODB_URL` | URL de conexão com o MongoDB | `mongodb://admin:admin@mongodb:27017/auctions?authSource=admin` |
| `MONGODB_DB` | Nome do banco de dados | `auctions` |

## Como rodar os testes

Com o Docker em execução, rode todos os testes do projeto:

```bash
go test ./...
```

Para rodar apenas o teste de fechamento automático de leilão:

```bash
go test ./internal/infra/database/auction -run TestCreateAuction_ShouldCloseAuctionAutomatically -v
```

Esse teste cria um container MongoDB temporário, cria um leilão, aguarda o tempo configurado em `AUCTION_DURATION` e valida se o status foi alterado automaticamente para finalizado.

## Como rodar o projeto com Docker Compose

Na raiz do projeto, execute:

```bash
docker compose up --build
```

Se a sua instalação usa o comando legado do Docker Compose, execute:

```bash
docker-compose up --build
```

A aplicação ficará disponível em:

```text
http://localhost:8080
```

Para parar os containers:

```bash
docker compose down
```

Para parar os containers e remover também o volume do MongoDB:

```bash
docker compose down -v
```

## Como rodar o projeto localmente

Caso queira executar a aplicação fora do container, suba um MongoDB localmente e ajuste a variável `MONGODB_URL` para apontar para `localhost`.

Exemplo usando Docker apenas para o MongoDB:

```bash
docker run --name auction-mongodb \
  -p 27017:27017 \
  -e MONGO_INITDB_ROOT_USERNAME=admin \
  -e MONGO_INITDB_ROOT_PASSWORD=admin \
  -d mongo:latest
```

Depois, execute a aplicação com a URL local do MongoDB:

```bash
MONGODB_URL="mongodb://admin:admin@localhost:27017/auctions?authSource=admin" go run cmd/auction/main.go
```

A aplicação será iniciada na porta `8080`.

Para remover o container local do MongoDB quando não precisar mais dele:

```bash
docker rm -f auction-mongodb
```

## Endpoints principais

| Método | Rota | Descrição |
| --- | --- | --- |
| `GET` | `/auction` | Lista leilões |
| `GET` | `/auction/:auctionId` | Busca um leilão por ID |
| `POST` | `/auction` | Cria um novo leilão |
| `GET` | `/auction/winner/:auctionId` | Busca o lance vencedor de um leilão |
| `POST` | `/bid` | Cria um lance |
| `GET` | `/bid/:auctionId` | Lista lances de um leilão |
| `GET` | `/user/:userId` | Busca um usuário por ID |

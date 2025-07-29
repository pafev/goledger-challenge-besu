# GoLedger Besu Challenge

## Descrição do Projeto

Implementação pessoal (Paulo Fernando Vilarim) para o desafio técnico da GoLedger, que consiste em criar uma aplicação Go que interage com uma rede blockchain Hyperledger Besu QBFT. A aplicação deve ser capaz de interagir com smart contracts deployados na rede Besu, gerenciar valores de variáveis do smart contract, sincronizar dados entre blockchain e banco de dados SQL (Postgres na atual implementação), e expor funcionalidades através de uma API REST.

## Como Executar

### 1. Configuração do Ambiente

Gere o arquivo de variáveis de ambiente a partir do template:

```bash
cp .env.example .env
```

### 2. Inicialização do Cliente Besu

Execute os seguintes comandos para inicializar a rede Besu local:

```bash
cd scripts/besu
chmod +x ./startDev.sh
./startDev.sh
cd ../..
```

Este comando irá:
- Inicializar uma rede Besu local com 4 nós
- Deployar o smart contract SimpleStorage
- Exibir o endereço do contrato (anote este endereço para configurar no .env)

### 3. Inicialização do Banco de Dados

Execute o script para inicializar o PostgreSQL:

```bash
chmod +x ./scripts/db/startDev.sh
./scripts/db/startDev.sh
```

Este comando irá:
- Inicializar uma container com um banco de dados Postgres 1.17 operando
- Exibir a URL do banco de dados (automatizei a passagem da URL para o .env, mas anote esta URL para confirmar no .env)
Obs.: se optar por rodar o banco de dados locamente, deve haver no .env a query param `?sslmode=disable` adicionada à URL do banco de dados. Verifique

### 4. Configuração das Variáveis de Ambiente

Complete o arquivo `.env` com as informações obtidas nos passos anteriores. Dele deve ficar mais ou menos assim:

```env
# Informações da aplicação
APP_NAME=goledger-challenge-besu
APP_DOMAIN="localhost" # localhost for development env
APP_ENV="development"  # development or stage or production
APP_PORT=5000

# Configuração do Banco de Dados
DATABASE_URL=postgresql://user:pass@localhost:5432/goledger_challenge?sslmode=disable

# Configurações da rede Besu
BESU_URL=http://localhost:8545 # http://localhost:8545 for development env
SMART_CONTRACT_ADDR="<endereco_do_contrato_deployado>"
SMART_CONTRACT_ABI_PATH="scripts/besu/artifacts/contracts/SimpleStorage.sol/SimpleStorage.json"
```

### 5. Instalação de Dependências

Baixe as dependências do Go:

```bash
go mod download
go mod tidy
```

### 6. Execução da Aplicação

Execute a aplicação:

```bash
go run cmd/main.go
```

A aplicação estará disponível em `http://localhost:8080`.

## Sobre as Tecnologias

### Hyperledger Besu

Hyperledger Besu é um cliente Ethereum de código aberto desenvolvido pela ConsenSys e posteriormente para a Hyperledger Foundation. Besu implementa a Ethereum Virtual Machine (EVM) e suporta os principais algoritmos de consenso do Ethereum, incluindo Proof of Work (PoW) e Proof of Stake (PoS).

As principais características do Besu incluem:
- Compatibilidade total com a Ethereum MainNet
- Suporte para redes privadas e de consórcio
- Implementação completa da EVM
- APIs JSON-RPC compatíveis com Ethereum
- Suporte para diferentes algoritmos de consenso (IBFT, QBFT, Clique)

### Ethereum e Smart Contracts

Ethereum é uma plataforma blockchain descentralizada que permite a execução de contratos inteligentes, que são programas autoexecutáveis com os termos do acordo diretamente escritos em código, de modo a proporcionar segurança e autoridade sem possibilidade de censura, fraude ou interferência de terceiros.

### Besu como Cliente Ethereum

Como cliente Ethereum, o Besu oferece:
- Sincronização com a rede Ethereum
- Validação de transações e blocos
- Execução de smart contracts
- Capacidade de mineração/validação de blocos

Para ambientes de desenvolvimento e testes, o Besu pode ser configurado para executar redes privadas locais, permitindo desenvolvimento e teste de aplicações sem custos de gas da rede principal.

## Arquitetura da Aplicação

A aplicação foi desenvolvida de maneira a se basear nos princípios da Clean Architecture, porém evitando uma carga de over-engineering dado o escopo reduzido do projeto, ao mesmo tempo que abre portas para escalabilidade e fácil manutenção, devido a modularização, padrões de projeto aplicados e tratamento de erros adequado. O projeto possui uma divisão clara entre as camadas de aplicação e domínio. A estrutura organizacional segue o padrão de separação por funcionalidades dentro de cada camada.

### Camadas da Aplicação

**Camada de Domínio (`domain/`)**
- Contém as regras de negócio e entidades centrais da aplicação
- Implementa o padrão Repository para desacoplar a lógica de negócio dos detalhes de implementação dos clientes da aplicação (SGBD Postgresql e cliente Ethereum Besu)
- Organizada por funcionalidades (ex: `smart_contract/`)

**Camada de Aplicação (`app/`)**
- Responsável pela coordenação entre a camada de domínio e interfaces externas (requisições http)
- Contém os casos de uso e serviços da aplicação
- Organizada por funcionalidades seguindo a mesma estrutura de pastas da camada de domínio

### Padrão Repository

Na camada de domínio, foi implementado o padrão Repository para abstrair a lógica de interação com os clientes externos:
- **PostgreSQL**: Para persistência de dados
- **Besu (Ethereum Client)**: Para interação com a blockchain

Este padrão permite que a lógica de negócio seja independente dos detalhes de implementação dess sistemas externos, facilitando testes e manutenção do código.

```
├── internal/domain/
│   └── smart_contract/
│       ├── repository-besu.go # Interface do repository para o cliente Besu
│       ├── repository-db.go   # Interface do repository para o cliente Postgres
│       ├── model.go           # Entidades de domínio
├── internal/app/
│   └── smart_contract/
│       ├── handler.go         # HTTP handlers
│       └── service.go         # Coordenação de casos de uso
└── scripts/                   # semelhante a uma seção de "infrastructure"
    ├── db/                    # Implementação PostgreSQL
    └── besu/                  # Implementação Besu client
```

## Funcionalidades e Endpoints

A aplicação expõe quatro endpoints REST com tratamento personalizado de erros:

### GET /api/v1/smart-contract/
**Função**: Recupera o valor atual armazenado no smart contract
**Ação**: Realiza uma chamada de leitura diretamente na blockchain
**Retorno**: JSON com o valor atual do contrato

### GET /api/v1/smart-contract/check-value/:value
**Função**: Compara o valor do smart contract com o valor fornecido na rota
**Parâmetros**: `value` - valor a ser comparado via URL parameter
**Retorno**: JSON indicando se os valores são iguais (`true`/`false`)

### POST /api/v1/smart-contract/set-value
**Função**: Define um novo valor no smart contract
**Body**: JSON contendo:
```json
{
  "value": 42,
  "privateKey": "..."
}
```
**Ação**: Envia uma transação para a blockchain atualizando o valor
**Retorno**: JSON confirmando a transação

### POST /api/v1/smart-contract/sync
**Função**: Sincroniza o valor do smart contract com o banco de dados PostgreSQL
**Ação**: Lê o valor atual da blockchain e armazena/atualiza no banco de dados
**Retorno**: JSON confirmando a sincronização

## Tecnologias Utilizadas

- **Linguagem**: Go
- **Blockchain**: Hyperledger Besu
- **Cliente Ethereum Lib**: go-ethereum
- **Banco de Dados**: PostgreSQL
- **API**: REST
- **Containerização**: Docker & Docker Compose
- **Smart Contract**: SimpleStorage

## Pré-requisitos

Antes de executar a aplicação, certifique-se de ter instalado:

- Go (versão 1.24)
- Node.js e NPM (recomendo via NVM)
- Docker e Docker Compose (verificar se Docker Compose está atualizado)
- Hardhat
- Hyperledger Besu

## Estrutura do Projeto

```
.
├── cmd/
│   └── main.go                       # Ponto de entrada da aplicação
├── configs/                          # Arquivos de setup e configuracao das dependencias do projeto
│   └── app/
│       ├── config.go
│   └── besu/
│       ├── config.go
│   └── db/
│       ├── config.go
│   └── http/
│       └── config.go                 # aqui há as rotas e área para Dependency Injection (DI) da api
├── internal/
│   └── app/
│       └── smart_contract/
│           ├── handler.go            # HTTP handlers
│           └── service.go            # Casos de uso
│   └── domain/                       # Camada de domínio
│       └── smart_contract/
│           ├── repository-besu.go    # Interface do repository para o cliente Besu
│           ├── repository-db.go      # Interface do repository para o cliente Postgres
│           └── model.go              # Entidades
├── scripts/                          # Scripts de inicialização
│   ├── besu/
│   └── db/
├── .env.example
├── go.mod
├── go.sum
└── README.md
```

## Exemplos de Uso

### Recuperar valor do contrato

```bash
curl -X GET http://localhost:8080/api/v1/smart-contracts/
```

### Definir novo valor

```bash
curl -X POST http://localhost:8080/api/v1/smart-contracts/set-value \
  -H "Content-Type: application/json" \
  -d '{
    "value": 123,
    "privateKey": "8f2a55949038a9610f50fb23b5883af3b4ecb3c3bb792cbcefbd1542c692be63"
  }'
```

### Verificar valor

```bash
curl -X GET http://localhost:8080/api/v1/smart-contracts/check-value/123
```

### Sincronizar com banco de dados

```bash
curl -X POST http://localhost:8080/api/v1/smart-contracts/sync
```

## Tratamento de Erros

A aplicação implementa tratamento de erros abrangente para:

- Falhas de conexão com a rede Besu
- Erros de transação blockchain
- Problemas de conectividade com banco de dados
- Validação de dados de entrada
- Timeouts em operações de rede
- Parsing de parâmetros de rota
- Autenticação de chaves privadas

Todos os erros retornam respostas HTTP apropriadas com mensagens descritivas ou log de erro

## Notas Técnicas

### Interação com Smart Contract

A aplicação utiliza a biblioteca go-ethereum para interagir com a rede Besu:

- **Transações de Escrita**: Utilizam `bind.NewKeyedTransactorWithChainID` para assinar transações
- **Chamadas de Leitura**: Utilizam `bind.CallOpts` para chamadas que não modificam estado
- **ABI**: Carregado automaticamente dos artifacts gerados pelo Hardhat

### Segurança

- Chaves privadas são fornecidas via requisição (não armazenadas)
- Validação de entrada em todos os endpoints
- Timeouts configurados para operações blockchain
- Tratamento de conexões com banco de dados

### Performance

- Pool de conexões com banco de dados
- Reutilização de clientes Ethereum
- Timeouts apropriados para evitar travamentos
- Estruturas de dados otimizadas

## Desenvolvimento

### Logs

A aplicação utiliza logging estruturado para facilitar debugging e monitoramento em produção.

## Contribuição

Este projeto faz parte de um desafio técnico da GoLedger, desevolvido por mim, Paulo Fernando Vilarim (pafev).

**Desenvolvido para o GoLedger Challenge**

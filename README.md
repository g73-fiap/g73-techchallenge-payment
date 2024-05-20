# g73-techchallenge-payment

[![Go Reference](https://pkg.go.dev/badge/golang.org/x/example.svg)](https://pkg.go.dev/golang.org/x/example)

Este é um microsserviço responsável por gerenciar pagamentos de pedidos em uma lanchonete. Ele oferece endpoints para criar novos pedidos de pagamento, notificar o status do pagamento e interagir com um serviço de pagamento de terceiros.


## Tecnologias Utilizadas

- Linguagem de Programação: Go
- Banco de Dados: DynamoDB
- Framework Web: Gin
- Outras Dependências: AWS SDK para Go, lib/pq (para conexão com PostgreSQL)
- Docker



## Requisitos

- go 1.20 ou superior
- Docker
- kubernetes cluster (Docker desktop)
- kubectl



## Funcionalidades
- Criar novos pedidos de pagamento.
- Notificar o status do pagamento.
- Gerar QR code de pagamento usando um serviço de pagamento de terceiros.
- Armazenar informações de pedidos de pagamento no DynamoDB.



## Como Executar
Para executar este microsserviço, siga estas etapas:

**1. Configuração do Ambiente:**

- Certifique-se de ter Go instalado na sua máquina.
- Se estiver usando o DynamoDB localmente, certifique-se de ter Docker e Docker Compose instalados e configurados.


**2. Clone o Repositório:**

```bash
git clone https://github.com/seu-usuario/g73-techchallenge-payment.git

```

**3. Configuração do DynamoDB (Opcional):**

- Navegue até o diretório do projeto e execute o seguinte comando para iniciar o DynamoDB localmente:

```bash
docker-compose up -d
```

**4. Compilação e Execução do Microsserviço:**

- Navegue até o diretório do projeto e execute o seguinte comando para compilar o microsserviço:

```bash
go build -o g73-techchallenge-payment ./cmd/g73-techchallenge-payment
```

- Após a compilação, execute o seguinte comando para iniciar o microsserviço:
```bash
./g73-techchallenge-payment
```

**5. Testtando a API:**
- Uma vez que o microsserviço esteja em execução, você pode acessar a API em http://localhost:8080.
- Consulte a seção de endpoints abaixo para ver os endpoints disponíveis e suas descrições.




## Endpoints

- **POST /api/payment/create:** Cria um novo pedido de pagamento.

- **PUT /api/payment/notify/{orderId}/{paymentId}:** Notifica o status do pagamento para um pedido específico.

##  Documentação e Coverage
[Documentation](https://github.com/IgorRamosBR/g73-techchallenge-payment/tree/master/docs)


## Arquitetura
Clean Architecture com a estrutura de pastas baseada no [Standard Go Project Layout](https://github.com/golang-standards/project-layout#go-directories) 

```bash
├── cmd
├── configs
├── docs
├── internal
|   |── api
|   |── controllers
|   ├── core
|   │   ├── entities
|   │   ├── usecases
|   ├── infra
|   │   ├── drivers
|   │   ├── gateways
├── k8s
├── migrations
```
# GoFlow - Motor de Processamento de Lotes Concorrente

![GoFlow Banner](banner.png)

**GoFlow** é um motor de processamento de lotes (batch processing engine) robusto e eficiente, desenvolvido em Go. Ele foi projetado para lidar com grandes volumes de dados através de uma arquitetura baseada em orquestração concorrente e persistência em PostgreSQL.

## Funcionalidades

- **Processamento Concorrente**: Orquestrador inteligente que gerencia múltiplos trabalhadores (workers) para processar lotes em paralelo.
- **Ingestão via API**: Endpoint HTTP pronto para receber dados e enfileirar para processamento.
- **Dashboard em Tempo Real**: Interface CLI (Terminal) para monitoramento instantâneo das métricas de sucesso, falha e status dos lotes.
- **Persistência Confiável**: Integração completa com PostgreSQL para rastreamento de execuções e integridade dos dados.
- **Arquitetura Modular**: Separado em camadas claras de Ingestão, Processamento, Integração e Reporting.

## Tecnologias

- **Linguagem**: [Go](https://go.dev/) (v1.25+)
- **Banco de Dados**: [PostgreSQL](https://www.postgresql.org/)
- **Driver DB**: [pgx](https://github.com/jackc/pgx)
- **Logging**: [slog](https://pkg.go.dev/log/slog) (Estruturado)
- **Gestão de Ambiente**: [godotenv](https://github.com/joho/godotenv)
- **Orquestração**: Nativa com Goroutines e Channels.

## Como Começar

### Pré-requisitos

- Go 1.25 ou superior
- Docker (para o banco de dados) || também pode rodar localmente em um postgres instalado na máquina

### Configuração

1. Clone o repositório:
   ```bash
   git clone https://github.com/rphmauriciodev/goflow.git
   cd goflow
   ```

2. Configure as variáveis de ambiente:
   ```bash
   cp .env.example .env
   # Edite o .env com suas credenciais do banco de dados
   ```

3. Inicie o banco de dados (Exemplo via Docker):
   ```bash
   docker run --name goflow-db -e POSTGRES_USER=seu_usuario -e POSTGRES_PASSWORD=sua_senha -e POSTGRES_DB=goflow -p 5432:5432 -d postgres
   ```

4. Execute as migrações:
   ```bash
   make migrate-up
   ```

### Executando a Aplicação

Para iniciar o motor de processamento:
```bash
make run
```

Para visualizar o dashboard de métricas:
```bash
make dashboard
```

Ou em modo "watch" (atualização automática a cada 2s):
```bash
make watch-dashboard
```

## Estrutura do Projeto

```text
.
├── cmd/
│   ├── goflow/       # Ponto de entrada do motor principal
│   └── dashboard/    # Aplicação de monitoramento CLI
├── db/
│   └── migrations/   # Scripts SQL de evolução do banco
├── internal/
│   ├── ingestion/    # Handlers HTTP e recebimento de dados
│   ├── processing/   # Lógica do Orquestrador e Processadores
│   ├── platform/     # Infraestrutura (DB, Logger, Repositórios)
│   └── reporting/    # Lógica de agregação para o dashboard
└── Makefile          # Atalhos para comandos comuns
```

## Endpoints de Ingestão (Exemplo)

O servidor roda por padrão na porta `:8080`.

- `POST /batches`: Recebe um novo lote para processamento.

---

Desenvolvido por [Maurício](https://github.com/rphmauriciodev)

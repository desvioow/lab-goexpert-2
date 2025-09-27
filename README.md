# Lab Go Expert - Sistema de Consulta de Temperatura por CEP

Este projeto implementa uma arquitetura de microsserviços em Go que permite consultar a temperatura atual de uma cidade através do CEP brasileiro. O sistema utiliza as APIs ViaCEP e WeatherAPI para obter informações de localização e clima.

## Arquitetura do Sistema

O projeto é composto por três serviços principais:

### Service-A (Porta 8080)
- **Função**: Validação e roteamento de CEPs
- **Responsabilidades**:
  - Recebe requisições com CEP via POST
  - Valida formato do CEP (8 dígitos numéricos)
  - Encaminha CEPs válidos para o Service-B
  - Implementa rastreamento distribuído com Zipkin

### Service-B (Porta 8081)
- **Função**: Processamento de dados de localização e clima
- **Responsabilidades**:
  - Consulta dados de endereço via API ViaCEP
  - Busca informações meteorológicas via WeatherAPI
  - Converte temperaturas para Celsius, Fahrenheit e Kelvin
  - Retorna dados consolidados de temperatura

### Zipkin (Porta 9411)
- **Função**: Sistema de rastreamento distribuído
- **Responsabilidades**:
  - Coleta traces de requisições entre serviços
  - Fornece interface web para visualização de traces
  - Monitora performance e latência das chamadas

## Tecnologias Utilizadas

- **Go**: Linguagem principal do projeto
- **Chi Router**: Framework web para roteamento HTTP
- **Zipkin**: Sistema de rastreamento distribuído
- **Docker & Docker Compose**: Containerização e orquestração
- **APIs Externas**:
  - ViaCEP: Consulta de endereços por CEP
  - WeatherAPI: Dados meteorológicos

## Estrutura do Projeto

```
lab-goexpert-2/
├── service-a/
│   ├── main.go              # Serviço A
│   ├── Dockerfile           # Container do Serviço A
│   ├── go.mod              # Dependências Go
│   └── requests.http       # Exemplos de requisições
├── service-b/
│   ├── main.go              # Serviço B
│   ├── Dockerfile           # Container do Serviço B
│   ├── go.mod              # Dependências Go
│   ├── clients/            # Clientes para APIs externas
│   │   ├── viacep.go       # Cliente ViaCEP
│   │   ├── weather.go      # Cliente WeatherAPI
│   │   └── temperature.go  # Conversões de temperatura
│   └── dtos/               # Estruturas de dados
│       ├── requests.go     # DTOs de requisição
│       └── responses.go    # DTOs de resposta
└── docker-compose.yaml     # Orquestração dos serviços
```

## Como Executar

### Pré-requisitos
- Docker
- Docker Compose

### Execução
```bash
# Clone o repositório
git clone <url-do-repositorio>
cd lab-goexpert-2

# Inicie todos os serviços
docker-compose up --build

# Para executar em background
docker-compose up -d --build
```

### Verificação dos Serviços
```bash
# Service-A
curl http://localhost:8080

# Service-B  
curl http://localhost:8081

# Zipkin UI
# Acesse: http://localhost:9411
```

## Exemplos de Uso - Curl

### CEP Válido
```bash
curl -X POST http://localhost:8080 \
  -H "Content-Type: application/json" \
  -d '{"cep": "29902555"}'
```

### CEP Inválido (Formato)
```bash
curl -X POST http://localhost:8080 \
  -H "Content-Type: application/json" \
  -d '{"cep": "12345"}'
```

### CEP Inválido (Não Numérico)
```bash
curl -X POST http://localhost:8080 \
  -H "Content-Type: application/json" \
  -d '{"cep": "abcd1234"}'
```

## Fluxo de Dados

1. **Cliente** envia CEP para **Service-A** (porta 8080)
2. **Service-A** valida formato do CEP
3. **Service-A** encaminha CEP válido para **Service-B** (porta 8081)
4. **Service-B** consulta endereço na **API ViaCEP**
5. **Service-B** obtém dados meteorológicos na **WeatherAPI**
6. **Service-B** converte temperaturas para diferentes escalas
7. **Service-B** retorna dados consolidados para **Service-A**
8. **Service-A** repassa resposta para o **Cliente**


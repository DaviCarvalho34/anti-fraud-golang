# Sistema Anti-Fraude em Go

Sistema de detecção e prevenção de fraudes para transações financeiras.

## Características

- 🔍 Análise em tempo real de transações
- 🎯 Motor de regras configurável
- 📊 Sistema de pontuação de risco
- 🚨 Detecção de padrões suspeitos
- 📈 Análise comportamental
- 🌍 Validação de geolocalização
- 💳 Detecção de cartões roubados

## Estrutura do Projeto

```
anti-fraud-golang/
├── cmd/
│   └── api/           # Aplicação principal
├── internal/
│   ├── models/        # Modelos de dados
│   ├── rules/         # Motor de regras anti-fraude
│   ├── services/      # Lógica de negócio
│   └── handlers/      # Handlers HTTP
├── pkg/
│   └── utils/         # Utilitários compartilhados
└── tests/             # Testes
```

## Instalação

```bash
go mod download
```

## Executar

```bash
go run cmd/api/main.go
```

## API Endpoints

### Analisar Transação
```bash
POST /api/v1/transaction/analyze
```

### Verificar Status
```bash
GET /api/v1/health
```

## Exemplos

### Análise de Transação
```bash
curl -X POST http://localhost:8080/api/v1/transaction/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "transaction_id": "TXN123",
    "user_id": "USER456",
    "amount": 1000.00,
    "currency": "BRL",
    "merchant": "Loja XYZ",
    "location": {
      "country": "BR",
      "city": "São Paulo",
      "latitude": -23.55,
      "longitude": -46.63
    }
  }'
```

## Regras de Detecção

1. **Valor Alto**: Transações acima de R$ 10.000
2. **Velocidade**: Múltiplas transações em curto período
3. **Localização**: Mudanças geográficas impossíveis
4. **Horário Suspeito**: Transações em horários incomuns
5. **Padrão de Compra**: Desvio do comportamento normal

## Níveis de Risco

- **LOW** (0-30): Transação aprovada automaticamente
- **MEDIUM** (31-70): Requer revisão manual
- **HIGH** (71-100): Bloqueada automaticamente

## Licença

MIT

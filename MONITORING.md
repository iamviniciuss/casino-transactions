# SigNoz Monitoring Setup

## Como usar

### Executar apenas a aplicação (padrão)
```bash
docker compose up --build -d
```

### Executar aplicação + SigNoz (monitoramento)
```bash
docker compose --profile monitoring up --build -d
```

### Parar apenas o profile de monitoramento
```bash
docker compose --profile monitoring down
```

## ⚠️ Troubleshooting

### Schema Migrator em "Waiting"
Se o `signoz-schema-migrator-sync` ficar em status "Waiting" mas nos logs mostrar "finished":

1. **Aguarde**: O SigNoz tem um delay de 30s para aguardar a migração
2. **Verificar se o binário foi encontrado**:
   ```bash
   docker compose logs signoz
   ```
3. **Forçar restart se necessário**:
   ```bash
   docker compose --profile monitoring restart signoz
   ```

### OTel Collector parando inesperadamente
Se o `signoz-otel-collector` estiver com erro "collector stopped unexpectedly":

1. **Verificar se o SigNoz está funcionando**:
   ```bash
   docker compose logs signoz
   ```
2. **Reiniciar apenas o collector**:
   ```bash
   docker compose --profile monitoring restart signoz-otel-collector
   ```
3. **Verificar conectividade**:
   ```bash
   # Testar se as portas estão abertas
   curl -v http://localhost:4318/v1/traces
   ```

## Acesso

- **SigNoz Dashboard**: http://localhost:3301
- **Casino API**: http://localhost:9095  
- **OTLP Collector**: 
  - gRPC: http://localhost:4317
  - HTTP: http://localhost:4318

## Estrutura

- `signoz-clickhouse`: Banco de dados para métricas/traces
- `signoz-zookeeper`: Coordenação (separado do Kafka)
- `signoz`: Interface web do SigNoz
- `signoz-otel-collector`: Coletor OpenTelemetry

## Telemetria nas Aplicações

As aplicações `casino-api` e `casino-consumer` estão configuradas para enviar automaticamente:
- **Traces**: Rastreamento de requisições HTTP
- **Metrics**: Métricas de performance  
- **Logs**: Logs estruturados

As variáveis OTEL são configuradas automaticamente via docker-compose.

## Testando a Instrumentação

### 1. Fazer requests na API
```bash
# Fazer algumas requisições para gerar traces
curl "http://localhost:9095/transactions?user_id=573a37e7-832a-4ecd-9691-41ff29afb955&limit=5"
curl "http://localhost:9095/health"
```

### 2. Verificar traces no SigNoz
1. Acesse: http://localhost:3301
2. Vá em "Services" → "casino-api"
3. Visualize os traces das requisições HTTP

### 3. Verificar logs da instrumentação
```bash
# Ver logs da API
docker compose logs casino-api

# Ver logs do OTel Collector  
docker compose logs signoz-otel-collector
```
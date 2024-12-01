# Desafio técnico Go Expert FullCycle

## Funcionamento

Todas as requisições recebidas passam primeiramente no middleware que válida o acesso do IP ou token.
O middleware usa um algoritmo de janelas de tempo que toda a nova requisição é criado uma referência de tempo no redis com uma data/horário de expiração.

## Configuração

Copiar o arquivo `.env.exemplo` com o nome `.env`

### Alterar a váriaveis de ambiente

- Alterar a váriaveis no arquivo `.env`
- SERVER_PORT: Porta do servidor web
- REDIS_HOST: Endereço/Servidor de conexão do Redis
- REDIS_PORT: Porta de conexão do Redis 
- REDIS_PASSWORD: Senha conexão do Redis
- REDIS_DB= DB conexão do Redis 
- IP_MAX_REQUESTS: Máximo de requisições para um ip único
- TOKEN_MAX_REQUESTS: Máximo de requisições por token
- TIME_WINDOW_MILISECONDS: Janela de tempo

## Execução

- Executar o comando `docker compose up`
- O servidor deverá iniciar e responder na porta configurada

### Exemplos de requição

- Sem token
```bash
curl -X GET "http://localhost:8080" \
  -H "Host: localhost:8080" \
  -H "Accept: application/json"
```

- Com token
```bash
curl -X GET "http://localhost:8080" \
  -H "Host: localhost:8080" \
  -H "Accept: application/json" \
  -H "API_KEY: xUfdhCBLzcwudBzQr3r3Pp60HZAa13Q6"
```

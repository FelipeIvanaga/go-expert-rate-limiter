# Desafio técnico Go Expert FullCycle

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

`go run main.go  --url=http://google.com --requests=25 --concurrency=4`

## Executando com o Docker

`docker run --rm -it $(docker build -q .) --url=http://google.com --requests=25 --concurrency=4`
Isso fará o build da imagem, execução e quando concluído será deletada automaticamente.
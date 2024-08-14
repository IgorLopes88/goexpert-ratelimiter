# Go Expert

Desafio **Rate Limiter** do curso **Pós Go Expert**.

**Objetivo:**  Desenvolver um rate limiter em Go que possa ser configurado para limitar o número máximo de requisições por segundo com base em um endereço IP específico ou em um token de acesso.

### Execução da **aplicação**
Para executar a aplicação execute o comando:
```
git clone https://github.com/IgorLopes88/goexpert-ratelimiter.git
cd goexpert-ratelimiter
go mod tidy
```

Em seguida utilize o Docker para subir o **Redis**
```
docker compose up -d
```

O resultado deverá ser esse:

```
 ✔ Container redis-rlgo                  Started
```
Execute o comando para subir o **Rate Limiter**:
```
go run main.go
```

### Execução dos **testes**
Execute o comando:
```
go run test/main.go
```

### Informações do Funcionamento

Dentro da raiz do projeto, possui um arquivo `.env` onde estão as configurações básicas (bloqueio por IP e endereço do banco de dados).
```
IP_LIMIT_MAX_REQUEST=5
IP_BLOCK_TIME_SECONDS=20
WEBSERVER_PORT=8080
DB_ADDRESS=localhost
DB_PORT=6379
```

As informações dos **tokens** estão dento do arquivo `list-token.json`, onde cada token possui uma configuração própria.
```
{
    "as9i265ch4bcu": 
    {
        "LimitRPS": 10,
        "BlockForSec": 300
    },
    "bxp2sl28mv78p":
    {
        "LimitRPS": 10,
        "BlockForSec": 10
    },
    "cs79sf0hno6z7":
    {
        "LimitRPS": 5,
        "BlockForSec": 60
    }
}
```
O banco de dados foi implementado atravês de uma inteface, e pode ser facilmente subistituido por outro que também implemente `storage`.

Pronto!


### Correções de Bugs

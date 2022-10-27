# FotballPlayerTrivia Api
This is a portfolio-project and is not meant to be taken into production due to potential security risks.

### Description
FotballPlayerTrivia Api is a webserver built in Go, backed by a Postgres database for long term storage and Redis for caching. The idea of this project is to show how one could build an api with caching. This api has two endpoints, one to test connection and one to retrieve player-trivia.

### Quickstart
Use docker-compose with the supplied `compose.yml`.

`cp .env.example .env`

`docker-compose up`

First time starting also run:

`CURL GET localhost:8080/migrate`


### API

| Request method | Endpoint       | Parameters | Description                                                                                                                                                                                                                                                                                                                                                                    |   |
|----------------|----------------|------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|---|
| GET            | /player-trivia | name       | Returning fun fact of the player, if found in the database. ```{name: "name of player", country: "Country of player", country_official_name: "Official name of country, ex. Kingdom of Sweden", skin_area_coverage_ppm: "PPM of the players estimated skin coverage of the whole country.", population_share_ppm: "PPM of players contribution to the population of the country"}``` |   |
| GET            | /test          | None       | Tests if the server is reachable.                                                                                                                                                                                                                                                                                                                                              |   |
|                |                |            |       


### Database
The api is backed by a standard PostgresDB, which also will be deployed with Docker-Compose. 

### Cache
The api caches the results in a Redis-cache. To see this is action do the same request twice in a row.

#### Options
| Environment variables | Description                 |
|-----------------------|-----------------------------|
| DB_HOST               | host for postgres db.       |
| DB_PORT               | port for postgres db.       |
| DB_USER               | username for postgres db.   |
| DB_PASSWORD           | password for postgres db.   |
| DB_DEFAULT_DB         | default db for postgres db. |
| REDIS_HOST            | host for redis.             |
| REDIS_PORT            | port for redis.             |
| NINJA_API_TOKEN       | token for ninja-api.        |
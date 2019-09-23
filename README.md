# JWT Auth API using Kong 

## Requirements

Install Docker and Docker Compose (Docker Compose comes with Docker on Windows and MacOS)

## Run this Project

1. Build the docker services with `docker-compose build`
2. Run the auth service  database migrations with `docker-compose build auth && docker-compose run auth ./bin/goose up`
3. Start the services  in detached mode with `docker-compose up -d`
4. Setup Kong services, routes and plugins with `sh ./bootstrap.sh`

Kong Service will be live on `http://localhost:8000`

## Auth Service Health Check Endpoint

`http://localhost:8000/auth/healthcheck`

## Register a new user account(Account Creation)

```
curl -X POST http://localhost:8000/auth/accounts/register -d '{ "data": { "email_address": "test@mail.com", "password": "password" }}'
```

## Login to an existing account(Account Authentication)

```
curl -X POST http://localhost:8000/auth/accounts/login -d '{ "data": { "email_address": "test@mail.com", "password": "password" }}'
```

## Running Tests

1. Build Test Containers
```
docker-compose -f docker-compose.test.yaml build
```

2. Build goose in test containers 
```
docker-compose -f docker-compose.test.yaml run auth-test go build -o ./bin/goose ./cmd/goose/main.go
```

3. Run Test DB Migrations Containers
```
docker-compose -f docker-compose.test.yaml run auth-test ./bin/goose up
```

4. Run Tests

```
docker-compose -f docker-compose.test.yaml run auth-test go test
```


## Getting user credentials

## Routes
- Account Creation  `[POST] http://localhost:8000/auth/accounts/register`
- Authentication `[POST] http://localhost:8000/auth/accounts/login`
- Health Check `[GET] http://localhost:8000/auth/healthcheck`

## Migrations
Database migrations are handled by [goose](https://github.com/pressly/goose)

### Generate migrations

Migrations Schemas are created in `auth/db/migrations`

To generate sql migrations schemas, the prefered migration type execute in the root folder.
```
docker-compose build auth && docker-compose run auth ./bin/goose create <migration_name> sql
```

To generate golang migrations schemas, execute in the root folder.
```
docker-compose build auth && docker-compose run auth ./bin/goose create <migration_name> go
```

### Runing Migrations

```
docker-compose build auth && docker-compose run auth ./bin/goose up
```

### Migrations Status

To view executed and pending migrations run
```
docker-compose build auth && docker-compose run auth ./bin/goose status
```


# Para Services Assignment

[Assignment Details](https://github.com/ParaServices/coding-challenge-golang)

## Requirements

Install Docker and Docker Compose (Docker Compose comes with Docker on Windows and MacOS)

## Run this Project

1. Build the docker services with `docker-compose build`
2. Run the auth service  database migrations with `docker-compose build auth && docker-compose run auth ./bin/goose up`
3. Setup Kong services, routes and plugins with `sh ./bootstrap.sh`
4. Start the application with `docker-compose up`

Kong Service will be live on `http://localhost:8000`

## Auth Service Health Check Endpoint

`http://localhost:8000/healtcheck`

## Register a new user account(Account Creation)

```
curl -X POST http://localhost:8000/auth/accounts/register -d '{ "data": { "email_address": "test@mail.com", "password": "password" }}'
```

## Login to an existing account(Account Authentication)

```
curl -X POST http://localhost:8000/auth/accounts/login -d '{ "data": { "email_address": "test@mail.com", "password": "password" }}'
```

## Run Tests
```
docker-compose build -f docker-compose.yaml -f docker-compose.test.yaml
```


## Getting user credentials

## Routes
- Create User
- Get User Auth

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


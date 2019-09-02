# Para Services Assignment

[Assignment Details](https://github.com/ParaServices/coding-challenge-golang)

## Requirements

Install Docker and Docker Compose (Docker Compose comes with Docker on Windows and MacOS)

## Run Project

1. Build the docker services with `docker-compose up --build -d`
2. Run Migrations auth-db migratios
3. Run `bootstrap.sh` to setup kong server with plugins and routes
4. Start the application with ``

Application should be live on `http://localhost:8000`

## Routes
- Create User
- Get User Auth

## Migrations
Database migrations are handled by [goose](https://github.com/pressly/goose)

### Generate migrations

Migrations Schemas are created in `auth/db/migrations`

To generate sql migrations schemas, the prefered migration type execute in the root folder.
```
docker-compose run auth ./bin/goose create <migration_name> sql
```

To generate golang migrations schemas, execute in the root folder.
```
docker-compose run auth ./bin/goose create <migration_name> go
```

### Run Migrations

```
docker-compose build auth && docker-compose run auth ./bin/goose up
```

### Migrations Status

To view executed and pending migrations run
```
docker-compose run auth ./bin/goose status
```


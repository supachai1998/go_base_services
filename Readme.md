## Requirements

```
 1. Make for run command make file
 2. Docker & Docker compose for run docker container
 3. Env file in configs folder for environment variable
```

## Initial Setup

### 1. Create a secret.yaml file in the configs folder

` 1. cp ./configs/secret.example.yaml ./configs/secret.yaml`

### 2. Update the secret.yaml file with the required environment variables

` Update the secret.yaml file with the required environment variables`

### 3. Run docker compose

` make up`

# Structure Project

```
    1. configs: All environment variable and secret file
    2. cmd: All command line interface ex. go run cmd/server/main.go
    3. controllers: All controllers for handling request and response
        3.1 middleware: All middleware for handling request and response
        3.2 v1: All route for version 1
    4. database: All database communication
    5. domain: All domain logic for the application and database
    6. services: All services for handling business logic
    7. utils: All utility function
    8. storage: All storage for file and connection to storage
    9. validate: All validation and custom validation in domain for request and response
    10. server: All server configuration and setup project
```

# API Documentation

```
    Default Port: 8080
    Swagger: /api/v1/doc
    Health Check: /api/v1/health
    Staff Domain: /api/v1/staffs
        1. Create Staff: [POST] /api/v1/staffs
            1.1 Request
            {
                "email": "string",
                "password": "string",
                "first_name": "string",
                "last_name": "string",
                "phone": "string", // optional
                "role_id": "uuid" // optional ( with out role can't do anything, but you can update later)
            }
        2. Get Token: [POST] /api/v1/staffs/token ## Wait for implement to send token to email
        3. Get Verify: [POST] /api/v1/staffs/verify
        4. Login First Time: [POST] /api/v1/staffs/login
        5. Change Password: [POST] /api/v1/staffs/me/password
        6. Login: [POST] /api/v1/staffs/login
    User Domain: /api/v1/users
        1. Create User: [POST] /api/v1/users
        2. Get Token: [POST] /api/v1/users/token ## Wait for implement to send token to email
        3. Get Verify: [POST] /api/v1/users/verify
        4. Login First Time: [POST] /api/v1/users/login
        5. Change Password: [POST] /api/v1/users/me/password
        6. Login: [POST] /api/v1/users/login
    Role Domain: /api/v1/roles [restricted permission for staff]
        1. Create Role: [POST] /api/v1/roles
        2. Get All Role: [GET] /api/v1/roles
        3. Update Role: [PUT] /api/v1/roles/{id}
```

# Database

```
    1. Postgres
    2. Redis
```

# For Development

### don't forgot to change file in configs for development environment

```
    make dev
    curl http://localhost:3001/api/v1/health
test login
(don't seed or migrate database to production this for test and development only or you can change in storage/seed )
    curl --location 'http://localhost:3001/api/v1/staffs/login' \
        --form 'email="super_admin1@admin.com"' \
        --form 'password="Passw0rd!"'
```

# Url Shortener (a.k.a "Shorty") 

An applicacion that generates randome short URLS based on longer ones.

### Stack
- Go
- [Echo](https://echo.labstack.com/) - (web framework)
- [GoDotEnv](https://github.com/joho/godotenv) - Env variables
- [Compile Daemon](https://github.com/githubnemo/CompileDaemon) - watches files for changes
- Redis - Database
- Docker - Containers
- Docker compose - Orchestrate containers

## Development
Install Dependencies
```bash
go mod tidy
```

To run the backend in development mode with file watching run:
```bash
CompileDaemon -command="./urlshortener"
```

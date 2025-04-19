# bmstu_team_development

# Как запуститься ? 

**создай свой env в корне проекта**

```env
# .env
POSTGRES_USER=admin
POSTGRES_PASSWORD=secret123!
PGADMIN_EMAIL=admin@example.com
PGADMIN_PASSWORD=supersecret!
```

**запуститься** 
```bash 
docker-compose up --build
```

**выключить и удалить бд**
```bash 
docker-compose down -v
```

**запуск сервиса из командой строки**
```bash
POSTGRES_HOST="localhost" \
POSTGRES_PORT="5432" \
POSTGRES_USER="admin" \
POSTGRES_PASSWORD="secret123!" \
POSTGRES_DB_NAME="taskdb" \
go run ./cmd/main.go
```
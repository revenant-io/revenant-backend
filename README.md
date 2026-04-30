# Revenant Backend

Backend API en Go para el proyecto Revenant. Construido con Gin framework, PostgreSQL y Docker.

## Stack Tecnológico

- **Go 1.23+** - Lenguaje de programación
- **Gin** - Framework web HTTP
- **PostgreSQL 15** - Base de datos
- **Docker & Docker Compose** - Containerización
- **golang-migrate** - Migraciones de base de datos
- **Zap** - Logging estructurado
- **Validator** - Validación de datos
- **JWT** - Autenticación

## Estructura del Proyecto

```
revenant-backend/
├── internal/
│   ├── config/          # Configuración de la aplicación
│   ├── database/        # Conexión y migraciones
│   ├── logger/          # Logging centralizado
│   ├── models/          # Estructuras de datos
│   ├── services/        # Lógica de negocio
│   ├── server/
│   │   ├── handlers/    # Manejadores HTTP
│   │   └── middleware/  # Middlewares
│   └── utils/
│       ├── hash/        # Hashing de contraseñas
│       ├── jwt/         # Token JWT
│       └── validator/   # Validación
├── migrations/          # Migraciones SQL
├── go.mod              # Módulo Go
├── Dockerfile          # Imagen Docker
├── docker-compose.yml  # Orquestación de servicios
├── .env.example        # Variables de entorno ejemplo
└── main.go             # Punto de entrada
```

## Inicio Rápido

### Con Docker Compose

```bash
# Copiar variables de entorno
cp .env.example .env

# Iniciar la aplicación
make run

# Ver logs
make logs

# Detener la aplicación
make down
```

### Sin Docker (desarrollo local)

```bash
# Copiar variables de entorno
cp .env.example .env

# Instalar dependencias
make setup

# Iniciar PostgreSQL
docker run --name postgres-revenant \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=revenant \
  -p 5432:5432 \
  postgres:15-alpine

# Ejecutar aplicación
go run main.go
```

## Endpoints

### Health Check
- `GET /health` - Verificar estado de la aplicación

### Autenticación
- `POST /api/v1/auth/register` - Registro de usuario
- `POST /api/v1/auth/login` - Login de usuario

### Usuarios (requiere autenticación)
- `GET /api/v1/users/:id` - Obtener usuario por ID

## Variables de Entorno

Ver `.env.example` para referencia completa.

```env
SERVER_PORT=8080
ENVIRONMENT=development
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_NAME=revenant
DATABASE_USER=postgres
DATABASE_PASSWORD=postgres
JWT_SECRET=your-secret-key
```

## Comandos Útiles

```bash
# Setup inicial
make setup

# Correr aplicación con Docker
make run

# Ver logs
make logs

# Detener servicios
make down

# Correr tests
make test

# Lint
make lint

# Limpiar build
make clean
```

## Desarrollo

### Agregar una nueva migración

```bash
migrate create -ext sql -dir migrations -seq create_<tabla>_table
```

Editar los archivos `.up.sql` y `.down.sql` generados.

### Agregar un nuevo endpoint

1. Crear handler en `internal/server/handlers/`
2. Registrar ruta en `internal/server/server.go`
3. Crear servicio asociado en `internal/services/` si es necesario

### Testing

```bash
go test -v -cover ./...
```

## Notas de Arquitectura

- **Clean Architecture**: Separación clara entre capas
- **Modular**: Cada módulo es independiente y reutilizable
- **Testeable**: Servicios sin dependencia directa de frameworks
- **Escalable**: Estructura lista para crecer

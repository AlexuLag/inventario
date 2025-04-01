# Inventario API

API para el sistema de inventario.

## Requisitos

- Go 1.21 o superior
- MySQL 8.0 o superior
- Node.js 18 o superior (para el frontend)

## Configuración

1. Clonar el repositorio:
```bash
git clone https://github.com/tu-usuario/inventario.git
cd inventario
```

2. Configurar la base de datos MySQL:
```bash
# Crear la base de datos
mysql -u root -p
CREATE DATABASE inventario;
exit;

# Importar el esquema
mysql -u root -p inventario < cmd/schema.sql
```

3. Configurar las variables de entorno:
```bash
cp .env.example .env
# Editar .env con tus credenciales de MySQL
```

4. Instalar dependencias:
```bash
# Backend
cd API
go mod download

# Frontend
cd ../frontend
npm install
```

5. Iniciar el servidor:
```bash
# Backend
cd API
go run cmd/main.go

# Frontend
cd frontend
npm run dev
```

## Estructura del Proyecto

```
inventario/
├── API/
│   ├── cmd/
│   │   ├── main.go
│   │   └── schema.sql
│   ├── internal/
│   │   ├── domain/
│   │   ├── infrastructure/
│   │   │   └── repository/
│   │   ├── interface/
│   │   │   └── handler/
│   │   └── usecase/
│   ├── .env
│   └── go.mod
└── frontend/
    ├── src/
    ├── public/
    └── package.json
```

## API Endpoints

### Productos
- `POST /api/products` - Crear producto
- `GET /api/products` - Obtener todos los productos
- `GET /api/products/{id}` - Obtener producto por ID
- `PUT /api/products/{id}` - Actualizar producto
- `DELETE /api/products/{id}` - Eliminar producto

### Usuarios
- `POST /api/users` - Crear usuario
- `GET /api/users` - Obtener todos los usuarios
- `GET /api/users/{id}` - Obtener usuario por ID
- `PUT /api/users/{id}` - Actualizar usuario
- `DELETE /api/users/{id}` - Eliminar usuario

### Inventario
- `POST /api/stocks` - Crear item en inventario
- `GET /api/stocks` - Obtener todos los items
- `GET /api/stocks/{id}` - Obtener item por ID
- `PUT /api/stocks/{id}` - Actualizar item
- `DELETE /api/stocks/{id}` - Eliminar item
- `GET /api/stocks/product/{productId}` - Obtener items por producto
- `GET /api/stocks/serial/{serial}` - Obtener item por número de serie

### Proveedores
- `POST /api/providers` - Crear proveedor
- `GET /api/providers` - Obtener todos los proveedores
- `GET /api/providers/{id}` - Obtener proveedor por ID
- `PUT /api/providers/{id}` - Actualizar proveedor
- `DELETE /api/providers/{id}` - Eliminar proveedor

## Desarrollo

### Backend

La API está construida con Go y utiliza:
- Chi para el enrutamiento
- MySQL para la base de datos
- Clean Architecture para la estructura del proyecto

### Frontend

El frontend está construido con React y utiliza:
- TypeScript
- Tailwind CSS
- React Query para el manejo de estado y caché
- React Router para la navegación

## Pruebas

```bash
# Backend
cd API
go test ./...

# Frontend
cd frontend
npm test
```

## Licencia

MIT 
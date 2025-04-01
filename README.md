# Inventario Project

This project follows Clean Architecture principles and consists of a Go backend and React frontend.

## Project Structure

```
.
├── API/            # Go backend with Clean Architecture
│   ├── cmd/        # Application entry points
│   ├── internal/   # Private application code
│   │   ├── domain/     # Enterprise business rules
│   │   ├── usecase/    # Application business rules
│   │   ├── interface/  # Interface adapters
│   │   └── infrastructure/ # Frameworks and drivers
│   └── pkg/        # Public libraries
└── client-app/     # React frontend with Vite
```

## Backend Setup

1. Navigate to the API directory:
   ```bash
   cd API
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Run the server:
   ```bash
   go run cmd/main.go
   ```

The backend server will start on `http://localhost:8080`

## Frontend Setup

1. Navigate to the client-app directory:
   ```bash
   cd client-app
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Start the development server:
   ```bash
   npm run dev
   ```

The frontend development server will start on `http://localhost:5173`

## Clean Architecture

This project follows Clean Architecture principles:

- **Entities**: Core business rules and objects
- **Use Cases**: Application-specific business rules
- **Interface Adapters**: Controllers and presenters
- **Frameworks and Drivers**: External frameworks, tools, and drivers

The architecture ensures:
- Independence of frameworks
- Testability
- Independence of UI
- Independence of database
- Independence of any external agency 
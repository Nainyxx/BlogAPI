# BlogAPI (24.08.2026)

A RESTful API for a blogging platform built with **Golang**.

## Implemented Features:
- **Data Structures:** User, Post, and Comment domains along with their respective methods.
- **Blog Service:** The core data orchestration and controlling component of BlogAPI.

## TODO List:
- **Database Integration:** Set up a PostgreSQL database and connect it to the API.
- **Repository Layer:** Implement an interface between the Database and the Blog Service.
- **Routing & Handlers:** Set up the Gin HTTP request handlers.

## Tech Stack
- **Language:** Go (v1.26+)
- **Web Framework:** Gin Gonic
- **UUID Generation:** google/uuid
- **Security:** golang.org/x/crypto/bcrypt
- **Database:** PostgreSQL

## Project Structure

```text
├── main.go          # Application entry point
├── go.mod           # Project dependencies
├── Gin/             # Routing and HTTP handlers (Controllers)
├── service/         # Business logic and data orchestration (Services)
├── models/          # Data structures and domain validation (Models)
└── utils/           # Helper functions (Validators, Hashers)
```

## 💻 Getting Started

### Clone the Repository
```bash
git clone https://github.com
cd BlogAPI
```

### Download Dependencies
```bash
go mod download
```

### Run the Server
```bash
go run main.go
```
The server will spin up locally (typically listening on port `:8080`).

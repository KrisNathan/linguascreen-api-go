# Linguascreen Go Backend

**Not intended for production!**

1. Auth is non existent
2. No rate limiting
3. No logging

The backend is split into 2:

1. The Go Service
2. Tesseract.js Service

Rationale for using tesseract.js:

1. The existing go library for tesseract doesn't seem to be very active.
2. Tesseract.js doesn't depend on system dependencies therefore it's much more flexible. We can pick the versions we want.
3. This isn't for production use

## Setup

1. **Start all services (database, Tesseract.js, and Go backend):**
   ```sh
   docker-compose up -d
   ```

   This will start:
   - MySQL database on port 3306
   - Tesseract.js service on port 3000
   - Go backend API on port 8080

2. **Alternative: Run Go backend locally (requires database and Tesseract.js running in Docker):**
   ```sh
   # Start database and Tesseract.js
   docker-compose up -d db tesseract-service
   
   # Run Go backend locally
   go run .
   ```

## API Documentation

Once the server is running, you can access the Swagger UI for API documentation at:

`http://localhost:8080/swagger/index.html`

This provides interactive documentation for all endpoints, including request/response examples and the ability to test the API directly from the browser.

Generating OpenAPI JSON

```sh
go run github.com/swaggo/swag/cmd/swag@latest init
```
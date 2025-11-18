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

1. **Start the database and Tesseract.js service:**
   ```sh
   docker-compose up -d
   ```

2. **Run the Go backend:**
   ```sh
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
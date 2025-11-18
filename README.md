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

## Running

```sh
go run main.go
```
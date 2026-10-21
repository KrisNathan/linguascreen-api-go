FROM golang:1.27-alpine AS base

WORKDIR /app

FROM base AS deps

COPY go.mod go.sum ./
RUN go mod download

FROM deps AS build

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-linguascreen-api-go

FROM build AS runner

CMD ["/docker-linguascreen-api-go"]

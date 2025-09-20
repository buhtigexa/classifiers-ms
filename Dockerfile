# Che, this is our Dockerfile, alta optimizacion papa!
# First we build everything in a temporary container

FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install these dependencies si o si
# Without these no compile amigo
RUN apk add --no-cache gcc musl-dev

# Dale, first we get all our dependencies sorted out, viste?
COPY go.mod go.sum ./
RUN go mod download

# Ahora metemos todo el codigo, que se yo
COPY . .

# Compilamos la app con todos los chiches
RUN CGO_ENABLED=1 GOOS=linux go build -o classifier ./cmd/web

# Che, ahora si, la imagen final re livianita
FROM alpine:latest

WORKDIR /app

# Necesitamos estas cositas para que funcione todo
RUN apk add --no-cache ca-certificates mysql-client

# Copiamos el binario nomas, todo lo demas al tacho
COPY --from=builder /app/classifier .

# Dale, ponemos las variables de entorno asi anda todo piola
ENV GO_ENV=production
ENV SERVER_ADDR=:4000
ENV DB_MAX_OPEN_CONNS=25
ENV DB_MAX_IDLE_CONNS=25
ENV DB_MAX_IDLE_TIME=15m

# Puerto 4000, no te olvides eh!
EXPOSE 4000

# Y a correr nomas!
CMD ["./classifier"]
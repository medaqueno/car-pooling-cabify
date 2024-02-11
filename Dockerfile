# Etapa de compilación: Usa una imagen de Go oficial para compilar.
FROM golang:1.22-alpine3.19 AS builder
WORKDIR /app
# Asegúrate de copiar el go.mod y go.sum si estás utilizando go modules.
COPY go.mod go.sum ./
RUN go mod download
# Copia el código fuente de tu aplicación al contenedor.
COPY . .
# Configura las variables de entorno para la compilación cruzada si es necesario.
# En este caso, parece que quieres compilar para amd64, así que asegúrate de que esto coincide con tu objetivo.
# Si estás compilando para la misma arquitectura que tu máquina host, estas líneas pueden no ser necesarias.
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
# Compila tu aplicación. Asegúrate de ajustar el camino de tus fuentes si es necesario.
# El punto de entrada de tu aplicación Go parece estar en cmd/main.go.
RUN go build -o /car-pooling-challenge cmd/main.go

# Etapa de ejecución: Copia el binario compilado a una nueva imagen basada en Alpine.
FROM alpine:3.19
# Instala las dependencias necesarias en tu imagen final, si hay alguna.
RUN apk --no-cache add ca-certificates libc6-compat
# Copia el binario compilado desde el builder al contenedor final.
COPY --from=builder /car-pooling-challenge /car-pooling-challenge
# Expone el puerto que tu aplicación utiliza.
EXPOSE 9091
# Configura el contenedor para ejecutar tu aplicación.
ENTRYPOINT ["/car-pooling-challenge"]
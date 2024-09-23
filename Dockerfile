#Build
FROM golang:1.23-alpine AS build

WORKDIR /app
COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY ./ ./ 

RUN CGO_ENABLED=0 GOOS=linux go build -o ./server_exe ./cmd/server/server.go

# Deploy
 FROM debian

 WORKDIR /

COPY --from=build /app/server_exe /server_exe

EXPOSE 8080

USER root

ENTRYPOINT [ "/server_exe" ]
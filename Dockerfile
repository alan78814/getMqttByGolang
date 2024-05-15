FROM golang:1.19-bullseye AS build

ENV TZ=Asia/Taipei

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build \
  -ldflags="-linkmode external -extldflags -static" \
  -tags netgo \
  -o myapp

FROM alpine:3.18

ENV TZ=Asia/Taipei

WORKDIR /

COPY --from=build /app/myapp myapp

EXPOSE 8080

CMD ["/myapp"]

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

FROM scratch

ENV TZ=Asia/Taipei

WORKDIR /

COPY --from=build /app/myapp myapp

CMD ["/myapp"]

FROM golang:alpine AS builder

ENV TZ=Asia/Taipei

WORKDIR /go

COPY . .

RUN go build -o myapp

FROM scratch

ENV TZ=Asia/Taipei

WORKDIR /app

COPY --from=builder /go/myapp .

CMD ["./myapp"]
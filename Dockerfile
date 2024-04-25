FROM golang:alpine
ENV TZ=Asia/Taipei
COPY . /go
WORKDIR /go
RUN go build -o myapp
CMD ["./myapp"]

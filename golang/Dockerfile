FROM golang:1.19
# FROM golang:1.19.1-alpine3.15
RUN apt update
WORKDIR /go/src

# COPY . .

RUN go install github.com/cosmtrek/air@v1.40.0
# RUN air -c .air.toml

# CMD ["go", "run", "main.go"]
CMD ["air", "-c", ".air.toml"]

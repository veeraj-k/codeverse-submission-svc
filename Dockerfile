FROM golang:1.24-bookworm AS base


WORKDIR /build


COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o codeverse-submission-svc


EXPOSE 8082

CMD ["./codeverse-submission-svc"]

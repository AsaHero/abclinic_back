FROM golang:1.19-alpine3.17 AS builder

COPY . /github.com/AsaHero/abclinic
WORKDIR /github.com/AsaHero/abclinic

RUN go mod download

RUN CGO_ENABLED=0 GOARCH="amd64" GOOS=linux go build -ldflags="-s -w" -o ./bin/abclinic ./cmd/main.go 

FROM alpine:latest 

WORKDIR /root/

COPY --from=builder /github.com/AsaHero/abclinic/bin/abclinic .
COPY --from=builder /github.com/AsaHero/abclinic/bin/abclinic/auth_model.conf .
COPY --from=builder /github.com/AsaHero/abclinic/bin/abclinic/policy.csv .

EXPOSE 80

CMD [ "./abclinic" ]
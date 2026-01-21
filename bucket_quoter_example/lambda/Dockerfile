# ---- BUILD IMAGE FOR LAMBDA -----
FROM golang:1.25 AS build

WORKDIR /lambda
COPY certs/ certs/
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN env GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o main main.go

# ---- RUNTIME IMAGE FOR LAMBDA ---
FROM public.ecr.aws/lambda/provided:al2023

COPY --from=build /lambda/certs/dns-api.crt /etc/ssl/certs/dns-api.crt
COPY --from=build /lambda/main ./main

ENTRYPOINT [ "/usr/local/bin/aws-lambda-rie", "./main" ]

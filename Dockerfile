FROM golang:1.23 AS build

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./internal ./internal
COPY ./cmd ./cmd

RUN CGO_ENABLED=0 GOOS=linux go build -o zitadel_test_app ./cmd/zitadel_test_app

FROM gcr.io/distroless/base-debian11 AS build-release

COPY --from=build /app/zitadel_test_app /app/zitadel_test_app

ENV ZITADEL_TEST_APP_PORT=8080
EXPOSE 8080

CMD ["/app/zitadel_test_app"]

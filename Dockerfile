FROM golang:alpine AS build
COPY . /app

WORKDIR /app
RUN go build

FROM alpine
COPY --from=build /app/intra /app/intra

ENTRYPOINT /app/intra
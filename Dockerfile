FROM golang:1.13-alpine

RUN apk update && apk add --no-cache git ca-certificates

WORKDIR /src

COPY . ./

RUN go mod download

WORKDIR /src/server

RUN CGO_ENABLED=0 go build -a -o /src/app


FROM alpine

RUN apk update && apk add --no-cache tzdata

WORKDIR /app

COPY --from=0 /src/app /app/
COPY --from=0 /src/key /key/

EXPOSE 3000

ENTRYPOINT ["./app"]

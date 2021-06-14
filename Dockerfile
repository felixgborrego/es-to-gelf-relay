FROM golang:1.14.3-alpine AS build
WORKDIR /src
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /out/relay

FROM scratch AS bin
EXPOSE 8080
COPY --from=build /out/relay /
ENTRYPOINT ["/relay"]
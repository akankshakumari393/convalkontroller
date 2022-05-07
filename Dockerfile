##
## Build binary
##
FROM golang:1.18-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

COPY *.go ./

COPY ./pkg ./pkg

# RUN go mod vendor

# RUN go mod download

RUN CGO_ENABLED=0 go build

##
## RUN the binary
##

FROM alpine

COPY --from=build /app/convalkontroller /usr/local/bin

#USER root:root

ENTRYPOINT ["convalkontroller"]
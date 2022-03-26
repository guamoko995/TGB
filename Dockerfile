# syntax=docker/dockerfile:1

##
## Build
##

FROM golang:1.16-buster AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o /tgb

##
## Deploy
##

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /tgb /tgb
COPY --from=build /app/Game/mediaFiles/Murakami.txt /

USER nonroot:nonroot

ARG TGBtoken=your_token
ENV TGBtoken="${TGBtoken}"

ENTRYPOINT ["/tgb"]

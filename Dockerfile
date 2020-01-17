FROM golang:1.13-buster as build

ARG app_name=discord_rainbow_bot

WORKDIR /go/src/app
ADD . /go/src/app
RUN go build -o /go/bin/app

FROM gcr.io/distroless/base-debian10
COPY --from=build /go/bin/app /${app_name}
CMD ["/app"]
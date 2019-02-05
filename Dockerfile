FROM golang as builder
WORKDIR /go/src/github.com/graphql-services/graphql-files
COPY . .
RUN go get ./... 
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /tmp/app .

FROM alpine:3.5

WORKDIR /app

COPY --from=builder /tmp/app /usr/local/bin/app

# RUN apk --update add docker

# https://serverfault.com/questions/772227/chmod-not-working-correctly-in-docker
RUN chmod +x /usr/local/bin/app
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

ENTRYPOINT []
CMD [ "app", "start" ]
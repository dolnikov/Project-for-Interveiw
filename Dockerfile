FROM golang:1.22.3-alpine AS builder

ENV GOPRIVATE="gitlab.com/*"
ENV GOINSECURE="gitlab.com/*"

RUN apk update && apk upgrade && \
     apk add --no-cache git openssl ca-certificates

WORKDIR /app
COPY . /app/

RUN echo -e "machine gitlab.com\nlogin <CI_NETRC_LOGIN>\npassword <CI_NETRC_TOKEN>" > ~/.netrc
RUN chmod 600 ~/.netrc

RUN git config --global url."https://<CI_GITLAB_USER>:<CI_GITLAB_USER>@gitlab.com".insteadOf "http://gitlab.com"

RUN ls -lah /app

RUN go build -o /bin/gateway-api /app/cmd/server

FROM alpine:3.18
COPY --from=builder /bin/gateway-api /bin/gateway-api
COPY --from=builder /app/.env /bin/.env
WORKDIR /bin

RUN set -exu && \
    env && \
    chmod -R a+x /bin/* && \
    ls -lah /usr

EXPOSE 8000
EXPOSE 18000
ENTRYPOINT ["./gateway-api"]

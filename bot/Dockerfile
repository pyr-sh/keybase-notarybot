FROM golang:1.14.2 AS base
COPY . /go/src/github.com/pyr-sh/keybase-notarybot/bot
RUN go install -v github.com/pyr-sh/keybase-notarybot/bot

FROM keybaseio/client:stable
WORKDIR /app
RUN mkdir -p ./.cache/storage \
    && mkdir -p ./logs \
    && chown -R keybase /app
COPY --from=base /go/bin/bot /usr/bin/bot

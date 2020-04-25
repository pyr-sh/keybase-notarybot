#!/usr/bin/env bash
docker rm -f notarybot || true
docker build -t "pyr-sh/keybase-notarybot" .
docker run -it -p 4001:4001 --env-file ./.env --name notarybot pyr-sh/keybase-notarybot bot -d

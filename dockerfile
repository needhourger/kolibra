FROM node:20 AS node-builder

COPY ./frontend /frontend
WORKDIR /frontend

RUN yarn install && yarn build

FROM golang:1.22.5 AS go-builder

WORKDIR /kolibra

COPY . .

COPY --from=node-builder /frontend/dist /kolibra/static/dist

RUN  go env -w GOPROXY=https://goproxy.cn,direct

RUN make build_backend

FROM alpine:latest

WORKDIR /kolibra
COPY --from=go-builder /kolibra/kolibra ./kolibra
COPY ./config.yaml.example ./config.yaml
RUN mkdir data && mkdir data/library

EXPOSE 8080
CMD [ "./kolibra" ]

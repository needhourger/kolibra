FROM ubuntu:22.04 AS builder
RUN sed -i 's@//.*archive.ubuntu.com@//mirrors.ustc.edu.cn@g' /etc/apt/sources.list \
  && cat /etc/apt/sources.list \
  && apt-get update \
  && apt-get install -y golang nodejs make npm git\
  && go version && node -v && make --version && npm -v && git --version

RUN npm install --global yarn && yarn --version

WORKDIR /

RUN git clone --recurse-submodules https://github.com/needhourger/kolibra.git \
  && cd kolibra && ls

RUN make all

FROM alpine:latest

WORKDIR /kolibra
COPY --from=builder /kolibra/kolibra ./kolibra
COPY --from=builder /kolibra/config.yaml.example ./config.yaml
RUN mkdir data && mkdir data/library

EXPOSE 8080
CMD [ "./kolibra" ]

FROM golang:1.20

WORKDIR /go/src

ENV PATH="/go/bin:${PATH}"
ENV GO111MODULE=on
ENV CGO_ENABLED=1

RUN apt-get update -y && \
    apt-get install -y build-essential librdkafka-dev

CMD ["tail", "-f", "/dev/null"]

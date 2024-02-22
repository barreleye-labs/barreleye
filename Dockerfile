FROM golang:1.22-alpine
LABEL authors="nayoung, youngmin"

RUN apk add --no-cache gcc musl-dev

COPY go.mod /barreleye/
COPY go.sum /barreleye/

RUN cd /barreleye && go mod download

ADD . /barreleye

RUN cd /barreleye && go build -o ./bin/barreleye

#ENTRYPOINT ["/barreleye/bin/barreleye"]
#CMD ["-nodeName=$NODENAME"]
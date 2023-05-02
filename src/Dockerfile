FROM golang:1.20.2

RUN apt update && apt upgrade -y

RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

WORKDIR /opt/app/proxy
CMD ["air"]
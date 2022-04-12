FROM golang:1.18-buster

ENV GolangCILint_Version=1.45.2

RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
        sh -s -- -b $(go env GOPATH)/bin v${GolangCILint_Version}


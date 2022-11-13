FROM golang:1.19-bullseye

ENV GolangCILint_Version=1.50.1

RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
        sh -s -- -b $(go env GOPATH)/bin v${GolangCILint_Version}


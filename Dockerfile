# syntax=docker/dockerfile:1.7

ARG GOLANGCI_LINT_VERSION=v2.13.1
ARG KUBE_API_LINTER_VERSION=latest

FROM golangci/golangci-lint:${GOLANGCI_LINT_VERSION}-alpine AS builder

ARG KUBE_API_LINTER_VERSION

WORKDIR /src

RUN apk --no-cache add binutils gcc git mercurial musl-dev && \
    git config --global --add safe.directory '*'

RUN CGO_ENABLED=0 go install sigs.k8s.io/kube-api-linter/cmd/golangci-lint-kube-api-linter@${KUBE_API_LINTER_VERSION}

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 golangci-lint custom -v --name golangci-lint

FROM golangci/golangci-lint:${GOLANGCI_LINT_VERSION}-alpine

RUN apk --no-cache add binutils gcc git mercurial musl-dev && \
    git config --global --add safe.directory '*'

COPY --from=builder /src/golangci-lint /usr/bin/golangci-lint
COPY --from=builder /go/bin/golangci-lint-kube-api-linter /usr/bin/golangci-lint-kube-api-linter

ENTRYPOINT ["golangci-lint"]

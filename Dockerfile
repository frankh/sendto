FROM node:19.8.1 as frontend-build
WORKDIR /work

RUN npm i -g pnpm@7.29.3

COPY frontend/package.json frontend/pnpm-lock.yaml ./
RUN pnpm i

COPY frontend ./
RUN node esbuild.mjs

FROM golang:1.20 as build
WORKDIR /work

COPY go.mod go.sum ./
RUN go mod download -x

COPY --from=frontend-build /work/public ./frontend/public
COPY cmd cmd
COPY pkg pkg
COPY *.go .
RUN CGO_ENABLED=0 go build -v ./cmd/sendto


FROM alpine:latest
WORKDIR /data
COPY --from=build /work/sendto /sendto

ENTRYPOINT ["/sendto"]

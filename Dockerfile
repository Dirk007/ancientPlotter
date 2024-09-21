FROM golang:1.23-alpine AS builder

WORKDIR /src

ENV CGO_ENABLED=0

COPY . /src

RUN go build -o /bin/ancientPlotter


FROM alpine:3.20

RUN apk update && \
    apk add inkscape python3 py3-lxml py3-cssselect py3-numpy

WORKDIR /plotter

COPY --from=builder /bin/ancientPlotter /plotter/ancientPlotter
COPY assets /plotter/assets

ENTRYPOINT ["/plotter/ancientPlotter", "--serve"]


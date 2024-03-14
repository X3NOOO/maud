FROM golang:latest as builder
LABEL builder=true

WORKDIR /src

COPY . .

RUN mkdir -p out

RUN go mod download
RUN go build -o out/maud github.com/X3NOOO/maud

FROM photon
LABEL builder=false

WORKDIR /maud

COPY --from=builder /src/out/* .
COPY --from=builder /src/maud.toml .

ENTRYPOINT ["/maud/maud", "--config", "/maud/maud.toml"]

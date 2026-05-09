FROM --platform=${BUILDOS}/${BUILDARCH} golang:1-alpine AS src
ARG TARGETARCH
ARG TARGETOS

WORKDIR /build
COPY go.* .
COPY main.go .
COPY solver ./solver
COPY static ./static

FROM src AS test

RUN go test ./...

FROM src AS builder

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -trimpath -o numerology-${TARGETOS}-${TARGETARCH}

FROM --platform=${TARGETOS}/${TARGETARCH} golang:1-alpine
ARG TARGETARCH
ARG TARGETOS

COPY --from=builder /build/numerology-${TARGETOS}-${TARGETARCH} /app/numerology

CMD ["/app/numerology"]

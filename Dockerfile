FROM golang:1.26-trixie

WORKDIR /app

RUN go run github.com/mxschmitt/playwright-go/cmd/playwright@v0.6100.0 install chromium --with-deps

COPY lidl-coupons lidl-coupons

USER 65532:65532
ENTRYPOINT ["/app/lidl-coupons"]

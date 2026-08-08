FROM golang:1.26-trixie

WORKDIR /app

ENV HOME=/app
ENV PLAYWRIGHT_BROWSERS_PATH=/app/.cache

RUN useradd --uid 65532 -m -d /app --shell /bin/bash nonroot
RUN go run github.com/mxschmitt/playwright-go/cmd/playwright@v0.6100.0 install chromium --with-deps

COPY lidl-coupons lidl-coupons

RUN chown -R nonroot:nonroot /app

USER nonroot:nonroot
ENTRYPOINT ["/app/lidl-coupons"]

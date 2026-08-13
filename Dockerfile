FROM golang:1.26-trixie

WORKDIR /app

ENV HOME=/app
ENV PLAYWRIGHT_BROWSERS_PATH=/app/.cache

RUN useradd --uid 65532 -m -d /app --shell /bin/bash nonroot
COPY go.mod go.mod
RUN bash -c "go run github.com/mxschmitt/playwright-go/cmd/playwright@\$(awk '/mxschmitt\/playwright-go/ {print \$2}' go.mod) install chromium --with-deps && rm go.mod"

COPY lidl-coupons lidl-coupons

RUN chown -R nonroot:nonroot /app

USER nonroot:nonroot
ENTRYPOINT ["/app/lidl-coupons"]

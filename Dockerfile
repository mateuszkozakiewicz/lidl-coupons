FROM golang:1.27-trixie@sha256:9baa6b4187bbb98d240372a8a235ac0bb6b5ddd52bba1431dc2f7c0705862728

WORKDIR /app

ENV HOME=/app
ENV PLAYWRIGHT_BROWSERS_PATH=/app/.cache

RUN apt-get update && apt-get install -y --no-install-recommends xvfb=2:21.1.16-1.3+deb13u3 && rm -rf /var/lib/apt/lists/* && \
  useradd --uid 65532 -m -d /app --shell /bin/bash nonroot && \
  mkdir -p /tmp/.X11-unix && chmod 1777 /tmp/.X11-unix

COPY go.mod go.mod
RUN bash -c "go run github.com/mxschmitt/playwright-go/cmd/playwright@\$(awk '/mxschmitt\/playwright-go/ {print \$2}' go.mod) install chromium --with-deps && rm go.mod"

COPY lidl-coupons lidl-coupons
COPY entrypoint.sh entrypoint.sh
RUN mkdir -p /app/playwright-data && chown -R 65532:65532 /app && chmod +  /app/entrypoint.sh

VOLUME /app/playwright-data

USER 65532:65532
ENTRYPOINT ["/app/entrypoint.sh"]

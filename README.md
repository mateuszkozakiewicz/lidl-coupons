# lidl-coupons
List and activate all Lidl coupons, send a summary to Discord.

## Configuration options

Configuration happens via environment variables.

|Environment variable|Description                                                            |Default                             |
|:------------------:|:---------------------------------------------------------------------:|:----------------------------------:|
|API_URL             |URL of the Lidl API                                                    |`https://www.lidl.pl/prm/api/v1/PL/`|
|LOGIN_URL           |URL of the Lidl login page                                             |`https://www.lidl.pl/mla/`          |
|STORAGE_PATH        |Path to store Playwright data                                          |`./playwright-data`                 |
|TIMEOUT             |Timeout for Playwright operations                                      |`5s`                                |
|LOGIN               |Lidl account login (email)                                             |                                    |
|PASSWORD            |Lidl account password                                                  |                                    |
|TOKEN               |Lidl account token (optional, if not provided, login will be performed)|                                    |
|USERNAME            |Username for Discord notifications                                     |`Lidl Bot`                          |
|WEBHOOK_URL         |Discord webhook URL for notifications                                  |                                    |
|LOG_LEVEL           |Log level                                                              |`warn`                              |
|RETRIES             |Number of retries for login                                            |`10`                                |

## Run the application

You can use the binary from releases and export variables to run the application. For example:

```bash
export LOGIN="your-login"
export PASSWORD="your-password"
export WEBHOOK_URL="your-webhook-url"
./lidl-coupons
```

## Docker image

You can use the provided Docker image to run the application.

The container runs as an unprivileged user (uid `65532`). If you bind-mount a host directory make sure it is writable by that uid:

```bash
mkdir -p playwright-data
sudo chown -R 65532:65532 playwright-data
```

Then run the container with your credentials and optionally a webhook URL:

```bash
docker run -e LOGIN="your-login" -e PASSWORD="your-password" -e WEBHOOK_URL="your-webhook-url" -v $(pwd)/playwright-data:/app/playwright-data ghcr.io/mateuszkozakiewicz/lidl-coupons:latest
```

Alternatively, use a named Docker volume instead of a bind mount:

```bash
docker volume create lidl-coupons-data
docker run -e LOGIN="your-login" -e PASSWORD="your-password" -e WEBHOOK_URL="your-webhook-url" -v lidl-coupons-data:/app/playwright-data ghcr.io/mateuszkozakiewicz/lidl-coupons:latest
```

## Local development

Install dependencies:
* *OpenSUSE*: `sudo zypper install chromium`
* `go run github.com/mxschmitt/playwright-go/cmd/playwright@$(awk '/mxschmitt\/playwright-go/ {print $2}' go.mod) install chromium`
* `go mod tidy` to install Go dependencies
* `go run main.go` to run the application

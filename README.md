# lidl-coupons
List and activate all Lidl coupons, send coupon summary to Discord

## Local development

Install dependencies:
* *OpenSUSE*: `sudo zypper install chromium`
* `go run github.com/mxschmitt/playwright-go/cmd/playwright@$(awk '/playwright-community\/playwright-go/ {print $2}' go.mod) install chromium`
* `go run main.go` to run the application

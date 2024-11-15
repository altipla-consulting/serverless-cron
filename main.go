package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/altipla-consulting/errors"
	"github.com/altipla-consulting/telemetry"
	"github.com/altipla-consulting/telemetry/logging"
	"github.com/altipla-consulting/telemetry/sentry"
)

func main() {
	telemetry.Configure(sentry.Reporter(), logging.Standard())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	if err := run(ctx); err != nil {
		telemetry.ReportError(ctx, err)
	}
}

func run(ctx context.Context) error {
	u := os.Getenv("CRON_URL")
	if u == "" {
		return errors.New("CRON_URL env variable required")
	}
	token := os.Getenv("CRON_TOKEN")

	slog.Info("Invoke cron", slog.String("url", u))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, nil)
	if err != nil {
		return errors.Trace(err)
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Trace(err)
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Trace(err)
	}
	fmt.Println(string(content))

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("cron unexpected status %v", resp.Status)
	}

	slog.Info("Cron invoked successfully!")

	return nil
}

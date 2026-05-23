// cronitor-local is a self-hosted cron job monitoring daemon.
// It exposes a simple HTTP API for recording job runs and serves
// a web dashboard for visualizing job status and history.
package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/yourorg/cronitor-local/internal/monitor"
	"github.com/yourorg/cronitor-local/internal/store"
	"github.com/yourorg/cronitor-local/internal/web"
)

func main() {
	var (
		addr      = flag.String("addr", ":8080", "HTTP listen address")
		dbPath    = flag.String("db", "cronitor.db", "Path to the job state database file")
		interval  = flag.Duration("interval", 60*time.Second, "How often to check for overdue jobs")
		slackURL  = flag.String("slack-webhook", "", "Slack webhook URL for alerts (optional)")
		emailAddr = flag.String("alert-email", "", "Email address for alerts (optional, requires SMTP env vars)")
	)
	flag.Parse()

	// Initialize the persistent store.
	s, err := store.New(*dbPath)
	if err != nil {
		log.Fatalf("failed to open store at %s: %v", *dbPath, err)
	}

	// Build the list of alerters based on provided flags.
	var alerters []monitor.Alerter
	if *slackURL != "" {
		alerters = append(alerters, monitor.NewSlackAlerter(*slackURL))
	}
	if *emailAddr != "" {
		smtpHost := requireEnv("SMTP_HOST")
		smtpPort := requireEnv("SMTP_PORT")
		smtpFrom := requireEnv("SMTP_FROM")
		alerters = append(alerters, monitor.NewEmailAlerter(smtpHost, smtpPort, smtpFrom, *emailAddr))
	}

	// If no alerters are configured, fall back to logging only.
	if len(alerters) == 0 {
		log.Println("warning: no alerters configured — overdue jobs will only be logged")
		alerters = append(alerters, monitor.NewLogAlerter())
	}

	alerter := monitor.NewMultiAlerter(alerters...)

	// Start the background monitor that checks for overdue jobs.
	m := monitor.New(s, alerter)
	go m.Run(*interval)

	// Register HTTP handlers.
	h := web.NewHandler(s)
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", h.HandleHealthz)
	mux.HandleFunc("/api/jobs", h.HandleJobs)
	mux.HandleFunc("/api/jobs/", h.HandleJob)
	mux.HandleFunc("/", h.HandleIndex)

	log.Printf("cronitor-local listening on %s (db: %s, check interval: %s)", *addr, *dbPath, *interval)
	if err := http.ListenAndServe(*addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

// requireEnv returns the value of an environment variable or exits with a
// helpful message if it is not set.
func requireEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("required environment variable %s is not set", key)
	}
	return v
}

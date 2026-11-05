package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

type Event struct {
	Type string `json:"type"`
	Repo struct {
		Name string `json:"name"`
	} `json:"repo"`
	Payload json.RawMessage `json:"payload"`
}

type PushPayload struct {
	Ref string `json:"ref"`
}

type ActionPayload struct {
	Action string `json:"action"`
}

type CreateDeletePayload struct {
	Ref     string `json:"ref"`
	RefType string `json:"ref_type"`
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run() error {
	flag.Parse()
	if flag.NArg() == 0 {
		return errors.New("GitHub username not provided")
	}
	username := flag.Args()[0]

	url := fmt.Sprintf("https://api.github.com/users/%s/events", url.PathEscape(username))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	events, err := fetchEvents(ctx, url)
	if err != nil {
		return err
	}

	for _, e := range events {
		fmt.Println("-", describe(e))
	}

	return nil
}

func fetchEvents(ctx context.Context, url string) ([]Event, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("User-Agent", "github-activity-cli")
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	switch {
	case resp.StatusCode == http.StatusNotFound:
		return nil, fmt.Errorf("user not found")
	case resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusTooManyRequests:
		return nil, fmt.Errorf("rate limit exceeded: %s", resp.Status)
	case resp.StatusCode < 200 || resp.StatusCode >= 300:
		return nil, fmt.Errorf("unexpected HTTP status: %s", resp.Status)
	}

	var events []Event
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return nil, fmt.Errorf("decode response: %v", err)
	}

	return events, nil
}

func describe(e Event) string {
	switch e.Type {
	case "PushEvent":
		var p PushPayload
		if err := json.Unmarshal(e.Payload, &p); err != nil || p.Ref == "" {
			return fmt.Sprintf("Pushed to %s", e.Repo.Name)
		}
		branch := strings.TrimPrefix(p.Ref, "refs/heads/")
		return fmt.Sprintf("Pushed to %s (%s)", e.Repo.Name, branch)

	case "WatchEvent":
		return fmt.Sprintf("Starred %s", e.Repo.Name)

	case "IssuesEvent":
		var p ActionPayload
		_ = json.Unmarshal(e.Payload, &p)
		return fmt.Sprintf("%s an issue in %s", capitalize(p.Action), e.Repo.Name)

	case "CreateEvent", "DeleteEvent":
		var p CreateDeletePayload
		_ = json.Unmarshal(e.Payload, &p)
		verb := "Created"
		if e.Type == "DeleteEvent" {
			verb = "Deleted"
		}
		if p.RefType == "repository" || p.Ref == "" {
			return fmt.Sprintf("%s repository %s", verb, e.Repo.Name)
		}
		return fmt.Sprintf("%s %s %s in %s", verb, p.RefType, p.Ref, e.Repo.Name)

	case "ReleaseEvent":
		var p ActionPayload
		_ = json.Unmarshal(e.Payload, &p)
		return fmt.Sprintf("%s a release in %s", capitalize(p.Action), e.Repo.Name)

	case "ForkEvent":
		return fmt.Sprintf("Forked repository %s", e.Repo.Name)

	case "PublicEvent":
		return fmt.Sprintf("Made %s public", e.Repo.Name)

	default:
		return fmt.Sprintf("%s in %s", capitalize(e.Type), e.Repo.Name)
	}
}

func capitalize(s string) string {
	if s == "" {
		return ""
	}
	r, size := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[size:]
}

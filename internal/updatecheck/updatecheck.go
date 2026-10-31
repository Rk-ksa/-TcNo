package updatecheck

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"TcNo-Acc-Switcher/internal/api"
	"TcNo-Acc-Switcher/internal/appclient"
)

const (
	failStateFile    = "TcNo-Acc-Switcher.lastUpdateCheckFail.json"
	httpTimeout      = 15 * time.Second
	LaunchCheckDelay = 700 * time.Millisecond
)

type failStateJSON struct {
	At string `json:"at"`
}

type LaunchAPICheckResult struct {
	Latest  string
	Message string
	Err     error
}

// ParseVersionClock parses a semver version string (e.g. "4.0.0" or "v4.0.0")
// into a comparable time.Time. MAJOR.MINOR.PATCH maps to year.month.day.
func ParseVersionClock(s string) (time.Time, error) {
	s = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(s, "\r", ""), "\n", ""))
	s = strings.TrimPrefix(s, "v")

	parts := strings.Split(s, ".")
	if len(parts) == 3 {
		values := [3]int{}
		for i, part := range parts {
			value, err := strconv.Atoi(part)
			if err != nil || value < 0 {
				return time.Time{}, fmt.Errorf("updatecheck: invalid semver format: %q", s)
			}
			values[i] = value
		}
		return time.Date(values[0], time.Month(values[1]), values[2], 0, 0, 0, 0, time.UTC), nil
	}
	return time.Time{}, fmt.Errorf("updatecheck: invalid semver format: %q", s)
}

func IsUpToDate(current, latest string) bool {
	cur, errC := ParseVersionClock(current)
	lat, errL := ParseVersionClock(latest)
	if errC != nil || errL != nil {
		return false
	}
	return lat.Equal(cur) || cur.After(lat)
}

// IsUpdateAvailable reports only a strictly newer, parseable release. An
// invalid version from an API or updater provider must never become an update
// announcement merely because it cannot be compared.
func IsUpdateAvailable(current, latest string) bool {
	cur, errC := ParseVersionClock(current)
	lat, errL := ParseVersionClock(latest)
	return errC == nil && errL == nil && lat.After(cur)
}

func FetchLatestVersion(ctx context.Context, client *http.Client, currentVersion string) (version string, message string, err error) {
	if appclient.IsOfflineMode() {
		return "", "", appclient.ErrOfflineMode
	}
	if client == nil {
		client = appclient.Shared
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, updateAPIURL(currentVersion), nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("User-Agent", api.UserAgent(strings.TrimSpace(currentVersion)))
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if err != nil {
		return "", "", err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", "", fmt.Errorf("updatecheck: HTTP %d", resp.StatusCode)
	}
	lines := strings.SplitN(strings.TrimSpace(string(body)), "\n", 2)
	version = strings.TrimSpace(lines[0])
	if len(lines) > 1 {
		message = strings.TrimSpace(lines[1])
	}
	return version, message, nil
}

// FetchLaunchAPICheck hits the tcno.co API endpoint used for launch telemetry
// and update fallback data.
func FetchLaunchAPICheck(ctx context.Context, currentVersion string) LaunchAPICheckResult {
	if appclient.IsOfflineMode() {
		return LaunchAPICheckResult{Err: appclient.ErrOfflineMode}
	}
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, cancel := context.WithTimeout(ctx, httpTimeout)
	defer cancel()

	latest, message, err := FetchLatestVersion(ctx, appclient.Shared, currentVersion)
	return LaunchAPICheckResult{
		Latest:  latest,
		Message: message,
		Err:     err,
	}
}

// HandleLaunchAPICheckResult applies launch fallback UI behaviour to an API
// result that may already have been fetched for launch telemetry.
// Fail toasts are throttled to once per day.
func HandleLaunchAPICheckResult(result LaunchAPICheckResult, exeDir string, currentVersion string, onUpdateAvailable func(message string), onCheckFailed func()) {
	err := result.Err
	if err != nil {
		if shouldEmitFailToast(exeDir) {
			_ = writeFailTimestamp(exeDir)
			if onCheckFailed != nil {
				onCheckFailed()
			}
		}
		return
	}
	// A successful HTTP response is not necessarily a version response. Do not
	// advertise an update when either side is missing or malformed.
	if !IsUpdateAvailable(currentVersion, result.Latest) {
		return
	}
	if onUpdateAvailable != nil {
		onUpdateAvailable(result.Message)
	}
}

// RunLaunchAPICheck runs the tcno.co API check used as an updater fallback on launch.
func RunLaunchAPICheck(ctx context.Context, exeDir string, currentVersion string, onUpdateAvailable func(message string), onCheckFailed func()) {
	HandleLaunchAPICheckResult(FetchLaunchAPICheck(ctx, currentVersion), exeDir, currentVersion, onUpdateAvailable, onCheckFailed)
}

// RunManualCheck checks for updates on user request. Returns "available", "up-to-date", or "failed".
func RunManualCheck(ctx context.Context, currentVersion string, onUpdateAvailable func(message string)) string {
	if appclient.IsOfflineMode() {
		return "failed"
	}
	latest, message, err := FetchLatestVersion(ctx, appclient.Shared, currentVersion)
	if err != nil {
		return "failed"
	}
	if !IsUpdateAvailable(currentVersion, latest) {
		return "up-to-date"
	}
	if onUpdateAvailable != nil {
		onUpdateAvailable(message)
	}
	return "available"
}

package updatecheck

import (
	"errors"
	"testing"
)

func TestIsUpdateAvailable(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, current, latest string
		want                  bool
	}{
		{name: "newer release", current: "4.0.8", latest: "v4.0.9", want: true},
		{name: "same release", current: "4.0.8", latest: "4.0.8", want: false},
		{name: "older release", current: "4.0.9", latest: "4.0.8", want: false},
		{name: "invalid current", current: "local-build", latest: "4.0.9", want: false},
		{name: "invalid latest", current: "4.0.8", latest: "temporarily unavailable", want: false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsUpdateAvailable(tc.current, tc.latest); got != tc.want {
				t.Fatalf("IsUpdateAvailable(%q, %q) = %t, want %t", tc.current, tc.latest, got, tc.want)
			}
		})
	}
}

func TestHandleLaunchAPICheckResult_UpdateAvailable(t *testing.T) {
	t.Parallel()
	var got string
	HandleLaunchAPICheckResult(LaunchAPICheckResult{
		Latest:  "4.1.2",
		Message: "update ready",
	}, "", "4.1.1", func(message string) {
		got = message
	}, func() {
		t.Fatal("did not expect failure callback")
	})
	if got != "update ready" {
		t.Fatalf("got update message %q, want %q", got, "update ready")
	}
}

func TestHandleLaunchAPICheckResult_UpToDate(t *testing.T) {
	t.Parallel()
	HandleLaunchAPICheckResult(LaunchAPICheckResult{
		Latest:  "4.1.1",
		Message: "ignored",
	}, "", "4.1.2", func(string) {
		t.Fatal("did not expect update callback")
	}, func() {
		t.Fatal("did not expect failure callback")
	})
}

func TestHandleLaunchAPICheckResult_InvalidVersionDoesNotAnnounceUpdate(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, current, latest string
	}{
		{name: "invalid latest", current: "4.0.8", latest: "temporarily unavailable"},
		{name: "trailing message", current: "4.0.8", latest: "4.0.9 service unavailable"},
		{name: "empty latest", current: "4.0.8", latest: ""},
		{name: "invalid current", current: "development", latest: "4.0.9"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			HandleLaunchAPICheckResult(LaunchAPICheckResult{Latest: tc.latest}, "", tc.current, func(string) {
				t.Fatal("did not expect update callback for an invalid version")
			}, nil)
		})
	}
}

func TestHandleLaunchAPICheckResult_ThrottlesFailureCallback(t *testing.T) {
	exeDir := t.TempDir()
	calls := 0
	onFailed := func() {
		calls++
	}

	result := LaunchAPICheckResult{Err: errors.New("network down")}
	HandleLaunchAPICheckResult(result, exeDir, "4.1.1", nil, onFailed)
	HandleLaunchAPICheckResult(result, exeDir, "4.1.1", nil, onFailed)

	if calls != 1 {
		t.Fatalf("failure callback calls = %d, want 1", calls)
	}
}

package service

import (
	"crawlquery/api/domain"
	"testing"
	"time"
)

func TestCheckAndThrottle(t *testing.T) {
	service := NewService()

	// Test adding a new URL
	url := domain.URL("http://example.com")
	throttled, err := service.CheckAndThrottle(url)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if !throttled {
		t.Errorf("expected to throttle URL %s, but it was not throttled", url)
	}

	// Test throttling the same URL again
	throttled, err = service.CheckAndThrottle(url)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if throttled {
		t.Errorf("expected not to throttle URL %s, but it was throttled", url)
	}

	// Test adding a different URL
	anotherURL := domain.URL("http://another-example.com")
	throttled, err = service.CheckAndThrottle(anotherURL)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if !throttled {
		t.Errorf("expected to throttle URL %s, but it was not throttled", anotherURL)
	}

	// Ensure the second URL is not throttled again
	throttled, err = service.CheckAndThrottle(anotherURL)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if throttled {
		t.Errorf("expected not to throttle URL %s, but it was throttled", anotherURL)
	}
}

func TestCheckAndThrottleInvalidURL(t *testing.T) {
	service := NewService()

	// Test adding an invalid URL
	invalidURL := domain.URL("http//google.com")
	_, err := service.CheckAndThrottle(invalidURL)

	if err == nil {
		t.Errorf("expected an error, got nil")
	}
}

func TestWithRateLimit(t *testing.T) {
	rateLimit := 100 * time.Millisecond
	service := NewService()

	// Apply rate limit option
	WithRateLimit(rateLimit)(service)

	// Add a URL
	url := domain.URL("http://example.com")
	throttled, err := service.CheckAndThrottle(url)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if !throttled {
		t.Errorf("expected to throttle URL %s, but it was not throttled", url)
	}

	// Ensure URL is throttled immediately
	throttled, err = service.CheckAndThrottle(url)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if throttled {
		t.Errorf("expected not to throttle URL %s immediately after, but it was throttled", url)
	}

	// Wait for rate limit duration
	time.Sleep(rateLimit + 10*time.Millisecond)

	// Ensure URL is throttled again after rate limit duration
	throttled, err = service.CheckAndThrottle(url)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if !throttled {
		t.Errorf("expected to throttle URL %s after rate limit duration, but it was not throttled", url)
	}
}

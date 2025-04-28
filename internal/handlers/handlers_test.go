package handlers

import (
	"testing"
)

func TestGenerateSessionID(t *testing.T) {
	username := "marcelo"
	expectedSessionID := "session_marcelo"

	sessionID := generateSessionID(username)

	if sessionID != expectedSessionID {
		t.Errorf("expected session ID %s but got %s", expectedSessionID, sessionID)
	}
}

func TestIsSessionValid(t *testing.T) {
	sessionID := "session_marcelo"

	mu.Lock()
	sessionMap[sessionID] = "marcelo"
	mu.Unlock()

	if !isSessionValid(sessionID) {
		t.Errorf("expected session ID %s to be valid, but it was invalid", sessionID)
	}

	if isSessionValid("nonexistent_session") {
		t.Errorf("expected nonexistent session to be invalid, but it was reported valid")
	}
}

func TestGetUsername(t *testing.T) {
	sessionID := "session_marcelo"

	mu.Lock()
	sessionMap[sessionID] = "marcelo"
	mu.Unlock()

	username := getUsername(sessionID)

	if username != "marcelo" {
		t.Errorf("expected username 'marcelo', but got '%s'", username)
	}
}

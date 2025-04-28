package handlers

import (
	"html/template"
	"net/http"
	"sync"
)

var (
	tpl *template.Template

	// hadcoded users, for now ! TODO : Use SQLite to store users
	users = map[string]string{
		"marcelo": "secret123",
		"admin":   "admin123",
	}
	sessionMap = make(map[string]string) // sessionID -> username
	mu sync.Mutex
)

func LoadTemplates() {
	tpl = template.Must(template.ParseGlob("web/templates/*.html"))
}

// generateSessionID creates a basic session ID (could be improved later).
func generateSessionID(username string) string {
	return "session_" + username
}

// LoginHandler handles user login.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tpl.ExecuteTemplate(w, "login.html", nil)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	mu.Lock()
	defer mu.Unlock()

	if users[username] == password {
		sessionID := generateSessionID(username)
		sessionMap[sessionID] = username

		http.SetCookie(w, &http.Cookie{
			Name:  "session",
			Value: sessionID,
			Path:  "/",
		})

		http.Redirect(w, r, "/chat", http.StatusSeeOther)
		return
	}

	tpl.ExecuteTemplate(w, "login.html", "Invalid credentials")
}

// LogoutHandler handles user logout.
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err == nil {
		mu.Lock()
		delete(sessionMap, cookie.Value)
		mu.Unlock()

		http.SetCookie(w, &http.Cookie{
			Name:   "session",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// ChatHandler renders the chat page after checking authentication.
func ChatHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil || !isSessionValid(cookie.Value) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	username := getUsername(cookie.Value)
	tpl.ExecuteTemplate(w, "chat.html", username)
}

// isSessionValid checks if a session ID is valid.
func isSessionValid(sessionID string) bool {
	mu.Lock()
	defer mu.Unlock()
	_, exists := sessionMap[sessionID]
	return exists
}

// getUsername retrieves the username for a given session ID.
func getUsername(sessionID string) string {
	mu.Lock()
	defer mu.Unlock()
	return sessionMap[sessionID]
}

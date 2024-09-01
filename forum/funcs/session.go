package funcs

import (
	"context"
	"database/sql"
	"net/http"
	"time"
	"log"
	"crypto/rand"
	"encoding/base64"
	"forum/DB"
)

var db *sql.DB

func init() {
	var err error
	db, err = DB.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	DB.CreateSessionsTable()
}

func createSessionID() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func invalidateExistingSessions(userID int) error {
    _, err := db.Exec("DELETE FROM sessions WHERE user_id = ?", userID)
    return err
}
// SetSession sets a session cookie for the logged-in user by creating a new session ID and inserting it into the sessions table in the database. The session cookie is set to expire in 24 hours.
func SetSession(w http.ResponseWriter, userID int) {
	// Invalidate existing sessions for the user
	err := invalidateExistingSessions(userID)
	if err != nil {
		log.Println("Error invalidating existing sessions:", err)
	}

	// Create a new session ID
	sessionID := createSessionID()
	expiry := time.Now().Add(24 * time.Hour)

	// Insert the new session into the database
	_, err = db.Exec("INSERT INTO sessions (session_id, expiration, user_id) VALUES (?, ?, ?)", sessionID, expiry, userID)
	if err != nil {
		log.Println("Error creating session:", err)
		return
	}

	// Set the session cookie
	cookie := http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Expires:  expiry,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}

// GetSession checks if the user is logged in by checking the session cookie. If the user is logged in, it returns the user ID and true. If the user is not logged in, it returns 0 and false.
func GetSession(r *http.Request) (int, bool) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return 0, false
	}

	var userID int
	var expiration time.Time
	err = db.QueryRow("SELECT user_id, expiration FROM sessions WHERE session_id = ?", cookie.Value).Scan(&userID, &expiration)
	if err != nil || expiration.Before(time.Now()) {
		return 0, false
	}

	return userID, true
}
// ClearSession deletes the session from the database and clears the session cookie.  This is used when the user logs out.
func ClearSession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == nil {
		_, err := db.Exec("DELETE FROM sessions WHERE session_id = ?", cookie.Value)
		if err != nil {
			log.Println("Error clearing session:", err)
		}

		cookie := http.Cookie{
			Name:     "session_id",
			Value:    "",
			Path:     "/",
			Expires:  time.Unix(0, 0),
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
	}
}
// SessionMiddleware is a middleware that checks if the user is logged in and If the user is not logged in, it redirects them to the login page.
func SessionMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, exists := GetSession(r)
		if !exists {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), "userID", userID))
		next.ServeHTTP(w, r)
	})
}

package funcs

import (
	"forum/DB"
	"net/http"
	"strings"
	//"time"
	"golang.org/x/crypto/bcrypt"
	//"log"
)

// Login handles the login process for the user
func Login(w http.ResponseWriter, r *http.Request) {
	// If the http method is GET, then show the login page
	if r.Method == http.MethodGet {
		err := tpl.ExecuteTemplate(w, "login.html", nil)
		if err != nil {
			res := Error{}
			res.Code = 500
			res.Status = "Internal server error"
			ErrorHandler(w, r, &res)
			return
		}
	} else if r.Method == http.MethodPost {
		// If the method is POST, then take the data from the form and validate the user

		// Get the form data
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Check if the data is not empty
		if strings.Replace(email, " ", "", -1) == "" || strings.Replace(password, " ", "", -1) == "" {
			msg := PopUp{}
			msg.MessageEmpty = "Please fill the form!"
			err := tpl.ExecuteTemplate(w, "login.html", &msg)
			if err != nil {
				res := Error{}
				res.Code = 500
				res.Status = "Internal server error"
				ErrorHandler(w, r, &res)
				return
			}
			return
		}

		// Check if the entered email exists in the database
		emailExist, errMail := DB.CheckEmailExists(email)
		if errMail != nil {
			res := Error{}
			res.Code = 500
			res.Status = "Internal server error"
			ErrorHandler(w, r, &res)
			return
		}
		if emailExist {
			// If the email is found in the DB, then get the user and check passwords
			user, err := DB.GetUser(email)
			if err != nil {
				res := Error{}
				res.Code = 500
				res.Status = "Internal server error"
				ErrorHandler(w, r, &res)
				return
			}

			// Check the passwords
			passMatch := CheckPasswordHash(password, user.UserPassword)
			if !passMatch {
				// Show error message if passwords do not match
				msg := PopUp{}
				msg.MessageNoUser = "Wrong email or password!"
				err := tpl.ExecuteTemplate(w, "login.html", &msg)
				if err != nil {
					res := Error{}
					res.Code = 500
					res.Status = "Internal server error"
					ErrorHandler(w, r, &res)
					return
				}
				return
			} else {
				// If the email is found and the password matches, then set the session and redirect to home

				// Set session for the user
				SetSession(w, user.UserID)

				// Redirect the user to the home page
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}
		} else {
			msg := PopUp{}
			msg.MessageNoUser = "Wrong email or password!"
			err := tpl.ExecuteTemplate(w, "login.html", &msg)
			if err != nil {
				res := Error{}
				res.Code = 500
				res.Status = "Internal server error"
				ErrorHandler(w, r, &res)
				return
			}
		}
	}
}

// CheckPasswordHash checks if the given plain-text password matches the hashed password 
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil 
}

// SetSession sets a session cookie for the logged-in user
// func SetSession(w http.ResponseWriter, userID int) {
// 	sessionID := createSessionID()
// 	expiry := time.Now().Add(24 * time.Hour)

// 	_, err := db.Exec("INSERT INTO sessions (session_id, expiration, user_id) VALUES (?, ?, ?)", sessionID, expiry, userID)
// 	if err != nil {
// 		log.Println("Error creating session:", err)
// 		return
// 	}

// 	cookie := http.Cookie{
// 		Name:     "session_id",
// 		Value:    sessionID,
// 		Path:     "/",
// 		Expires:  expiry,
// 		HttpOnly: true,
// 	}
// 	http.SetCookie(w, &cookie)
// }

// // createSessionID creates a new session ID
// func createSessionID() string {
// 	b := make([]byte, 32)
// 	_, err := rand.Read(b)
// 	if err != nil {
// 		return ""
// 	}
// 	return base64.URLEncoding.EncodeToString(b)
// }

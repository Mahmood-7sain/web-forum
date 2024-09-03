package funcs

import (
	"fmt"
	"forum/DB"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	//If the http method is GET then show the login page only
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
		//If the method is POST then take the data from the form and validate and singup user

		//Get the form data
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		//Check if the data is not empty
		if strings.Replace(username, " ", "", -1) == "" || strings.Replace(username, " ", "", -1) == "" || strings.Replace(username, " ", "", -1) == "" {
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
		}

		//Check if the entered email and username are already taken or not
		emailTaken, errMail := DB.CheckEmailExists(email)
		userTaken, errName := DB.CheckUserExists(username)

		//Handle errors
		if errMail != nil || errName != nil {
			res := Error{}
			res.Code = 500
			res.Status = "Internal server error"
			ErrorHandler(w, r, &res)
			return
		}

		//Set up error msgs
		if emailTaken || userTaken || !IsStrongPassword(password) {
			msgs := PopUp{}
			if emailTaken {
				msgs.MessageMail = "Email already taken!"
			}
			if userTaken {
				msgs.MessageName = "Username already taken!"
			}
			if !IsStrongPassword(password) {
				msgs.MessagePass = "Min 8 chars, 1 Upper, 1 Lower, 1 number!"
			}
			err := tpl.ExecuteTemplate(w, "login.html", &msgs)
			if err != nil {
				res := Error{}
				res.Code = 500
				res.Status = "Internal server error"
				ErrorHandler(w, r, &res)
				return
			}
		} else {
			//If all of the data is valid then add user to the database and redirect to login page

			//Hash the user's password
			hashedPassword, err := HashPassword(password)
			if err != nil {
				res := Error{}
				res.Code = 500
				res.Status = "Internal server error"
				ErrorHandler(w, r, &res)
				return
			}

			//Create the user struct
			NewUser := DB.User{
				UserName:     username,
				UserEmail:    email,
				UserPassword: hashedPassword,
			}

			//Add user to the database and redirect to login
			errAdd := DB.AddUser(NewUser)
			if errAdd != nil {
				res := Error{}
				res.Code = 500
				res.Status = "Internal server error"
				ErrorHandler(w, r, &res)
				return
			} else {
				msg := PopUp{}
                msg.MessageSuccess = "Registered Successfully! Please login to continue."
				err := tpl.ExecuteTemplate(w, "login.html", &msg)
				if err != nil {
					fmt.Println(err)
					res := Error{}
					res.Code = 500
					res.Status = "Internal server error"
					ErrorHandler(w, r, &res)
					return
				}
			}
		}
	}
}

// HashPassword hashes a given password using bcrypt and returns the hashed password.
func HashPassword(password string) (string, error) {
	// Convert the password string to a byte slice
	passwordBytes := []byte(password)

	// Generate the bcrypt hash from the password with a cost of bcrypt.DefaultCost
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	// Convert the hashed password byte slice back to a string and return it
	return string(hashedPasswordBytes), nil
}

// IsStrongPassword checks if the given password is strong.
// A strong password must contain at least one uppercase letter, one lowercase letter, one digit, and be at least 8 characters long.
func IsStrongPassword(password string) bool {
	// Define the regex patterns for each requirement
	uppercasePattern := `[A-Z]`
	lowercasePattern := `[a-z]`
	digitPattern := `\d`
	minLengthPattern := `.{8,}`

	// Compile the regex patterns
	uppercaseRe := regexp.MustCompile(uppercasePattern)
	lowercaseRe := regexp.MustCompile(lowercasePattern)
	digitRe := regexp.MustCompile(digitPattern)
	minLengthRe := regexp.MustCompile(minLengthPattern)

	// Check each requirement
	return minLengthRe.MatchString(password) &&
		uppercaseRe.MatchString(password) &&
		lowercaseRe.MatchString(password) &&
		digitRe.MatchString(password)
}

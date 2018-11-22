package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
)

var baseUser = []User{
	User{
		Password: "12345678",
		Email:    "andres@hotmail.com",
	},
	User{
		Password: "987654321",
		Email:    "toni@hotmail.com",
	},
}

var emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type User struct {
	Email    string
	Password string
	Exito    string
	BadLogin string
	// Name        string
	// NewPassword string
	// Errors      map[string]string
}

// func (user User) Validate() bool {
// 	// if user.Password != baseUser.Password {
// 	// 	user.Errors["Password"] = "Password does not match"
// 	// }
// 	// if user.NewPassword == baseUser {
// 	// 	user.Errors["PasswordEqual"] = "Password is equal to new password"
// 	// }
// 	// if len(user.NewPassword) < 8 {
// 	// 	user.Errors["PasswordShort"] = "Very short password"
// 	// }
// 	// if user.Name == "" {
// 	// 	user.Errors["Name"] = "Please put your full name"
// 	// }

// 	// if !emailRegexp.MatchString(user.Email) {
// 	// 	user.Errors["Email"] = "Your email is invalid"
// 	// }

// 	return len(user.Errors) == 0
//}

func main() {
	http.HandleFunc("/signup", CreateFormHandler)
	http.HandleFunc("/validate", ValidateCredentialsHandler)
	http.HandleFunc("/login", CreateLoginHandler)
	http.HandleFunc("/validate/login/ok", ValidateLoginHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Fatal(http.ListenAndServe(":8084", nil))
}
func CreateFormHandler(w http.ResponseWriter, r *http.Request) {
	err := Render(w, "static/index.html.tmpl", User{})
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}
func CreateLoginHandler(w http.ResponseWriter, r *http.Request) {
	err := Render(w, "static/login.html.tmpl", User{})
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}

func ValidateCredentialsHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	user := User{
		Email:    r.Form.Get("Email"),
		Password: r.Form.Get("Password"),
	}
	for i := 0; i < len(baseUser); i++ {
		if user.Email == baseUser[i].Email {
			user.BadLogin = "Email already exits"
			Render(w, "static/index.html.tmpl", user)
			return
		}
	}
	user.Exito = "Credentials Save"
	baseUser = append(baseUser, user)
	Render(w, "static/index.html.tmpl", user)
}
func Validators(w http.ResponseWriter, r *http.Request) {
	f, err := os.OpenFile("users.json")
}
func ValidateLoginHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	user := User{
		Email:    r.Form.Get("Email"),
		Password: r.Form.Get("Password"),
	}

	for _, savedUser := range baseUser {
		if user.Email == savedUser.Email && user.Password == savedUser.Password {
			w.Write([]byte("Ok"))
			return
		}
	}
	user.BadLogin = "Bad Login"
	Render(w, "static/login.html.tmpl", user)
}

func Render(w http.ResponseWriter, tmpl string, information interface{}) error {
	template, err := template.ParseFiles(tmpl)
	if err != nil {
		return err
	}

	err = template.Execute(w, information)
	if err != nil {
		return err
	}

	return nil
}

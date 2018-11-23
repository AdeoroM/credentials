package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

var baseUser = []User{}

var emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type User struct {
	Email    string `json:"Email"`
	Password string `json:"Password"`
	Exito    string
	BadLogin string
}

func main() {
	file, err := os.Open("baseUsers.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(byteValue, &baseUser)

	http.HandleFunc("/signup", CreateFormHandler)
	http.HandleFunc("/validate", ValidateCredentialsHandler)
	http.HandleFunc("/login", CreateLoginHandler)
	http.HandleFunc("/validate/login/ok", ValidateLoginHandler)
	http.HandleFunc("/users", TableUsersHandler)
	http.HandleFunc("/users/delete", DeleteHandler)
	http.HandleFunc("/users/edit", EditUserHandler)
	http.HandleFunc("/users/update", ChangeUserHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Fatal(http.ListenAndServe(":8084", nil))
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

func CreateFormHandler(w http.ResponseWriter, r *http.Request) {
	err := Render(w, "static/index.html.tmpl", User{})
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

	baseUserJson, err := json.Marshal(baseUser)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	err = ioutil.WriteFile("baseUsers.json", baseUserJson, 0644)
	Render(w, "static/index.html.tmpl", user)
}

func TableUsersHandler(w http.ResponseWriter, r *http.Request) {
	Render(w, "static/list.html.tmpl", baseUser)
}

func CreateLoginHandler(w http.ResponseWriter, r *http.Request) {

	err := Render(w, "static/login.html.tmpl", User{})
	if err != nil {
		w.Write([]byte(err.Error()))
	}

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

func DeleteHandler(w http.ResponseWriter, r *http.Request) {

	email := r.URL.Query().Get("email")
	index := -1

	for i, user := range baseUser {
		if user.Email != email {
			continue
		}

		index = i
		break
	}

	if index >= 0 {
		baseUser = append(baseUser[:index], baseUser[index+1:]...)
	}
	baseUserJson, err := json.Marshal(baseUser)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	err = ioutil.WriteFile("baseUsers.json", baseUserJson, 0644)

	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

func EditUserHandler(w http.ResponseWriter, r *http.Request) {

	email := r.URL.Query().Get("email")
	u := User{}

	for _, user := range baseUser {
		if user.Email != email {
			continue
		}
		u = user
		break
	}

	if u.Email != "" {
		Render(w, "static/edit.html.tmpl", u)
	}

}

func ChangeUserHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	OriginalEmail := r.Form.Get("originalEmail")

	u := User{}
	index := -1

	for i, user := range baseUser {
		if user.Email != OriginalEmail {
			continue
		}
		u = user
		index = i
		break
	}

	if index >= 0 {
		u.Email = r.Form.Get("Email")
		u.Password = r.Form.Get("Password")
		baseUser[index] = u
	}

	baseUserJson, err := json.Marshal(baseUser)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	err = ioutil.WriteFile("baseUsers.json", baseUserJson, 0644)

	http.Redirect(w, r, "/users", http.StatusSeeOther)

}

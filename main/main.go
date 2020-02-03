package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// kode SALT
const saltnew string = "Test_GOl4NG_@_SaLT"

// User variabel user
type User struct {
	ID       int
	Name     string
	UserName string
	Email    string
	Address  string
	Password string
	Token    string
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "golang_test"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp(127.0.0.1:3306)/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

var tmpl = template.Must(template.ParseGlob("../form/*"))

// Index fungsi halaman awal list
func Index(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if len(token) > 10 {
		db := dbConn()
		defer db.Close()
		UserID := CheckToken(token, db)
		if UserID > 0 {
			selDB, err := db.Query("SELECT id,username,nama,email FROM user ORDER BY id DESC")
			if err != nil {
				panic(err.Error())
			}
			emp := User{}
			res := []User{}
			emp.Token = token
			for selDB.Next() {
				err = selDB.Scan(&emp.ID, &emp.UserName, &emp.Name, &emp.Email)
				if err != nil {
					panic(err.Error())
				}
				res = append(res, emp)
			}
			tmpl.ExecuteTemplate(w, "Index", res)
		} else {
			tmpl.ExecuteTemplate(w, "Login", nil)
		}
	} else {
		tmpl.ExecuteTemplate(w, "Login", nil)
	}

}

// Show fungsi untuk menampilkan data
func Show(w http.ResponseWriter, r *http.Request) {

	token := r.URL.Query().Get("token")
	if len(token) > 10 {
		db := dbConn()
		defer db.Close()
		UserID := CheckToken(token, db)
		if UserID > 0 {
			nID := r.URL.Query().Get("id")
			selDB, err := db.Query("SELECT id,username,nama,email FROM user WHERE id=?", nID)
			if err != nil {
				panic(err.Error())
			}
			emp := User{}
			emp.Token = token
			for selDB.Next() {
				err = selDB.Scan(&emp.ID, &emp.UserName, &emp.Name, &emp.Email)
				if err != nil {
					panic(err.Error())
				}
			}
			tmpl.ExecuteTemplate(w, "Show", emp)
		} else {
			tmpl.ExecuteTemplate(w, "Login", nil)
		}
	} else {
		tmpl.ExecuteTemplate(w, "Login", nil)
	}

}

// New fungsi untuk tambah
func New(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if len(token) > 10 {
		db := dbConn()
		defer db.Close()
		UserID := CheckToken(token, db)
		if UserID > 0 {
			emp := User{}
			emp.Token = token
			tmpl.ExecuteTemplate(w, "New", emp)
		} else {
			tmpl.ExecuteTemplate(w, "Login", nil)
		}
	} else {
		tmpl.ExecuteTemplate(w, "Login", nil)
	}
}

// Edit fungsi untuk rubah
func Edit(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if len(token) > 10 {
		db := dbConn()
		defer db.Close()
		UserID := CheckToken(token, db)
		if UserID > 0 {
			nID := r.URL.Query().Get("id")
			selDB, err := db.Query("SELECT id,username,nama,email FROM user WHERE id=?", nID)
			if err != nil {
				panic(err.Error())
			}
			emp := User{}
			emp.Token = token
			for selDB.Next() {
				err = selDB.Scan(&emp.ID, &emp.UserName, &emp.Name, &emp.Email)
				if err != nil {
					panic(err.Error())
				}
			}
			tmpl.ExecuteTemplate(w, "Edit", emp)
		} else {
			tmpl.ExecuteTemplate(w, "Login", nil)
		}
	} else {
		tmpl.ExecuteTemplate(w, "Login", nil)
	}

}

// Insert fungsi untuk tambah
func Insert(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if len(token) > 10 {
		db := dbConn()
		defer db.Close()
		UserID := CheckToken(token, db)
		if UserID > 0 {
			if r.Method == "POST" {
				name := r.FormValue("name")
				email := r.FormValue("email")
				username := r.FormValue("username")
				password := r.FormValue("password")
				insForm, err := db.Prepare("INSERT INTO user(username, nama,email,password) VALUES(?,?,?,MD5(?))")
				if err != nil {
					panic(err.Error())
				}
				insForm.Exec(username, name, email, password)
				log.Println("INSERT: UserName: " + username + " | Name: " + name)
			}
			http.Redirect(w, r, "/?token="+token, 301)
		} else {
			tmpl.ExecuteTemplate(w, "Login", nil)
		}
	} else {
		tmpl.ExecuteTemplate(w, "Login", nil)
	}

}

// Update fungsi untuk update
func Update(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if len(token) > 10 {
		db := dbConn()
		defer db.Close()
		UserID := CheckToken(token, db)
		if UserID > 0 {
			if r.Method == "POST" {
				name := r.FormValue("name")
				username := r.FormValue("username")
				email := r.FormValue("email")
				id := r.FormValue("uid")
				insForm, err := db.Prepare("UPDATE user SET username=?,nama=?,email=? WHERE id=?")
				if err != nil {
					panic(err.Error())
				}
				insForm.Exec(username, name, email, id)
				log.Println("UPDATE: Name: " + name + " | Email: " + email)
			}
			http.Redirect(w, r, "/?token="+token, 301)
		} else {
			tmpl.ExecuteTemplate(w, "Login", nil)
		}
	} else {
		tmpl.ExecuteTemplate(w, "Login", nil)
	}
}

// Delete fungsi untuk hapus
func Delete(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if len(token) > 10 {
		db := dbConn()
		defer db.Close()
		UserID := CheckToken(token, db)
		if UserID > 0 {
			emp := r.URL.Query().Get("id")
			hasilint, _ := strconv.Atoi(emp)
			if hasilint > 0 && hasilint != UserID {
				// memastikan tidak delete diri sendiri
				delForm, err := db.Prepare("DELETE FROM user WHERE id=?")
				if err != nil {
					panic(err.Error())
				}
				delForm.Exec(emp)
				log.Println("DELETE")
			}

			http.Redirect(w, r, "/?token="+token, 301)
		} else {

			tmpl.ExecuteTemplate(w, "Login", nil)
		}
	} else {

		tmpl.ExecuteTemplate(w, "Login", nil)
	}
}

// CheckToken fungsi untuk cek token
func CheckToken(token string, db *sql.DB) int {
	UserID := 0
	if len(token) > 10 {
		db.QueryRow(`SELECT id FROM user where token=?`, token).
			Scan(
				&UserID,
			)
	}
	return UserID
}

// ActionLogin funtion untuk cek login
func ActionLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		username := r.FormValue("username")
		password := r.FormValue("password")

		if len(username) > 0 && len(password) > 0 {
			UserID := 0
			db := dbConn()
			defer db.Close()
			db.QueryRow(`SELECT id FROM user where username=? AND password=MD5(?)`, username, password).
				Scan(
					&UserID,
				)
			if UserID > 0 {
				timeInt := time.Now().UnixNano() / 1000000000
				hashdata := sha256.Sum256([]byte(fmt.Sprintf("%v", timeInt) + saltnew))
				tokenTeks := hex.EncodeToString(hashdata[:])
				insForm, err := db.Prepare("UPDATE user SET token=? WHERE id=?")
				if err != nil {
					panic(err.Error())
				}
				insForm.Exec(tokenTeks, UserID)
				http.Redirect(w, r, "/?token="+tokenTeks, 301)
			} else {
				tmpl.ExecuteTemplate(w, "Login", nil)
			}
		} else {
			tmpl.ExecuteTemplate(w, "Login", nil)
		}
	}
}

// Login fungsi untuk login
func Login(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if len(token) > 10 {
		db := dbConn()
		defer db.Close()
		UserID := CheckToken(token, db)

		if UserID > 0 {
			Index(w, r)
		} else {
			tmpl.ExecuteTemplate(w, "Login", nil)
		}

	} else {
		tmpl.ExecuteTemplate(w, "Login", nil)
	}
}

// main fungsi yang pertama kali di eksekusi
func main() {
	http.HandleFunc("/", Login)
	http.HandleFunc("/actlogin", ActionLogin)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)
	http.HandleFunc("/logout", Index)
	log.Println("Server started on: http://localhost:8181")
	http.ListenAndServe(":8181", nil)

}

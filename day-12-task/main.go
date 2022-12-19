package main

import (
	"context"
	"day-12/connection"
	"day-12/middleware"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

// struktur object

// const SESSION_ID = "id"

type Metadata struct{
	Title string
	IsLogin bool
	UserName string
	FlashData string
	UserID    int
}

var Data = Metadata{
	Title: "Personal Web",
}

type Blog struct {
	Id int
	Title string
	Description string
	Start_date, End_date time.Time
	Format_start	string
	Technologies [4]string
	Image string
	Author string
	IsLogin bool
}

type User struct {
	Id       int
	Name     string
	Email    string
	Password string
}


// dummy data

var Blogs = []Blog{
}

func main() {
	// assign mux kedalam variable route
	route := mux.NewRouter()

	// database
	connection.DatabaseConnect()

	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	route.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads/"))))

	// ROUTE >>>
	route.HandleFunc("/", home).Methods("GET") //Oke
	route.HandleFunc("/contact", contact).Methods("GET")

	// untuk blog

	route.HandleFunc("/blog", blogIndex).Methods("GET") //Oke
	route.HandleFunc("/blog/{id}", blogShow).Methods("GET") //Oke
	route.HandleFunc("/blog-create", blogCreate).Methods("GET") //Oke
	// route.HandleFunc("/blog-store", blogStore).Methods("POST") //Oke
	route.HandleFunc("/blog-store", middleware.UploadFile(blogStore)).Methods("POST")
	route.HandleFunc("/blog-edit/{id}", blogEdit).Methods("GET")
	route.HandleFunc("/blog-update/{id}", blogUpdate).Methods("POST")
	route.HandleFunc("/blog-delete/{id}", blogDelete).Methods("GET") //Oke

	//untuk register 

	route.HandleFunc("/register", register).Methods("GET")
	route.HandleFunc("/register-form", registerForm).Methods("POST")

	// untuk login
	route.HandleFunc("/login", login).Methods("GET")
	route.HandleFunc("/login-form", loginForm).Methods("POST")

	// untuk logout
	route.HandleFunc("/logout", logout).Methods("GET")

	// <<< ROUTE

	fmt.Println("Server running on port 5000")
	// fmt.Println(Blogs)
	http.ListenAndServe("localhost:5000", route)

	
}

func home(w http.ResponseWriter, r *http.Request) {
	// setting untuk doc text/html
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// setting route mengarah file html
	var tmpl, err = template.ParseFiles("views/index.html")
	// handle route error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}
	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.UserName = session.Values["Name"].(string)
	}

	fm := session.Flashes("message")

	var flashes []string
	if len(fm) > 0 {
		session.Save(r, w)
		// Initiate a strings slice to return messages.
		for _, fl := range fm {
			// Add message to the slice.
			flashes = append(flashes, fl.(string))
		}
	}

	Data.FlashData = strings.Join(flashes, "")

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

func contact(w http.ResponseWriter, r *http.Request) {
	// setting untuk doc text/html
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// setting route mengarah file html
	var tmpl, err = template.ParseFiles("views/contact.html")

	// handle route error

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	// session
	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.UserName = session.Values["Name"].(string)
	}
	// session

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

// 
func blogIndex(w http.ResponseWriter, r *http.Request) {
		// setting untuk doc text/html
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		// setting route mengarah file html
		var tmpl, err = template.ParseFiles("views/blog/blog-index.html")
		// handle route error
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("message : " + err.Error()))
			return
		}
		// session

		var store = sessions.NewCookieStore([]byte("SESSION_ID"))
		session, _ := store.Get(r, "SESSION_ID")

		if session.Values["IsLogin"] != true {
			Data.IsLogin = false
		} else {
			Data.UserID = session.Values["Id"].(int)
			Data.IsLogin = session.Values["IsLogin"].(bool)
			Data.UserName = session.Values["Name"].(string)
		}	

		// session

		rows, _ := connection.Conn.Query(context.Background(), `SELECT id, title, start_date, end_date, description, technologies, author, image FROM tb_blog`)
		// LEFT JOIN tb_user ON tb_blog.authorid = tb_user.id
		var result []Blog // array data
		
		for rows.Next() {
			var each = Blog{} // manggil struct
	
			err := rows.Scan(&each.Id, &each.Title, &each.Start_date, &each.End_date, &each.Description, &each.Technologies, &each.Author, &each.Image)

			if err != nil {
				fmt.Println(err.Error())
				return
			}
			// each.Author = "Dicky"
			each.Format_start = each.Start_date.Format("2 January 2006")
			// each.Format_start = each.Start_date.Format("2 January 2006")
			result = append(result, each)
		}
		respData := map[string]interface{}{
			"Blogs": result,
			"Data": Data,
		}
		
		// fmt.Println(respData)
		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, respData)
}

func blogShow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/blog/blog-detail.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	var BlogDetail = Blog{
	}

	// s := BlogDetail.Start_date.Format("2 January 2006")
	// e :=BlogDetail.End_date.Format("2 January 2006")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	err = connection.Conn.QueryRow(context.Background(), "SELECT id, title, start_date, end_date,  description , technologies , author , image FROM tb_blog WHERE id=$1", id).Scan(
		&BlogDetail.Id, &BlogDetail.Title, &BlogDetail.Start_date, &BlogDetail.End_date, &BlogDetail.Description, &BlogDetail.Technologies , &BlogDetail.Author, &BlogDetail.Image,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}
	
	// session
	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")
	
	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
		} else {
			Data.IsLogin = session.Values["IsLogin"].(bool)
			Data.UserName = session.Values["Name"].(string)
		}
		// session
		
	data := map[string]interface{}{
		"Blog": BlogDetail,
		"Data": Data,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, data)
}

func blogCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/blog/blog-create.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}
	
	// session
	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.UserID = session.Values["Id"].(int)
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.UserName = session.Values["Name"].(string)
	}
	// session
	fmt.Println(Data.UserName)
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

func blogStore(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	Node_js:= r.PostForm.Get("inlineCheckbox1") 
	React_js := r.PostForm.Get("inlineCheckbox2")
	Next_js := r.PostForm.Get("inlineCheckbox3")
	Type_Script := r.PostForm.Get("inlineCheckbox4")
	Title :=       r.PostForm.Get("inputTitle")
	Description := r.PostForm.Get("inputDescription")
	Start_date :=  r.PostForm.Get("startDate")
	End_date :=    r.PostForm.Get("endDate")
	// Author :=      "Dicky Joel Saputra"
	Technologies := [4]string{Node_js,React_js,Next_js,Type_Script}
	// Blogs.push(newBlog)

	dataContex := r.Context().Value("dataFile")
	image := dataContex.(string)

	// session
	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.UserID = session.Values["Id"].(int)
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.UserName = session.Values["Name"].(string)
	}
	// session
	
	_, err = connection.Conn.Exec(context.Background(),"INSERT into tb_blog(title,authorid,description,start_date,end_date,technologies,author,image) VALUES ($1, $2, $3 , $4 , $5, $6, $7, $8)", Title,Data.UserID,Description,Start_date,End_date,Technologies,Data.UserName,image)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/blog", http.StatusMovedPermanently)
}

func blogEdit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/blog/blog-edit.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	var BlogDetail = Blog{}
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	err = connection.Conn.QueryRow(context.Background(), "SELECT id, title, start_date, end_date, description , technologies FROM tb_blog WHERE id=$1", id).Scan(
		&BlogDetail.Id, &BlogDetail.Title, &BlogDetail.Start_date, &BlogDetail.End_date, &BlogDetail.Description, &BlogDetail.Technologies ,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	// session
	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.UserName = session.Values["Name"].(string)
	}
	// session

	data := map[string]interface{}{
		"Blog": BlogDetail,
		"Data": Data,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, data)
}

func blogUpdate(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	
	Node_js:= r.PostForm.Get("inlineCheckbox1") 
	React_js := r.PostForm.Get("inlineCheckbox2")
	Next_js := r.PostForm.Get("inlineCheckbox3")
	Type_Script := r.PostForm.Get("inlineCheckbox4")
	title :=       r.PostForm.Get("inputTitle")
	description := r.PostForm.Get("inputDescription")
	start_date :=  r.PostForm.Get("startDate")
	end_date :=    r.PostForm.Get("endDate")
	// Author :=      "Dicky Joel Saputra"
	Technologies := [4]string{Node_js,React_js,Next_js,Type_Script}
	
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	
	_, err = connection.Conn.Exec(context.Background(), 
	`UPDATE public.tb_blog SET "title"=$1, "start_date"=$2, "end_date"=$3, "description"=$4, "technologies"=$5 WHERE "id"=$6`, title,start_date,end_date,description,Technologies,id)

	// `UPDATE public.tb_blog SET "title"=$1, "description"=$2, "start_date"=$3, "end_date"=$4, "technologies"=$5, WHERE "id"=$6;`, Title, Description, Start_date, End_date, Technologies, id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	fmt.Print()
	http.Redirect(w, r, "/blog", http.StatusMovedPermanently)
}

func blogDelete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	// fmt.Println(id)

	// Blogs = append(Blogs[:index], Blogs[index+1:]...)
	// fmt.Println(Blogs)

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_blog WHERE id=$1", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/blog", http.StatusFound)
}

func register(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/blog/form-register.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

func registerForm(w http.ResponseWriter, r *http.Request){
	err := r.ParseForm()
	if err != nil{
		log.Fatal(err)
	}

	var name = r.PostForm.Get("inputName")
	var email = r.PostForm.Get("inputEmail")
	var password = r.PostForm.Get("inputPassword")

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	// fmt.Println(passwordHash)

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_user(name, email, password) VALUES($1, $2, $3)", name, email, passwordHash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	session.AddFlash("succesfull register", "message")

	session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
}

func login(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/blog/form-login.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	fm := session.Flashes("message")

	var flashes []string
	if len(fm) > 0 {
		session.Save(r, w)
		// Initiate a strings slice to return messages.
		for _, fl := range fm {
			// Add message to the slice.
			flashes = append(flashes, fl.(string))
		}
	}

	Data.FlashData = strings.Join(flashes, "")

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

func loginForm(w http.ResponseWriter, r *http.Request){
	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	email := r.PostForm.Get("inputEmail")
	password := r.PostForm.Get("inputPassword")

	user := User{}

	err = connection.Conn.QueryRow(context.Background(), "SELECT * FROM tb_user WHERE email=$1", email).Scan(
		&user.Id, &user.Name, &user.Email, &user.Password,
	)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	session.Values["IsLogin"] = true
	session.Values["Name"] = user.Name
	session.Values["Id"] = user.Id
	session.Options.MaxAge = 10800 // 3 hours

	session.AddFlash("Successfully Login!", "message")
	session.Save(r, w)

	// fmt.Println(user)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func logout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Logout!")
	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")
	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
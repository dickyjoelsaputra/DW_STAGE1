package main

import (
	"context"
	"day-10/connection"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// struktur object

type Blog struct {Title, Author, Description string
	Start_date, End_date time.Time
	Technologies [4]string
	Id int
}

// dummy data

var Blogs = []Blog{
	// {
	// 	Id: 0,
	// 	Title:       "Pasar koding menjamur",
	// 	Start_date:  "2006-01-02",
	// 	End_date:    "2006-01-02",
	// 	Author:      "Dicky Joel Saputra",
	// 	Description: "lorem lorem lorem",
	// 	Technologies: [4]string{"checked", "checked", "checked", "checked"},
	// },
	// {
	// 	Id: 1,
	// 	Title:       "Padasdsadasdsadadsadsadas",
	// 	Start_date:	 "2006-01-02",
	// 	End_date:    "2006-01-02",
	// 	Author:      "Dicky Joel Saputra",
	// 	Description: "lodqdqiwqqdkjdwqdwqqwqdjwqndjwqnjdwqnj",
	// 	Technologies: [4]string{"", "checked", "checked", ""},
	// },
	// {
	// 	Id: 2,
	// 	Title:       "asdasdsadasdsadsads",
	// 	Start_date:	 "2011-01-02",
	// 	End_date:    "2020-01-02",
	// 	Author:      "Dicky Joel Dickkk",
	// 	Description: "dsadsadsasadsadsa",
	// 	Technologies: [4]string{"checked", "checked", "checked", ""},
	// },
}

func main() {
	// assign mux kedalam variable route
	route := mux.NewRouter()

	// database
	connection.DatabaseConnect()

	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	// ROUTE >>>
	route.HandleFunc("/", home).Methods("GET") //Oke
	route.HandleFunc("/contact", contact).Methods("GET")

	// untuk blog

	route.HandleFunc("/blog", blogIndex).Methods("GET") //Oke
	route.HandleFunc("/blog/{id}", blogShow).Methods("GET") //Oke
	route.HandleFunc("/blog-create", blogCreate).Methods("GET") //Oke
	route.HandleFunc("/blog-store", blogStore).Methods("POST") //Oke
	route.HandleFunc("/blog-edit/{id}", blogEdit).Methods("GET")
	route.HandleFunc("/blog-update/{id}", blogUpdate).Methods("POST")
	route.HandleFunc("/blog-delete/{id}", blogDelete).Methods("GET") //Oke

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

	rows, _ := connection.Conn.Query(context.Background(), "SELECT id, title, start_date, end_date, description, technologies FROM tb_blog")

	var result []Blog // array data

	for rows.Next() {
		var each = Blog{} // manggil struct
		each.Start_date.Format("2006-01-02")
		each.End_date.Format("2006-01-02")
		// fmt.Println(each)

		err := rows.Scan(&each.Id, &each.Title, &each.Start_date, &each.End_date, &each.Description, &each.Technologies ,)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// fmt.Println(each)
		result = append(result, each)
		// fmt.Println(result)
	}

	respData := map[string]interface{}{
		"Blogs": result,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, respData)
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

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
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
	
		rows, _ := connection.Conn.Query(context.Background(), "SELECT id, title, start_date, end_date, description, technologies FROM tb_blog")

		var result []Blog // array data
	
		for rows.Next() {
			var each = Blog{} // manggil struct
	
			err := rows.Scan(&each.Id, &each.Title, &each.Start_date, &each.End_date, &each.Description, &each.Technologies ,)
			each.Start_date.Format("2006-01-02")
			each.End_date.Format("2006-01-02")
			if err != nil {
				fmt.Println(err.Error())
				return
			}
	
			// each.Author = "Dicky"
			// each.Format_date = each.Post_date.Format("2 January 2006")
	
			result = append(result, each)
		}
	
		respData := map[string]interface{}{
			"Blogs": result,
		}
	
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

	var BlogDetail = Blog{}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	err = connection.Conn.QueryRow(context.Background(), "SELECT id, title, start_date, end_date,  description , technologies FROM tb_blog WHERE id=$1", id).Scan(
		&BlogDetail.Id, &BlogDetail.Title, &BlogDetail.Start_date, &BlogDetail.End_date, &BlogDetail.Description, &BlogDetail.Technologies ,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	BlogDetail.Author = "Dicky Joel"

	data := map[string]interface{}{
		"Blog": BlogDetail,
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

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
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

	fmt.Println(Start_date,End_date)

	_, err = connection.Conn.Exec(context.Background(),"INSERT into tb_blog(title,description,start_date,end_date,technologies) VALUES ($1, $2, $3 , $4 , $5)", Title,Description,Start_date,End_date,Technologies)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	fmt.Println(err)

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

	data := map[string]interface{}{
		"Blog": BlogDetail,
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
	
	title := r.PostForm.Get("inputTitle")
	description := r.PostForm.Get("inputDescription")
	// startdate := r.PostForm.Get("startDate")
	// enddate:= r.PostForm.Get("endDate")
	// Node_js:= r.PostForm.Get("inlineCheckbox1") 
	// React_js := r.PostForm.Get("inlineCheckbox2")
	// Next_js := r.PostForm.Get("inlineCheckbox3")
	// Type_Script := r.PostForm.Get("inlineCheckbox4")
	
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	
	_, err = connection.Conn.Exec(context.Background(), `UPDATE public.tb_blog
	SET "title"=$1, "description"=$2 WHERE "id"=$3`, title,description,id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/blog", http.StatusMovedPermanently)
}

func blogDelete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	fmt.Println(id)

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

// func date(){

// }
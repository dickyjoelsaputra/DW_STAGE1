package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// struktur object

type Blog struct {Title, Author, Description, Start_date, End_date string
	Technologies [4]string
}

// dummy data

var Blogs = []Blog{
	{
		Title:       "Pasar koding menjamur",
		Start_date:  "2006-01-02",
		End_date:    "2006-01-02",
		Author:      "Dicky Joel Saputra",
		Description: "lorem lorem lorem",
		Technologies: [4]string{"checked", "checked", "checked", "checked"},
	},
	{
		Title:       "Padasdsadasdsadadsadsadas",
		Start_date:	 "2006-01-02",
		End_date:    "2006-01-02",
		Author:      "Dicky Joel Saputra",
		Description: "lodqdqiwqqdkjdwqdwqqwqdjwqndjwqnjdwqnj",
		Technologies: [4]string{"", "checked", "checked", ""},
	},
	{
		Title:       "asdasdsadasdsadsads",
		Start_date:	 "2011-01-02",
		End_date:    "2020-01-02",
		Author:      "Dicky Joel Dickkk",
		Description: "dsadsadsasadsadsa",
		Technologies: [4]string{"checked", "checked", "checked", ""},
	},
}

func main() {
	// assign mux kedalam variable route
	route := mux.NewRouter()

	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	// ROUTE >>>
	route.HandleFunc("/", home).Methods("GET")
	route.HandleFunc("/contact", contact).Methods("GET")

	// untuk blog

	route.HandleFunc("/blog", blogIndex).Methods("GET")
	route.HandleFunc("/blog/{index}", blogShow).Methods("GET")
	route.HandleFunc("/blog-create", blogCreate).Methods("GET")
	route.HandleFunc("/blog-store", blogStore).Methods("POST")
	route.HandleFunc("/blog-edit/{index}", blogEdit).Methods("GET")
	route.HandleFunc("/blog-update/{index}", blogUpdate).Methods("POST")
	route.HandleFunc("/blog-delete/{index}", blogDelete).Methods("GET")

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

	respData := map[string]interface{}{
		"Blogs": Blogs,
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

	respData := map[string]interface{}{
		"Blogs": Blogs,
	}

	fmt.Println(respData)

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

	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	for i, data := range Blogs {
		if index == i {
			BlogDetail = Blog{
				Title:       data.Title,
				Description: data.Description,
				Start_date:  data.Start_date,
				End_date:    data.End_date,
				Author:      "Dicky Joel Saputra",
			}
		}
	}

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

	// title := r.PostForm.Get("inputTitle")
	// description := r.PostForm.Get("inputDescription")
	// startdate := r.PostForm.Get("startDate")
	// Node_js := r.PostForm.Get("inlineCheckbox1")
	// fmt.Println(Node_js)

	Node_js:= r.PostForm.Get("inlineCheckbox1") 
	React_js := r.PostForm.Get("inlineCheckbox2")
	Next_js := r.PostForm.Get("inlineCheckbox3")
	Type_Script := r.PostForm.Get("inlineCheckbox4")

	var newBlog = Blog{
		Title:       r.PostForm.Get("inputTitle"),
		Description: r.PostForm.Get("inputDescription"),
		Start_date:  r.PostForm.Get("startDate"),
		End_date:    r.PostForm.Get("endDate"),
		Author:      "Dicky Joel Saputra",
		Technologies: [4]string{Node_js,React_js,Next_js,Type_Script},
	}

	// Blogs.push(newBlog)

	Blogs = append(Blogs, newBlog)

	fmt.Println(Blogs)

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

	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	for i, data := range Blogs {
		if index == i {
			BlogDetail = Blog{
				Title:       data.Title,
				Description: data.Description,
				Start_date:  data.Start_date,
				End_date:    data.End_date,
				Author:      "Dicky Joel Saputra",
				Technologies: data.Technologies,
			}
		}
	}

	data := map[string]interface{}{
		"Blog": BlogDetail,
		"Index": index,
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
	startdate := r.PostForm.Get("startDate")
	enddate:= r.PostForm.Get("endDate")
	Node_js:= r.PostForm.Get("inlineCheckbox1") 
	React_js := r.PostForm.Get("inlineCheckbox2")
	Next_js := r.PostForm.Get("inlineCheckbox3")
	Type_Script := r.PostForm.Get("inlineCheckbox4")

	BlogDetail := Blog{				
		Title:       title,
		Description: description,
		Start_date:  startdate,
		End_date:    enddate,
		Author:      "Dicky Joel Saputra",
		Technologies: [4]string{Node_js,React_js,Next_js,Type_Script},
	}
	index, _ := strconv.Atoi(mux.Vars(r)["index"])


	for i, _ := range Blogs {
		if index == i {
			Blogs[index] = BlogDetail
		}
	}


	http.Redirect(w, r, "/blog", http.StatusMovedPermanently)
}

func blogDelete(w http.ResponseWriter, r *http.Request) {
	index, _ := strconv.Atoi(mux.Vars(r)["index"])
	fmt.Println(index)

	Blogs = append(Blogs[:index], Blogs[index+1:]...)

	http.Redirect(w, r, "/blog", http.StatusFound)
}

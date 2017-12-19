package main

import (
	"database/sql"
	"html/template"
	"fmt"
	"net/http"

    _ "github.com/go-sql-driver/mysql"
    "github.com/gorilla/mux"
)

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/view/", viewHandler)
	r.HandleFunc("/home/", HomeHandler)
	r.HandleFunc("/edit/{id:[0-9]+}", editHandler)
    r.HandleFunc("/editt/{id:[0-9]+}", edittHandler)
	r.HandleFunc("/save/", saveHandler)
    r.HandleFunc("/simpan/", simpanHandler)
    r.HandleFunc("/update/", updateHandler)
    r.HandleFunc("/updatee/", updateeHandler)
    r.HandleFunc("/delete/{id:[0-9]+}", deleteHandler)
    r.HandleFunc("/hapus/{id:[0-9]+}", hapusHandler)
    r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
    http.ListenAndServe(":8080", r)
}

var templates = template.Must(template.ParseFiles("templates/edit.html", "templates/editt.html", "templates/view.html", "templates/home.html"))

type Page struct {
	Title string
	Body  []byte
}
type Pagedata struct {
	// Title string
	Datasql []Entry
    Userdata []Entry
}

type Entry struct {
    Number int
    ID, Name, Nama, Umur, Email, Address string
}



func renderTemplate(w http.ResponseWriter, tmpl string, p *Pagedata) {
	// p.Datasql = template.HTML()

    err := templates.ExecuteTemplate(w, tmpl+".html", p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func editHandler(w http.ResponseWriter, r *http.Request) {
    
    vars := mux.Vars(r)
    fmt.Println(vars["id"]);
    db, err := sql.Open("mysql", "root@tcp(localhost:3306)/staff") //localhost:3306
    if err != nil {
        panic(err.Error())  // Just for example purpose. You should use proper error handling instead of panic
    }
    defer db.Close()

    //SELECT DATA BY ID
    users, err := db.Query("SELECT * FROM users WHERE id="+vars["id"])
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    defer users.Close()

    results := []Entry{}
    setData := Entry{}

    for users.Next() {
        var id, name, nama, email, address string
        users.Scan(&id, &name, &nama, &email, &address)
        setData.ID = id
        setData.Name = name
        setData.Nama = nama
        setData.Email = email
        setData.Address = address
        results = append(results, setData)
    }

    fmt.Println(results)


    //SELECT ALL DATA
    rows, err := db.Query("SELECT * FROM users")
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    defer rows.Close()

    // Fetch rows

    resultsRow := []Entry{}
    setDataRow := Entry{}
    i := 0
    for rows.Next() {
        var id, name, nama, email, address string
        rows.Scan(&id, &name, &nama, &email, &address)
        i += 1
        setDataRow.Number = i
        setDataRow.ID = id
        setDataRow.Name = name
        setDataRow.Nama = nama
        setDataRow.Email = email
        setDataRow.Address = address
        resultsRow = append(resultsRow, setDataRow)
    }

    if err = rows.Err(); err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    p := &Pagedata{Datasql: resultsRow, Userdata: results}
    w.Header().Set("Content-Type", "text/html")
    renderTemplate(w, "edit", p)
}

func edittHandler(w http.ResponseWriter, r *http.Request) {

    vars := mux.Vars(r)
    fmt.Println(vars["id"]);
    db, err := sql.Open("mysql", "root@tcp(localhost:3306)/staff") //localhost:3306
    if err != nil {
        panic(err.Error())  // Just for example purpose. You should use proper error handling instead of panic
    }
    defer db.Close()

    //SELECT DATA BY ID
    man, err := db.Query("SELECT * FROM man WHERE id="+vars["id"])
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    defer man.Close()

    results := []Entry{}
    setData := Entry{}

    for man.Next() {
        var id, name, age string
        man.Scan(&id, &name, &age)
        setData.ID = id
        setData.Nama = name
        setData.Umur = age
        results = append(results, setData)
    }

    fmt.Println(results)


    //SELECT ALL DATA
    rows, err := db.Query("SELECT * FROM man")
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    defer rows.Close()

    // Fetch rows

    resultsRow := []Entry{}
    setDataRow := Entry{}
    i := 0
    for rows.Next() {
        var id, name, age string
        rows.Scan(&id, &name, &age)
        i += 1
        setDataRow.Number = i
        setDataRow.ID = id
        setDataRow.Nama = name
        setDataRow.Umur = age
        resultsRow = append(resultsRow, setDataRow)
    }

    if err = rows.Err(); err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    p := &Pagedata{Datasql: resultsRow, Userdata: results}
    w.Header().Set("Content-Type", "text/html")
    renderTemplate(w, "editt", p)
}

func updateeHandler(w http.ResponseWriter, r *http.Request){
    id := r.FormValue("id")
    name := r.FormValue("nama")
    age := r.FormValue("umur")

    db, err := sql.Open("mysql", "root@tcp(localhost:3306)/staff") //localhost:3306
    if err != nil {
        panic(err.Error())  // Just for example purpose. You should use proper error handling instead of panic
    }
    defer db.Close()

    // Execute the query
    stmtIns, err := db.Query("UPDATE man SET name='"+name+"', age='"+age+"' WHERE id="+id)

    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    defer stmtIns.Close()

    http.Redirect(w, r, "/home/", http.StatusFound)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {

    db, err := sql.Open("mysql", "root@tcp(localhost:3306)/staff") //localhost:3306
    if err != nil {
        panic(err.Error())  // Just for example purpose. You should use proper error handling instead of panic
    }
    defer db.Close()

    rows, err := db.Query("SELECT * FROM man")
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }

    // text := ""
    resultsRow := []Entry{}
    setDataRow := Entry{}
    i := 0
    for rows.Next() {
        var id, nama, age string
        rows.Scan(&id, &nama, &age)
        i += 1
        setDataRow.Number = i
        setDataRow.ID = id
        setDataRow.Nama = nama
        setDataRow.Umur = age
        resultsRow = append(resultsRow, setDataRow)
    }

    if err = rows.Err(); err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    p := &Pagedata{Datasql: resultsRow}
    w.Header().Set("Content-Type", "text/html")
    renderTemplate(w, "home", p)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	
	db, err := sql.Open("mysql", "root@tcp(localhost:3306)/staff") //localhost:3306
    if err != nil {
        panic(err.Error())  // Just for example purpose. You should use proper error handling instead of panic
    }
    defer db.Close()

	rows, err := db.Query("SELECT * FROM users")
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }

    // text := ""
    resultsRow := []Entry{}
    setDataRow := Entry{}
    i := 0
    for rows.Next() {
        var id, name, nama, email, address string
        rows.Scan(&id, &name, &nama, &email, &address)
        i += 1
        setDataRow.Number = i
        setDataRow.ID = id
        setDataRow.Name = name
        setDataRow.Nama = nama
        setDataRow.Email = email
        setDataRow.Address = address
        resultsRow = append(resultsRow, setDataRow)
    }

    if err = rows.Err(); err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
	p := &Pagedata{Datasql: resultsRow}
	w.Header().Set("Content-Type", "text/html")
	renderTemplate(w, "view", p)
}

func updateHandler(w http.ResponseWriter, r *http.Request){
    id := r.FormValue("id")
    name := r.FormValue("name")
    nama := r.FormValue("nama")
    email := r.FormValue("email")
    address := r.FormValue("address")

    db, err := sql.Open("mysql", "root@tcp(localhost:3306)/staff") //localhost:3306
    if err != nil {
        panic(err.Error())  // Just for example purpose. You should use proper error handling instead of panic
    }
    defer db.Close()

        // Execute the query
    stmtIns, err := db.Query("UPDATE users SET name='"+name+"', nama='"+nama+"', email='"+email+"', address='"+address+"' WHERE id="+id)
    
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    defer stmtIns.Close()

    http.Redirect(w, r, "/view/", http.StatusFound)
}


func simpanHandler(w http.ResponseWriter, r *http.Request) {

    name := r.FormValue("nama")
    age := r.FormValue("umur")

    db, err := sql.Open("mysql", "root@tcp(localhost:3306)/staff") //localhost:3306
    if err != nil {
        panic(err.Error())  // Just for example purpose. You should use proper error handling instead of panic
    }
    defer db.Close()

    // Execute the query
    stmtIns, err := db.Query("INSERT INTO man (name, age) VALUES ('"+name+"','"+age+"')")

    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    defer stmtIns.Close()

    http.Redirect(w, r, "/home/", http.StatusFound)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {

    name := r.FormValue("name")
    nama := r.FormValue("nama")
    email := r.FormValue("email")
    address := r.FormValue("address")

    db, err := sql.Open("mysql", "root@tcp(localhost:3306)/staff") //localhost:3306
    if err != nil {
        panic(err.Error())  // Just for example purpose. You should use proper error handling instead of panic
    }
    defer db.Close()

        // Execute the query
    stmtIns, err := db.Query("INSERT INTO users (name, nama, email, address) VALUES ('"+name+"','"+nama+"','"+email+"','"+address+"')")
    
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    defer stmtIns.Close()

    http.Redirect(w, r, "/view/", http.StatusFound)
}

func deleteHandler(w http.ResponseWriter, r *http.Request){

    vars := mux.Vars(r)
    fmt.Println(vars["id"]);

    db, err := sql.Open("mysql", "root@tcp(localhost:3306)/staff") //localhost:3306
    if err != nil {
        panic(err.Error())  // Just for example purpose. You should use proper error handling instead of panic
    }
    defer db.Close()

        // Execute the query
    stmtIns, err := db.Query("DELETE FROM users WHERE id="+vars["id"])
    
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    defer stmtIns.Close()

    http.Redirect(w, r, "/view/", http.StatusFound)
}

func hapusHandler(w http.ResponseWriter, r *http.Request){

    vars := mux.Vars(r)
    fmt.Println(vars["id"]);

    db, err := sql.Open("mysql", "root@tcp(localhost:3306)/staff") //localhost:3306
    if err != nil {
        panic(err.Error())  // Just for example purpose. You should use proper error handling instead of panic
    }
    defer db.Close()

    // Execute the query
    stmtIns, err := db.Query("DELETE FROM man WHERE id="+vars["id"])

    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    defer stmtIns.Close()

    http.Redirect(w, r, "/home/", http.StatusFound)
}

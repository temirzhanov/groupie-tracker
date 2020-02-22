package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var myClient = &http.Client{Timeout: 10 * time.Second}

var parsed []ParsedArtist
var searchData SearchData

func init() {

	log.Println("Fetching data from api...")

	var artists []Artist
	var relations Relations

	// get artist array
	err := getJSON("https://groupietrackers.herokuapp.com/api/artists", &artists)
	if err != nil {
		log.Println("API artists data fetch failed")
		os.Exit(1)
	}

	// get relations array
	err = getJSON("https://groupietrackers.herokuapp.com/api/relation", &relations)
	if err != nil {
		log.Println("API relations data fetch failed")
		os.Exit(1)
	}

	// insert custom images
	for id, url := range customImages {
		artists[id-1].Image = url
	}

	// build CombinedData
	for index, artist := range artists {
		parsed = append(parsed, combineData(artist, relations.Index[index]))
	}

	// build SearchData
	buildSearchData()

}

func main() {

	// main handler
	http.HandleFunc("/", mainHandler)

	// static folder
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// init server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	port = ":" + port
	log.Print("Server started at http://localhost", port, "\n")
	err := http.ListenAndServe(port, nil)

	// open browser
	if len(os.Args) > 1 && os.Args[1] == "open" {
		openbrowser("http://localhost" + port)
	}

	// in case of error
	panic(err)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	url := strings.Trim(r.URL.String(), "/")

	switch r.Method {
	case "GET":
		if url == "" {
			handleHome(w, r)
		} else if artist, found := getArtist(url); found {
			handleArtist(w, r, artist)
		} else {
			handle404(w, r)
		}

	case "POST":

		if url == "search" {
			handleSearch(w, r)
		} else if url == "filter" {
			handleFilter(w, r)
		} else {
			handle400(w, r)
		}

	default:
		handle400(w, r)
	}

}

// route: /filter/
func handleFilter(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	fdata := FilterInput{}
	json.Unmarshal(body, &fdata)
	fmt.Printf("%+v\n", fdata)
}

// route: /
func handleHome(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/base.html", "templates/home.html")

	if err != nil {
		handle500(w, err)
		return
	}

	err = t.Execute(w, HomeData{parsed, searchData})
	if err != nil {
		handle500(w, err)
		return
	}
}

// route: /:artistID
func handleArtist(w http.ResponseWriter, r *http.Request, artist ParsedArtist) {
	t, err := template.ParseFiles("templates/base.html", "templates/artist.html")

	if err != nil {
		handle500(w, err)
		return
	}

	if err != nil {
		handle500(w, err)
		return
	}

	err = t.Execute(w, artist)

	if err != nil {
		handle500(w, err)
		return
	}
}

// route: /search
func handleSearch(w http.ResponseWriter, r *http.Request) {
	input := r.FormValue("input")

	input = trimWhitespaces(input)
	input = leaveOneSpace(input)

	foundData := seachKeyword(input)

	t, err := template.ParseFiles("templates/base.html", "templates/search.html")

	if err != nil {
		handle500(w, err)
		return
	}

	foundData.Keyword = input

	err = t.Execute(w, foundData)
	if err != nil {
		handle500(w, err)
		return
	}
}

// route: any
func handle404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)

	t, err := template.ParseFiles("templates/base.html", "templates/404.html")

	if err != nil {
		handle500(w, err)
		return
	}

	err = t.Execute(w, nil)

	if err != nil {
		handle500(w, err)
		return
	}
}

// bad request
func handle400(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(400)

	t, err := template.ParseFiles("templates/base.html", "templates/400.html")

	if err != nil {
		handle500(w, err)
		return
	}

	err = t.Execute(w, nil)

	if err != nil {
		handle500(w, err)
		return
	}
}

func handle500(w http.ResponseWriter, err error) {
	w.WriteHeader(500)

	t, other := template.ParseFiles("templates/base.html", "templates/500.html")

	if other != nil {
		w.Write([]byte("Something went wrong\nError 500\n" + other.Error()))
		return
	}

	fmt.Println(err)
	t.Execute(w, err)
}

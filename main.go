package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type ApiResponse struct {
	Url  string `json:"url"`
	Desc string `json:"title"`
}

func getJson(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}

func handler(w http.ResponseWriter, r *http.Request) {
	path := "https://deviantart.com/" + r.URL.Path[1:]
	json := ApiResponse{}
	err := getJson("https://backend.deviantart.com/oembed?url="+path, &json)
	if err != nil {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}
	templateFile, err := ioutil.ReadFile("template.html")
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	template := string(templateFile)
	fmt.Fprintf(w, template, json.Desc, json.Url, path)
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

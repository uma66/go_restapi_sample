package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"log"
	"net/http"
	"strings"
	"time"
)

var decoder = schema.NewDecoder()

// enum
// 参考: https://text.baldanders.info/golang/enumeration/
type VideoType int

const (
	Unknown VideoType = iota
	VideoA
	VideoB
	VideoC
)

var videoTypeMap = map[VideoType]string{
	Unknown: "Unknown",
	VideoA:  "VideoA",
	VideoB:  "VideoB",
	VideoC:  "VideoC",
}

func (v VideoType) String() string {
	if s, ok := videoTypeMap[v]; ok {
		return s
	}
	return "Unknown"
}

func GetVideoType(s string) VideoType {
	for key, value := range videoTypeMap {
		if strings.ToLower(value) == strings.ToLower(s) {
			return key
		}
	}
	return Unknown
}

type Person struct {
	Name string
	Age  int
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func handleRequests() {

	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "root")
	}).
		Methods(http.MethodGet).
		Schemes("http")
	// Schemes("https")

	router.HandleFunc("/video/{type}/{title}/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		fmt.Println("requestParams: ", vars)

		title, videoType, id := vars["title"], vars["type"], vars["id"]

		fmt.Println("log => ", [...]string{title, videoType, id})

		respondWithJson(w, http.StatusOK, vars)
	}).Methods(http.MethodGet)

	router.HandleFunc("/person", func(w http.ResponseWriter, r *http.Request) {

		q := r.URL.Query()
		name := q.Get("name")
		fmt.Println("query_params: ", name)

		// 以下をリクエストするとnameとageがPersonにマッピングされる
		// http://127.0.0.1:8001/person?name=aaa&age=30
		var person Person
		err := decoder.Decode(&person, r.URL.Query())
		if err != nil {
			// Handle error
			respondWithError(w, http.StatusBadRequest, "decode error")
			return
		}
		fmt.Println("struct_query_params: ", person)

		respondWithJson(w, http.StatusOK, name)
	}).Methods(http.MethodGet)

	router.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {

		err := r.ParseForm()
		if err != nil {
			// Handle error
		}

		param1 := r.FormValue("param1")
		param2 := r.FormValue("param2")

		fmt.Println("post_params1: ", param1)
		fmt.Println("post_params2: ", param2)

		// Postパラメータをstruct Personにマッピング
		var person Person
		err = decoder.Decode(&person, r.PostForm)
		if err != nil {
			// Handle error
		}

	}).Methods(http.MethodPost)

	// Gorillaはnet/httpのラッパーで、net/httpはリクエスト毎に新しいgoroutineを開始し、ハンドラに処理を受け渡す。
	// https://www.reddit.com/r/golang/comments/641z3b/is_gorilla_mux_router_or_the_http_package/
	// https://stackoverflow.com/questions/49975616/golang-rest-api-concurrency
	srv := &http.Server{
		Addr:    "0.0.0.0:8000",
		Handler: router, // Pass our instance of gorilla/mux in.
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	log.Fatal(srv.ListenAndServe())
}

func main() {
	fmt.Println("start api")

	// stringから列挙型に変換
	a := GetVideoType("VideoB")
	fmt.Println("enum: ", a)

	handleRequests()
}


package main

import (
    "fmt"
    "log"
    "time"
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

type Person struct {
    Name  string
    Age int
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
    
        title := vars["title"]
        videoType := vars["type"]
        id := vars["id"]
        
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

        // PostパラメータをPersonのstructにマッピング
        var person Person
        err = decoder.Decode(&person, r.PostForm)
        if err != nil {
            // Handle error
        }

	}).Methods(http.MethodPost)
    

    srv := &http.Server{
        Addr:         "0.0.0.0:8000",
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
    handleRequests()
}

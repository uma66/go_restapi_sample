
package main

import (
    "fmt"
    "log"
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "root")
}

func videoHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)

    fmt.Println("requestParams: ", vars)

    title := vars["title"]
    videoType := vars["type"]
    id := vars["id"]
    
    fmt.Println("log => ", [...]string{title, videoType, id})

    // if error == nil {
        // respondWithError(w, http.StatusBadRequest, "errorMsg")
    // }

    respondWithJson(w, http.StatusOK, vars)
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
    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/", rootHandler).Methods("GET")
    router.HandleFunc("/video/{type}/{title}/{id:[0-9]+}", videoHandler).Methods("GET")
    log.Fatal(http.ListenAndServe(":8000", router))
}

func main() {
    fmt.Println("start api")
    handleRequests()
}

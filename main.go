package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"

    "github.com/gorilla/mux"
)

const baseURL = "https://jsonplaceholder.typicode.com/posts"

// Post represents the structure of a post.
type Post struct {
    UserID int    `json:"userId"`
    ID     int    `json:"id,omitempty"`
    Title  string `json:"title"`
    Body   string `json:"body"`
}

// Handler for creating a post
func createPost(w http.ResponseWriter, r *http.Request) {
    var post Post
    err := json.NewDecoder(r.Body).Decode(&post)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    jsonData, _ := json.Marshal(post)
    resp, err := http.Post(baseURL, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body)
    w.WriteHeader(resp.StatusCode)
    w.Write(body)
}

// Handler for reading all posts
func readPosts(w http.ResponseWriter, r *http.Request) {
    resp, err := http.Get(baseURL)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body)
    w.WriteHeader(resp.StatusCode)
    w.Write(body)
}

// Handler for updating a post
func updatePost(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    // Read the request body
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    defer r.Body.Close() // Close the request body when done

    // Create the new request for updating the post
    req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/%s", baseURL, id), bytes.NewBuffer(body))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    // Read the response body
    responseBody, _ := ioutil.ReadAll(resp.Body)
    w.WriteHeader(resp.StatusCode)
    w.Write(responseBody)
}


// Handler for deleting a post
func deletePost(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%s", baseURL, id), nil)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    w.WriteHeader(resp.StatusCode)
}

// Main function to start the server
func main() {
    r := mux.NewRouter()

    r.HandleFunc("/posts", createPost).Methods("POST")
    r.HandleFunc("/posts", readPosts).Methods("GET")
    r.HandleFunc("/posts/{id}", updatePost).Methods("PUT")
    r.HandleFunc("/posts/{id}", deletePost).Methods("DELETE")

    fmt.Println("Server is running on port 8080")
    if err := http.ListenAndServe(":8080", r); err != nil {
        panic(err)
    }
}

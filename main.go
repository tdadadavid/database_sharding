package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/serialx/hashring"
	"log"
	"net/http"
)

type Client struct {
	Port string
	Conn *pgx.Conn
}

type Response struct {
	URL         string `json:"url"`
	URLId       string `json:"url_id"`
	ServerShard string `json:"server_shard"`
}

var clients map[string]Client

var servers = []string{
	"5433",
	"5434",
	"5435",
}

// add servers to ring
var serverRing = hashring.New(servers)

func main() {

	mux := http.NewServeMux()
	clients = connectClients()

	mux.HandleFunc("/create", createUrl)
	mux.HandleFunc("/get", getUrl)

	fmt.Println("server started on port 8080")
	http.ListenAndServe(":8080", mux)
}

func getUrl(w http.ResponseWriter, r *http.Request) {

}

func createUrl(w http.ResponseWriter, r *http.Request) {
	// get the url
	url := r.URL.Query().Get("url")

	// Create a new SHA-256 hash
	hash := sha256.New()

	// Write the URL data to the hash
	hash.Write([]byte(url))

	// Calculate the SHA-256 hash and encode it as base64
	hashBytes := hash.Sum(nil)
	result := base64.URLEncoding.EncodeToString(hashBytes)
	fmt.Println("hash", result)

	// pick 5 first letter
	urlId := result[0:5]

	// get the server for this hash
	port, ok := serverRing.GetNode(urlId)
	if !ok {
		return
	}
	server := clients[port]
	if _, err := server.Conn.Exec(context.Background(), "INSERT INTO url_table (url, url_id) VALUES ($1, $2)", url, urlId); err != nil {
		return
	}

	resp := Response{
		URL:         url,
		URLId:       urlId,
		ServerShard: port,
	}

	jsonData, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Error generating JSON", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(201)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func connectClients() map[string]Client {
	constr1 := "postgres://postgres:rootpassword@localhost:5433/postgres"
	constr2 := "postgres://postgres:rootpassword@localhost:5434/postgres"
	constr3 := "postgres://postgres:rootpassword@localhost:5435/postgres"

	conn5433, err := pgx.Connect(context.Background(), constr1)
	conn5434, _ := pgx.Connect(context.Background(), constr2)
	conn5435, _ := pgx.Connect(context.Background(), constr3)
	if err != nil {
		log.Fatal(err)
	}

	// Test the connection
	err = conn5433.Ping(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Test the connection
	err = conn5434.Ping(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Test the connection
	err = conn5435.Ping(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	return map[string]Client{
		"5433": {Port: "5433", Conn: conn5433},
		"5434": {Port: "5434", Conn: conn5434},
		"5435": {Port: "5435", Conn: conn5435},
	}
}

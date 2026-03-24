// Command server exposes a Repository adapter over HTTP+JSON.
// Routes mirror the messages in ports/item/repository.avpr.
//
// Usage:
//
//	server [--addr host:port] [--adapter memory]
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/example/avro-adapters/item"
)

func main() {
	addr := flag.String("addr", "localhost:8080", "listen address")
	adapter := flag.String("adapter", "memory", "adapter backend: memory")
	flag.Parse()

	var repo item.Repository
	switch *adapter {
	case "memory":
		repo = item.NewMemoryRepository()
	default:
		log.Fatalf("unknown adapter: %s", *adapter)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/item/get", handleGet(repo))
	mux.HandleFunc("/item/list", handleList(repo))
	mux.HandleFunc("/item/save", handleSave(repo))
	mux.HandleFunc("/item/delete", handleDelete(repo))

	fmt.Printf("ItemRepository server (%s adapter) → http://%s\n", *adapter, *addr)
	log.Fatal(http.ListenAndServe(*addr, mux))
}

func respond(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func handleGet(repo item.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			ID string `json:"id"`
		}
		json.NewDecoder(r.Body).Decode(&req)
		it, err := repo.Get(req.ID)
		if err != nil {
			respond(w, http.StatusNotFound, err)
			return
		}
		respond(w, http.StatusOK, it)
	}
}

func handleList(repo item.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, _ := repo.List()
		if items == nil {
			items = []*item.Item{}
		}
		respond(w, http.StatusOK, items)
	}
}

func handleSave(repo item.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var it item.Item
		json.NewDecoder(r.Body).Decode(&it)
		repo.Save(&it)
		respond(w, http.StatusOK, nil)
	}
}

func handleDelete(repo item.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			ID string `json:"id"`
		}
		json.NewDecoder(r.Body).Decode(&req)
		repo.Delete(req.ID)
		respond(w, http.StatusOK, nil)
	}
}

// Command demo illustrates hexagonal architecture dependency injection.
//
// run() is the application core: it depends only on item.Repository (the port),
// never on a concrete adapter. The adapter is selected and injected at the
// composition root (here: main). Swapping memory↔RPC requires no change to run().
//
// Usage:
//
//	demo [--adapter memory]
//	demo [--adapter rpc] [--addr host:port]
package main

import (
	"errors"
	"flag"
	"log"

	"github.com/example/avro-adapters/item"
	"github.com/pterm/pterm"
)

func main() {
	adapter := flag.String("adapter", "memory", "adapter: memory | rpc")
	addr := flag.String("addr", "localhost:8080", "RPC server address (adapter=rpc only)")
	flag.Parse()

	// Composition root: resolve the port to a concrete adapter.
	var repo item.Repository
	switch *adapter {
	case "memory":
		repo = item.NewMemoryRepository()
	case "rpc":
		repo = item.NewRpcClient(*addr)
	default:
		log.Fatalf("unknown adapter: %s", *adapter)
	}

	pterm.DefaultHeader.Printf("ItemRepository demo — adapter: %s", *adapter)
	run(repo)
}

// run is the application core. It only knows about the Repository port.
// No adapter types appear here — this is the hexagonal boundary.
func run(repo item.Repository) {
	// Save
	for _, it := range []*item.Item{
		{ID: "1", Name: "Widget"},
		{ID: "2", Name: "Gadget"},
		{ID: "3", Name: "Doohickey"},
	} {
		if err := repo.Save(it); err != nil {
			log.Fatalf("save: %v", err)
		}
	}

	// List
	items, _ := repo.List()
	pterm.DefaultSection.Println("All items")
	rows := pterm.TableData{{"ID", "Name"}}
	for _, it := range items {
		rows = append(rows, []string{it.ID, it.Name})
	}
	pterm.DefaultTable.WithHasHeader().WithBoxed().WithData(rows).Render()

	// Get
	pterm.DefaultSection.Println("Get item 2")
	it, err := repo.Get("2")
	if err != nil {
		log.Fatalf("get: %v", err)
	}
	pterm.Info.Printf("Found: id=%s name=%s\n", it.ID, it.Name)

	// Delete, then confirm via typed error
	repo.Delete("2")
	pterm.DefaultSection.Println("After deleting item 2")
	_, err = repo.Get("2")
	var nf *item.NotFoundError
	if errors.As(err, &nf) {
		pterm.Warning.Printf("Not found (expected): %v\n", err)
	}
}

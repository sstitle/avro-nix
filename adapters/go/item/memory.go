package item

import "sync"

// MemoryRepository is an in-process adapter backed by a Go map.
// Suitable for tests, local development, and in-process DI.
type MemoryRepository struct {
	mu    sync.RWMutex
	store map[string]*Item
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{store: make(map[string]*Item)}
}

func (r *MemoryRepository) Get(id string) (*Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	it, ok := r.store[id]
	if !ok {
		return nil, &NotFoundError{ID: id}
	}
	return it, nil
}

func (r *MemoryRepository) List() ([]*Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]*Item, 0, len(r.store))
	for _, it := range r.store {
		items = append(items, it)
	}
	return items, nil
}

func (r *MemoryRepository) Save(it *Item) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[it.ID] = it
	return nil
}

func (r *MemoryRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.store, id)
	return nil
}

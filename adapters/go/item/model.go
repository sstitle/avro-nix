package item

// Item is the domain entity.
// Derived from com.example.Item in ports/item/repository.avpr.
type Item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// NotFoundError is returned by adapters when no Item exists for the given ID.
// Derived from com.example.errors.NotFound in ports/item/repository.avpr.
type NotFoundError struct {
	ID string `json:"id"`
}

func (e *NotFoundError) Error() string { return "item not found: " + e.ID }

// Package item defines the Repository port for Item entities.
//
// The interface is derived from the Avro Protocol at ports/item/repository.avpr.
// Every adapter — in-memory, file-backed, RPC client — must satisfy this
// interface. The application core and service layer depend only on this package;
// they never import a concrete adapter directly.
package item

// Repository is the port. Any struct implementing it is a valid adapter.
// Swap adapters by changing what you inject at the composition root (main or
// the Nix module); no business logic changes.
type Repository interface {
	Get(id string) (*Item, error)
	List() ([]*Item, error)
	Save(item *Item) error
	Delete(id string) error
}

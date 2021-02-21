package store

type Store interface {
	GetEntry(id string) (string, error)
	InsertEntry(id string, contents string) error
}

package store

import "fmt"

type Store interface {
	GetEntry(id string) (string, error)
	InsertEntry(id string, contents string) error
}

// EntryNotFoundError occurs when no entry exists with the given ID.
type EntryNotFoundError struct {
	ID string
}

func (f EntryNotFoundError) Error() string {
	return fmt.Sprintf("could not find entry with ID=%v", f.ID)
}

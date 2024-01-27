package barreldb

type Table struct {
	db     BarrelDatabase
	prefix string
}

func NewTable(db BarrelDatabase, prefix string) *Table {
	return &Table{
		db:     db,
		prefix: prefix,
	}
}

// Has retrieves if a prefixed version of a key is present in the database.
func (t *Table) Has(key []byte) (bool, error) {
	return t.db.Has(append([]byte(t.prefix), key...))
}

// Get retrieves the given prefixed key if it's present in the database.
func (t *Table) Get(key []byte) ([]byte, error) {
	return t.db.Get(append([]byte(t.prefix), key...))
}

func (t *Table) Put(key []byte, value []byte) error {
	return t.db.Put(append([]byte(t.prefix), key...), value)
}

//// Delete removes the given prefixed key from the database.
//func (t *Table) Delete(key []byte) error {
//	return t.db.Delete(append([]byte(t.prefix), key...))
//}

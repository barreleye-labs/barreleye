package barreldb

type Table struct {
	DB     *BarrelDatabase
	Prefix string
}

func NewTable(db *BarrelDatabase, prefix string) *Table {
	return &Table{
		DB:     db,
		Prefix: prefix,
	}
}

// Has retrieves if a prefixed version of a key is present in the database.
func (t *Table) Has(key []byte) (bool, error) {
	return t.DB.Has(append([]byte(t.Prefix), key...))
}

// Get retrieves the given prefixed key if it's present in the database.
func (t *Table) Get(key []byte) ([]byte, error) {
	return t.DB.Get(append([]byte(t.Prefix), key...))
}

func (t *Table) Put(key []byte, value []byte) error {
	return t.DB.Put(append([]byte(t.Prefix), key...), value)
}

//// Delete removes the given prefixed key from the database.
//func (t *Table) Delete(key []byte) error {
//	return t.db.Delete(append([]byte(t.prefix), key...))
//}

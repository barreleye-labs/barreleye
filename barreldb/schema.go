package barreldb

var (
	// key
	databaseVersionKey = []byte("DatabaseVersion")
	lastHeaderKey      = []byte("LastHeader")
	lastBlockKey       = []byte("LastBlock")

	// prefix
	blockPrefix  = []byte("b")
	headerPrefix = []byte("h")
)

package barreldb

import (
	"flag"
	"github.com/syndtr/goleveldb/leveldb"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"runtime"
)

type BarrelDatabase struct {
	db     *leveldb.DB
	tables map[string]*Table
}

func New() (*BarrelDatabase, error) {
	db, err := leveldb.OpenFile(DefaultDataDir(), nil)
	if err != nil {
		return nil, err
	}
	return &BarrelDatabase{db: db, tables: make(map[string]*Table)}, nil
}

func (barrelDB *BarrelDatabase) Close() error {
	err := barrelDB.db.Close()
	return err
}

func (barrelDB *BarrelDatabase) CreateTable(name string, prefix string) error {
	table := NewTable(barrelDB, prefix)
	barrelDB.tables[name] = table
	return nil
}

func (barrelDB *BarrelDatabase) GetTable(name string) *Table {
	return barrelDB.tables[name]
}

func (barrelDB *BarrelDatabase) Get(key []byte) ([]byte, error) {

	return barrelDB.db.Get(key, nil)
}

func (barrelDB *BarrelDatabase) Has(key []byte) (bool, error) {
	return barrelDB.db.Has(key, nil)
}

func (barrelDB *BarrelDatabase) Put(key []byte, value []byte) error {
	return barrelDB.db.Put(key, value, nil)
}

func (barrelDB *BarrelDatabase) Delete(key []byte) error {
	return barrelDB.db.Delete(key, nil)
}

func DefaultDataDir() string {
	_, filename, _, _ := runtime.Caller(0)
	root := path.Join(path.Dir(filename), "../barreldb")

	if flag.Lookup("name") == nil {
		return filepath.Join(root, "data")
	}

	nodeName := flag.Lookup("name").Value.(flag.Getter).Get().(string)

	return filepath.Join(root, nodeName)

	//// Try to place the data folder in the user's home dir
	//home := homeDir()
	//if home != "" {
	//	switch runtime.GOOS {
	//	case "darwin":
	//		return filepath.Join(home, "Library", "Barreleye")
	//	case "windows":
	//		// We used to put everything in %HOME%\AppData\Roaming, but this caused
	//		// problems with non-typical setups. If this fallback location exists and
	//		// is non-empty, use it, otherwise DTRT and check %LOCALAPPDATA%.
	//		fallback := filepath.Join(home, "AppData", "Roaming", "Barreleye")
	//		appdata := windowsAppData()
	//		if appdata == "" || isNonEmptyDir(fallback) {
	//			return fallback
	//		}
	//		return filepath.Join(appdata, "Barreleye")
	//	default:
	//		return filepath.Join(home, ".barreleye")
	//	}
	//}
	//// As we cannot guess a stable location, return empty and handle later
	//return ""
}

func windowsAppData() string {
	v := os.Getenv("LOCALAPPDATA")
	if v == "" {
		// Windows XP and below don't have LocalAppData. Crash here because
		// we don't support Windows XP and undefining the variable will cause
		// other issues.
		panic("environment variable LocalAppData is undefined")
	}
	return v
}

func isNonEmptyDir(dir string) bool {
	f, err := os.Open(dir)
	if err != nil {
		return false
	}
	names, _ := f.Readdir(1)
	f.Close()
	return len(names) > 0
}

func homeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return ""
}

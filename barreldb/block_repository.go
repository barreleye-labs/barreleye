package barreldb

func (barrelDB *BarrelDatabase) CreateBlock(key string, value string) error {
	if err := barrelDB.GetTable("block").Put([]byte(key), []byte(value)); err != nil {
		return err
	}
	return nil
}

func (barrelDB *BarrelDatabase) GetBlock(key string) ([]byte, error) {
	data, err := barrelDB.GetTable("block").Get([]byte(key))
	if err != nil {
		return nil, err
	}
	return data, nil
}

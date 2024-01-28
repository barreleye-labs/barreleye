package core

func (bc *Blockchain) GetBlockFromDB(key string) ([]byte, error) {
	data, err := bc.db.GetTable("block").Get([]byte(key))
	if err != nil {
		return nil, err
	}
	return data, nil
}

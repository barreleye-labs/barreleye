package core

func (bc *Blockchain) CreateBlock(key string, value string) error {
	if err := bc.db.GetTable("block").Put([]byte(key), []byte(value)); err != nil {
		return err
	}
	return nil
}

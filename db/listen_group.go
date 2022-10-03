package db

func (gb *GormDb) InsertListenGroupOne(lg ListenGroup) (int64, error) {
	res := gb.DB.Create(&lg)
	return res.RowsAffected, res.Error
}

func (gb *GormDb) FindListenGroupMany(lg ListenGroup) ([]ListenGroup, int64, error) {
	var lgs []ListenGroup
	result := gb.DB.Where(&lg).Find(&lgs)
	return lgs, result.RowsAffected, result.Error
}

func (gb *GormDb) DeleteListenGroupOne(lg ListenGroup) (int64, error) {
	res := gb.DB.Delete(&lg)
	return res.RowsAffected, res.Error
}

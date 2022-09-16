package db

import "github.com/solywsh/qqBot-revue/service/wzxy"

func (gb *GormDb) InsertWzxyTokenOne(t wzxy.TokenWzxy) (int64, error) {
	res := gb.DB.Create(&t)
	return res.RowsAffected, res.Error
}

func (gb *GormDb) UpdateWzxyTokenOne(t wzxy.TokenWzxy, saveZero bool) (int64, error) {
	if saveZero {
		res := gb.DB.Save(t)
		return res.RowsAffected, res.Error
	} else {
		res := gb.DB.Model(&t).Updates(t)
		return res.RowsAffected, res.Error
	}
}

func (gb *GormDb) DeleteWzxyTokenOne(t wzxy.TokenWzxy) (int64, error) {
	res := gb.DB.Delete(&t)
	return res.RowsAffected, res.Error
}

func (gb *GormDb) FindWzxyTokenMany(w wzxy.TokenWzxy) ([]wzxy.TokenWzxy, int64, error) {
	var ts []wzxy.TokenWzxy
	result := gb.DB.Where(&w).Find(&ts)
	return ts, result.RowsAffected, result.Error
}

func (gb *GormDb) InsertWzxyUserOne(u wzxy.UserWzxy) (int64, error) {
	res := gb.DB.Create(&u)
	return res.RowsAffected, res.Error
}

func (gb *GormDb) UpdateWzxyUserOne(u wzxy.UserWzxy, saveZero bool) (int64, error) {
	if saveZero {
		res := gb.DB.Save(u)
		return res.RowsAffected, res.Error
	} else {
		res := gb.DB.Model(&u).Updates(u)
		return res.RowsAffected, res.Error
	}
}

func (gb *GormDb) DeleteWzxyUserOne(u wzxy.UserWzxy) (int64, error) {
	res := gb.DB.Delete(&u)
	return res.RowsAffected, res.Error
}

func (gb *GormDb) FindWzxyUserMany(u wzxy.UserWzxy) ([]wzxy.UserWzxy, int64, error) {
	var us []wzxy.UserWzxy
	result := gb.DB.Where(&u).Find(&us)
	return us, result.RowsAffected, result.Error
}

func (gb *GormDb) GetAllWzxyUser() ([]wzxy.UserWzxy, int64, error) {
	var us []wzxy.UserWzxy
	result := gb.DB.Find(&us)
	return us, result.RowsAffected, result.Error
}

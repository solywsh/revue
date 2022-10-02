package db

import (
	"database/sql"
	"github.com/solywsh/qqBot-revue/service/wzxy"
)

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

func (gb *GormDb) InsertMonitorWzxyOne(m wzxy.MonitorWzxy) (int64, error) {
	res := gb.DB.Create(&m)
	return res.RowsAffected, res.Error
}

func (gb *GormDb) FindMonitorWzxyMany(m wzxy.MonitorWzxy) ([]wzxy.MonitorWzxy, int64, error) {
	var ms []wzxy.MonitorWzxy
	result := gb.DB.Where(&m).Find(&ms)
	return ms, result.RowsAffected, result.Error
}

func (gb *GormDb) DeleteMonitorWzxyOne(m wzxy.MonitorWzxy) (int64, error) {
	res := gb.DB.Delete(&m)
	return res.RowsAffected, res.Error
}

func (gb *GormDb) UpdateMonitorWzxyOne(m wzxy.MonitorWzxy, saveZero bool) (int64, error) {
	if saveZero {
		res := gb.DB.Save(m)
		return res.RowsAffected, res.Error
	} else {
		res := gb.DB.Model(&m).Updates(m)
		return res.RowsAffected, res.Error
	}
}

func (gb *GormDb) InsertClassStudentWzxyOne(csw wzxy.ClassStudentWzxy) (int64, error) {
	res := gb.DB.Create(&csw)
	return res.RowsAffected, res.Error
}

func (gb *GormDb) InsertClassStudentWzxyMany(csws []wzxy.ClassStudentWzxy) (int64, error) {
	res := gb.DB.Create(&csws)
	return res.RowsAffected, res.Error
}

func (gb *GormDb) FindClassStudentWzxyMany(csw wzxy.ClassStudentWzxy) ([]wzxy.ClassStudentWzxy, int64, error) {
	var csws []wzxy.ClassStudentWzxy
	result := gb.DB.Where(&csw).Find(&csws)
	return csws, result.RowsAffected, result.Error
}

//// Delete ClassStudentWzxy Many
//func (gb *GormDb) DeleteClassStudentWzxyMany(csw wzxy.ClassStudentWzxy) (int64, error) {
//	res := gb.DB.Where(&csw).Delete(&wzxy.ClassStudentWzxy{})
//	return res.RowsAffected, res.Error
//}

func (gb *GormDb) DeleteClassStudentWzxyOne(csw wzxy.ClassStudentWzxy) (int64, error) {
	res := gb.DB.Delete(&csw)
	return res.RowsAffected, res.Error
}

func (gb *GormDb) UpdateClassStudentWzxyOne(csw wzxy.ClassStudentWzxy, saveZero bool) (int64, error) {
	if saveZero {
		res := gb.DB.Save(csw)
		return res.RowsAffected, res.Error
	} else {
		res := gb.DB.Model(&csw).Updates(csw)
		return res.RowsAffected, res.Error
	}
}

func (gb *GormDb) FindClassStudentWzxyByKeywords(keywords string) ([]wzxy.ClassStudentWzxy, int64, error) {
	var csws []wzxy.ClassStudentWzxy
	results := gb.DB.Raw("SELECT * FROM class_student_wzxies WHERE "+
		"name = @name OR "+
		"student_id = @student_id OR "+
		"class_name = @class_name OR "+
		"user_id = @user_id",
		sql.Named("name", keywords),
		sql.Named("student_id", keywords),
		sql.Named("class_name", keywords),
		sql.Named("user_id", keywords)).
		Find(&csws)
	return csws, results.RowsAffected, results.Error
}

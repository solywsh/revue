package db

import "time"

func (gb *GormDb) FindOrCreateUserSession(userId, appName string) (UserSession, error) {
	var us UserSession
	res := gb.DB.Where(UserSession{UserId: userId, AppName: appName}).Attrs(UserSession{
		UserId:     userId,
		AppName:    appName,
		UpdateTime: time.Now(),
		Status:     0,
	}).FirstOrCreate(&us)
	if res.RowsAffected == 1 {
		return us, res.Error
	} else {
		return us, res.Error
	}
}

func (gb *GormDb) FindUserSessionMany(us UserSession) ([]UserSession, int64, error) {
	var many []UserSession
	result := gb.DB.Where(us).Find(&many)
	return many, result.RowsAffected, result.Error
}

func (gb *GormDb) DeleteUserSession(userId, appName string) (bool, error) {
	var us UserSession
	res := gb.DB.Where(UserSession{UserId: userId, AppName: appName}).First(&us)
	if res.RowsAffected >= 1 {
		gb.DB.Delete(&us)
		return true, res.Error
	} else {
		return false, res.Error
	}
}

func (gb *GormDb) UpdateUserSessionMany(us UserSession, saveZero bool) (bool, error) {
	if saveZero {
		res := gb.DB.Save(us)
		if res.RowsAffected >= 1 {
			return true, res.Error
		} else {
			return false, res.Error
		}
	} else {
		res := gb.DB.Model(&us).Updates(us)
		if res.RowsAffected >= 1 {
			return true, res.Error
		} else {
			return false, res.Error
		}
	}
}

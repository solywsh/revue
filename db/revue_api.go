package db

import uuid "github.com/satori/go.uuid"

// InsertRevueApiToken 查询并创建,查询到对应的信息,则不创建返回false,没有则创建,返回true
func (gb *GormDb) InsertRevueApiToken(userId string, permission uint) (bool, string) {
	var rt RevueApiToken
	res := gb.DB.Where(RevueApiToken{UserId: userId}).Attrs(RevueApiToken{
		UserId:     userId,
		Token:      uuid.NewV4().String(),
		Permission: permission,
	}).FirstOrCreate(&rt)
	if res.RowsAffected == 1 {
		return true, rt.Token
	} else {
		return false, rt.Token
	}
}

// GetRevueApiToken 根据qq号得到token,注意此函数必须在qq号存在的业务逻辑下使用
func (gb *GormDb) GetRevueApiToken(userId string) string {
	var rt RevueApiToken
	gb.DB.Where(RevueApiToken{UserId: userId}).First(&rt)
	return rt.Token
}

// SearchRevueApiToken 根据传入的qq号和token判断是否具有权限,并返回权限
func (gb *GormDb) SearchRevueApiToken(userId, token string) (bool, uint) {
	var rt RevueApiToken
	res := gb.DB.Where(RevueApiToken{UserId: userId, Token: token}).First(&rt)
	if res.RowsAffected >= 1 {
		return true, rt.Permission
	} else {
		return false, rt.Permission
	}
}

// FindRevueApiTokenMany 查询多个
func (gb *GormDb) FindRevueApiTokenMany(rt RevueApiToken) ([]RevueApiToken, int64, error) {
	var many []RevueApiToken
	result := gb.DB.Where(rt).Find(&many)
	return many, result.RowsAffected, result.Error
}

// ResetRevueApiToken 重新设置token
func (gb *GormDb) ResetRevueApiToken(userId string) (bool, string) {
	var rt RevueApiToken
	res := gb.DB.Where(RevueApiToken{UserId: userId}).First(&rt)
	if res.RowsAffected >= 1 {
		gb.DB.Model(&rt).Update("token", uuid.NewV4().String())
		return true, rt.Token
	} else {
		// 不存在
		return false, rt.Token
	}
}

// DeleteRevueApiToken 删除token
func (gb *GormDb) DeleteRevueApiToken(userId string) (bool, string) {
	var rt RevueApiToken
	res := gb.DB.Where(RevueApiToken{UserId: userId}).First(&rt)
	if res.RowsAffected >= 1 {
		gb.DB.Unscoped().Delete(&rt)
		return true, rt.Token
	} else {
		// 不存在
		return false, rt.Token
	}
}

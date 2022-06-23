package db

// GetAllKeywordsReply 获取所有的关键词回复
func (gb *GormDb) GetAllKeywordsReply(res *[]KeywordsReply) {
	gb.DB.Find(&res)
}

// UpdateKeywordsReply 插入/更新关键词回复,如果存在对应关键词,则更新,不存在则insert
func (gb *GormDb) UpdateKeywordsReply(info KeywordsReply) {
	var kr KeywordsReply
	if info.Flag == 1 {
		// 第一次创建
		gb.DB.Create(&info)
	} else if info.Flag == 2 {
		// 第二次添加关键词
		gb.DB.Where(KeywordsReply{ID: info.ID}).First(&kr)
		// 如果第二次为更新之前查询,需要更新设置人的qq
		gb.DB.Model(&kr).Updates(KeywordsReply{Keywords: info.Keywords, Userid: info.Userid, Flag: 2})
	} else if info.Flag == 3 {
		// 第三次添加回复和匹配模式
		gb.DB.Where(KeywordsReply{ID: info.ID}).First(&kr)
		gb.DB.Model(&kr).Updates(KeywordsReply{Msg: info.Msg, Mode: info.Mode, Flag: 3})
	}
	//fmt.Println(KR.ID, KR.Keywords, KR.Msg)
}

// GetKeywordsReplyFlag 根据userid查找该userid是否存在正在记录中的自动回复
func (gb *GormDb) GetKeywordsReplyFlag(userId string) (bool, KeywordsReply) {
	var kr KeywordsReply
	if res := gb.DB.Where("userid = ? AND flag <> ?", userId, 3).First(&kr); res.RowsAffected >= 1 {
		return true, kr
	} else {
		return false, kr
	}
}

// SearchKeywordsReply 根据关键词搜索是否存在记录,如果存在记录则返回对应回复
func (gb *GormDb) SearchKeywordsReply(keywords string) (bool, KeywordsReply) {
	var kr KeywordsReply
	// 存在关键词并且状态flag为0
	if res := gb.DB.Where(KeywordsReply{Keywords: keywords, Flag: 0}).First(&kr); res.RowsAffected >= 1 {
		return true, kr
	}
	return false, kr
}

// DeleteKeywordsReply 删除对应的关键词回复
func (gb *GormDb) DeleteKeywordsReply(id uint) {
	gb.DB.Delete(&KeywordsReply{}, id)
}

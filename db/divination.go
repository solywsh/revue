package db

import (
	"github.com/solywsh/qqBot-revue/service/divination"
	"time"
)

// GetDivination 查询或创建,如果查询到结果返回false + tag,没有查询到则创建之后返回true + tag
func (gb *GormDb) GetDivination(userid string) (bool, string) {
	var divin Divination
	date := time.Now().Format("20060102")
	res := gb.DB.Where(Divination{
		Date:   date,
		UserId: userid,
	}).Attrs(Divination{
		Date:   date,
		UserId: userid,
		Tag:    divination.NewDivination(), // 创建一个新的
	}).FirstOrCreate(&divin)
	if res.RowsAffected == 1 {
		return true, divin.Tag
	} else {
		return false, divin.Tag
	}
}

package db

import (
	"github.com/solywsh/qqBot-revue/service/pa"
	"time"
)

// GetProgrammerAlmanac 查询今天的黄历,如果没有,则存入数据库返回内容,如果有则直接返回查询结果
func (gb *GormDb) GetProgrammerAlmanac() string {
	date := time.Now().Format("20060102")
	var proAlmanac ProgrammerAlmanac
	gb.DB.Where(ProgrammerAlmanac{Date: date}).Attrs(ProgrammerAlmanac{
		Date:    date,
		Almanac: pa.NewCalendar(),
	}).FirstOrCreate(&proAlmanac)
	return proAlmanac.Almanac
}

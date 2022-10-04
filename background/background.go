package background

import (
	"github.com/solywsh/qqBot-revue/cq"
	"github.com/solywsh/qqBot-revue/db"
)

/*
	后台任务执行
*/
var gdb = db.NewDB() // 初始化操作数据库

func Services() {
	go func() {
		// service list
		go wzxyService()           // 我在校园后台服务
		go cq.ListenGroupService() // 监听群组同步服务
	}()
}

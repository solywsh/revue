package cq

import (
	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/solywsh/qqBot-revue/db"
	"log"
	"time"
)

var (
	ListenGroup = cmap.New[struct{}]()
)

func ListenGroupService() {
	for {
		manyListenGroup, _, err := gdb.FindListenGroupMany(db.ListenGroup{})
		if err != nil {
			log.Println("ListenGroupService FindListenGroupMany err:", err)
			return
		}
		ListenGroup.Clear()
		for _, group := range manyListenGroup {
			ListenGroup.Set(group.Group, struct{}{})
		}
		time.Sleep(time.Minute * 5)
	}
}

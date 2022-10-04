package cq

import (
	cmap "github.com/orcaman/concurrent-map"
	"github.com/solywsh/qqBot-revue/db"
	"log"
	"time"
)

var (
	ListenGroup = cmap.New()
)

func ListenGroupService() {
	for {
		manyListenGroup, _, err := gdb.FindListenGroupMany(db.ListenGroup{})
		if err != nil {
			log.Println("ListenGroupService FindListenGroupMany err:", err)
			return
		}
		for _, group := range manyListenGroup {
			ListenGroup.Clear()
			if !ListenGroup.Has(group.Group) {
				ListenGroup.Set(group.Group, struct{}{})
			}
		}
		time.Sleep(time.Minute * 5)
	}
}

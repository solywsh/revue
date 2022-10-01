package wzxy

import (
	"fmt"
	"testing"
)

func TestUserWzxy_GetUncheckList(t *testing.T) {
	ut := UserWzxy{
		Jwsession: "75b56a2c6592450b8bedc318de6dcdd0",
		UserAgent: "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.29(0x18001d2b) NetType/WIFI Language/zh_CN miniProgram/wxce6d08f781975d91",
	}
	list, err := ut.GetUncheckList(2)
	for i, wzxy := range list {
		fmt.Println(i, wzxy.ClassName, wzxy.StudentId)
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(list)
}

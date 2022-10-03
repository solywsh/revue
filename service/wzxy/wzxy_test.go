package wzxy

import (
	"fmt"
	"testing"
)

func TestUserWzxy_GetUncheckList(t *testing.T) {
	ut := UserWzxy{
		Jwsession: "ed8465b4a4f648e5b4b53c3eaa26f8ca",
		UserAgent: "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.29(0x18001d2b) NetType/WIFI Language/zh_CN miniProgram/wxce6d08f781975d91",
	}
	list, code := ut.GetUnSignedList()
	fmt.Println(code)
	fmt.Println(list)
}

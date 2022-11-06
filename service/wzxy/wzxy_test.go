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

func TestClassStudentCheck(t *testing.T) {

	uw := UserWzxy{
		Jwsession: "a7a63b0d3ccd47b28c262ba6cd8915b0",
	}

	uc := ClassStudentWzxy{
		StudentId: "",
		Name:      "",
		ID:        2,
		ClassName: "",
		UserId:    "",
		checkId:   "509359067697779115",
	}
	_, message := uw.ClassCheckOperate(2, uc)
	fmt.Println("message:" + message)

}

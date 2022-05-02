package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/thedevsaddam/gojsonq"
	"strconv"
)

//
//  formatAccessUrl
//  @Description: 格式化url，根据配置文件是否开启鉴权格式化
//  @param str
//  @return string
//
func formatAccessUrl(str string) string {
	if yamlConfig.fAuth.enable {
		return yamlConfig.urlHeader + str + "?access_token=" + yamlConfig.fAuth.tokenOrSecret
	} else {
		return yamlConfig.urlHeader + str
	}
}

// 发送私聊消息
func sendPrivateMsg(userId, groupId, message, autoEscape string) (string, error) {
	client := resty.New()
	post, err := client.R().SetQueryParams(map[string]string{
		"user_id":     userId,
		"group_id":    groupId,
		"message":     message,
		"auto_escape": autoEscape,
	}).Post(formatAccessUrl("/send_private_msg"))
	if err != nil {
		return "", err
	}
	rJson := gojsonq.New().JSONString(string(post.Body()))
	if rJson.Reset().Find("retcode") != nil && rJson.Reset().Find("retcode").(float64) != 0.0 {
		return "", fmt.Errorf(string(post.Body()))
	}
	messageId := strconv.Itoa(int(rJson.Reset().Find("data.message_id").(float64)))
	return messageId, err
}

// 发送群消息
func sendGroupMsg(groupId, message, autoEscape string) (string, error) {
	client := resty.New()
	post, err := client.R().SetQueryParams(map[string]string{
		"group_id":    groupId,
		"message":     message,
		"auto_escape": autoEscape,
	}).Post(formatAccessUrl("/send_group_msg"))
	if err != nil {
		return "", err
	}
	rJson := gojsonq.New().JSONString(string(post.Body()))
	if rJson.Reset().Find("retcode") != nil && rJson.Reset().Find("retcode").(float64) != 0.0 {
		return "", fmt.Errorf(string(post.Body()))
	}
	messageId := strconv.Itoa(int(rJson.Reset().Find("data.message_id").(float64)))
	return messageId, err
}

//
//  sendMsg
//  @Description: 发送消息
//  @param messageType 消息类型private、group,如果不传入根据传入id判断
//  @param userId private时对方的qq
//  @param groupId group时群号
//  @param message 消息
//  @param autoEscape 是否解析CQ码
//  @return string 返回message_id
//  @return error
//
func sendMsg(messageType, userId, groupId, message, autoEscape string) (string, error) {
	client := resty.New()
	post, err := client.R().SetQueryParams(map[string]string{
		"message_type": messageType,
		"user_id":      userId,
		"group_id":     groupId,
		"message":      message,
		"auto_escape":  autoEscape,
	}).Post(formatAccessUrl("/send_msg"))
	if err != nil {
		return "", err
	}
	rJson := gojsonq.New().JSONString(string(post.Body()))
	if rJson.Reset().Find("retcode") != nil && rJson.Reset().Find("retcode").(float64) != 0.0 {
		return "", fmt.Errorf(string(post.Body()))
	}
	messageId := strconv.Itoa(int(rJson.Reset().Find("data.message_id").(float64)))
	return messageId, err
}

//
//  deleteMsg
//  @Description: 撤回消息
//  @param messageId 需要撤回的消息Id
//
func deleteMsg(messageId string) {
	client := resty.New()
	_, _ = client.R().SetQueryParams(map[string]string{
		"message_id": messageId,
	}).Post(formatAccessUrl("/delete_msg"))
}

//
//  deleteFriend
//  @Description: 删除好友
//  @param friendId 好友qq号
//
func deleteFriend(friendId string) {
	client := resty.New()
	_, _ = client.R().SetQueryParams(map[string]string{
		"friend_id": friendId,
	}).Post(formatAccessUrl("/delete_friend"))
}

//
//  getMsg
//  @Description: 获取消息(信息)
//  @param messageId
//  @return string json的string格式
//	{
//		"data": {
//		"group": false,
//			"message": "[CQ:face,id=123]",
//			"message_id": -123150518,
//			"message_id_v2": "00000000b629494a00005983",
//			"message_seq": 22915,
//			"message_type": "private",
//			"real_id": 22915,
//			"sender": {
//			"nickname": "boogiepop",
//				"user_id": 121212121212
//		},
//		"time": 1651404575
//	},
//		"retcode": 0,
//		"status": "ok"
//	}
//  @return error
//
func getMsg(messageId string) (string, error) {
	client := resty.New()
	post, err := client.R().SetQueryParams(map[string]string{
		"message_id": messageId,
	}).Post(formatAccessUrl("/get_msg"))
	if err != nil {
		return "", err
	}
	rJson := gojsonq.New().JSONString(string(post.Body()))
	if rJson.Reset().Find("retcode") != nil && rJson.Reset().Find("retcode").(float64) != 0.0 {
		return "", fmt.Errorf(string(post.Body()))
	}
	return string(post.Body()), nil

}

//
//  getFriendList
//  @Description: 获取好友列表
//  @return string:json的string格式
//	{
//	   "data": [
//		   {
//			   "nickname": "QRZ?",
//			   "remark": "QRZ?",
//			   "user_id": 1212121212
//		   },
//		   {
//			   "nickname": "ACE OF SPADES",
//			   "remark": "XX",
//			   "user_id": 1212121212
//		   },
//		   {
//			   "nickname": "boogiepop",
//			   "remark": "boogiepop",
//			   "user_id": 12112121212
//		   }
//	   ],
//	   "retcode": 0,
//	   "status": "ok"
//	}
//  @return error
//
func getFriendList() (string, error) {
	client := resty.New()
	post, err := client.R().Post(formatAccessUrl("/get_friend_list"))
	if err != nil {
		return "", err
	}
	rJson := gojsonq.New().JSONString(string(post.Body()))
	if rJson.Reset().Find("retcode") != nil && rJson.Reset().Find("retcode").(float64) != 0.0 {
		return "", fmt.Errorf(string(post.Body()))
	}
	return string(post.Body()), nil

}

//
//  getGroupInfo
//  @Description:
//  @param groupId 群号
//  @param noCache 是否使用缓存
//  @return string
//	{
//    "data": {
//        "group_create_time": 0,
//        "group_id": 121212212,
//        "group_level": 0,
//        "group_memo": "",
//        "group_name": "XXXXXXX",
//        "max_member_count": 200,
//        "member_count": 7
//    },
//    "retcode": 0,
//    "status": "ok"
//	}
//  @return error
//
func getGroupInfo(groupId, noCache string) (string, error) {
	client := resty.New()
	post, err := client.R().SetQueryParams(map[string]string{
		"group_id": groupId,
		"no_cache": noCache,
	}).Post(formatAccessUrl("/get_group_info"))
	if err != nil {
		return "", err
	}
	rJson := gojsonq.New().JSONString(string(post.Body()))
	if rJson.Reset().Find("retcode") != nil && rJson.Reset().Find("retcode").(float64) != 0.0 {
		return "", fmt.Errorf(string(post.Body()))
	}
	return string(post.Body()), nil
}

//
//  getGroupList
//  @Description: 获取群列表
//  @return string
//  {
//	"data": [
//		{
//		"group_create_time": 0,
//		"group_id": 1212121212,
//		"group_level": 0,
//		"group_memo": "",
//		"group_name": "bot测试",
//		"max_member_count": 200,
//		"member_count": 3
//		},
//		{
//		"group_create_time": 0,
//		"group_id": 1212121212,
//		"group_level": 0,
//		"group_memo": "",
//		"group_name": "test",
//		"max_member_count": 200,
//		"member_count": 10
//		}
//	],
//	"retcode": 0,
//	"status": "ok"
//	}
//  @return error
//
func getGroupList() (string, error) {
	client := resty.New()
	post, err := client.R().Post(formatAccessUrl("/get_group_list"))
	if err != nil {
		return "", err
	}
	rJson := gojsonq.New().JSONString(string(post.Body()))
	fmt.Println(string(post.Body()))
	if rJson.Reset().Find("retcode") != nil && rJson.Reset().Find("retcode").(float64) != 0.0 {
		return "", fmt.Errorf(string(post.Body()))
	}
	return string(post.Body()), nil
}

//
//  getGroupMemberList
//  @Description: 获取群成员信息
//  @param groupId 群号
//  @return string
//	{
//    "data": [
//        {
//            "age": 0,
//            "area": "",
//            "card": "XXXX",
//            "card_changeable": false,
//            "group_id": 1212121212,
//            "join_time": 1635232398,
//            "last_sent_time": 1651401105,
//            "level": "1",
//            "nickname": "XXXX",
//            "role": "owner",
//            "sex": "male",
//            "shut_up_timestamp": 0,
//            "title": "",
//            "title_expire_time": 0,
//            "unfriendly": false,
//            "user_id": 121212121212
//        },
//        {
//            "age": 0,
//            "area": "",
//            "card": "",
//            "card_changeable": false,
//            "group_id": 1212121212,
//            "join_time": 1635232488,
//            "last_sent_time": 1651392760,
//            "level": "1",
//            "nickname": "XXXX",
//            "role": "admin",
//            "sex": "male",
//            "shut_up_timestamp": 0,
//            "title": "",
//            "title_expire_time": 0,
//            "unfriendly": false,
//            "user_id": 12112121212
//        }
//    ],
//    "retcode": 0,
//    "status": "ok"
//}
//  @return error
//
func getGroupMemberList(groupId string) (string, error) {
	client := resty.New()
	post, err := client.R().SetQueryParams(map[string]string{
		"group_id": groupId,
	}).Post(formatAccessUrl("/get_group_member_list"))
	if err != nil {
		return "", err
	}
	rJson := gojsonq.New().JSONString(string(post.Body()))
	if rJson.Reset().Find("retcode") != nil && rJson.Reset().Find("retcode").(float64) != 0.0 {
		return "", fmt.Errorf(string(post.Body()))
	}
	return string(post.Body()), nil
}

//
//  getGroupAtAllRemain
//  @Description: 获取群@全体成员 剩余次数
//  @param groupId 群号
//  @return string
//	{
//    "data": {
//        "can_at_all": true,//是否可以@全体成员
//        "remain_at_all_count_for_group": 19,//群内所有管理当天剩余 @全体成员 次数
//        "remain_at_all_count_for_uin": 9 //Bot 当天剩余 @全体成员 次数
//    },
//    "retcode": 0,
//    "status": "ok"
//}
//  @return error
//
func getGroupAtAllRemain(groupId string) (string, error) {
	client := resty.New()
	post, err := client.R().SetQueryParams(map[string]string{
		"group_id": groupId,
	}).Post(formatAccessUrl("/get_group_at_all_remain"))
	if err != nil {
		return "", err
	}
	rJson := gojsonq.New().JSONString(string(post.Body()))
	if rJson.Reset().Find("retcode") != nil && rJson.Reset().Find("retcode").(float64) != 0.0 {
		return "", fmt.Errorf(string(post.Body()))
	}
	return string(post.Body()), nil
}

package cq

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/thedevsaddam/gojsonq"
)

// GetMsg
//
//	 @Description: 获取消息(信息)
//	 @param messageId
//	 @return string json的string格式
//		{
//			"data": {
//			"group": false,
//				"message": "[CQ:face,id=123]",
//				"message_id": -123150518,
//				"message_id_v2": "00000000b629494a00005983",
//				"message_seq": 22915,
//				"message_type": "private",
//				"real_id": 22915,
//				"sender": {
//				"nickname": "boogiepop",
//					"user_id": 121212121212
//			},
//			"time": 1651404575
//		},
//			"retcode": 0,
//			"status": "ok"
//		}
//	 @return error
func GetMsg(messageId string) (string, error) {
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

// GetFriendList
//
//	 @Description: 获取好友列表
//	 @return string:json的string格式
//		{
//		   "data": [
//			   {
//				   "nickname": "QRZ?",
//				   "remark": "QRZ?",
//				   "user_id": 1212121212
//			   },
//			   {
//				   "nickname": "ACE OF SPADES",
//				   "remark": "XX",
//				   "user_id": 1212121212
//			   },
//			   {
//				   "nickname": "boogiepop",
//				   "remark": "boogiepop",
//				   "user_id": 12112121212
//			   }
//		   ],
//		   "retcode": 0,
//		   "status": "ok"
//		}
//	 @return error
func GetFriendList() (string, error) {
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

// GetGroupInfo
//
//	 getGroupInfo
//	 @Description:
//	 @param groupId 群号
//	 @param noCache 是否使用缓存
//	 @return string
//		{
//	   "data": {
//	       "group_create_time": 0,
//	       "group_id": 121212212,
//	       "group_level": 0,
//	       "group_memo": "",
//	       "group_name": "XXXXXXX",
//	       "max_member_count": 200,
//	       "member_count": 7
//	   },
//	   "retcode": 0,
//	   "status": "ok"
//		}
//	 @return error
func GetGroupInfo(groupId, noCache string) (string, error) {
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

// GetGroupList
//
//	 @Description: 获取群列表
//	 @return string
//	 {
//		"data": [
//			{
//			"group_create_time": 0,
//			"group_id": 1212121212,
//			"group_level": 0,
//			"group_memo": "",
//			"group_name": "bot测试",
//			"max_member_count": 200,
//			"member_count": 3
//			},
//			{
//			"group_create_time": 0,
//			"group_id": 1212121212,
//			"group_level": 0,
//			"group_memo": "",
//			"group_name": "test",
//			"max_member_count": 200,
//			"member_count": 10
//			}
//		],
//		"retcode": 0,
//		"status": "ok"
//		}
//	 @return error
func GetGroupList() (string, error) {
	client := resty.New()
	post, err := client.R().Post(formatAccessUrl("/get_group_list"))
	if err != nil {
		return "", err
	}
	rJson := gojsonq.New().JSONString(string(post.Body()))
	//fmt.Println(string(post.Body()))
	if rJson.Reset().Find("retcode") != nil && rJson.Reset().Find("retcode").(float64) != 0.0 {
		return "", fmt.Errorf(string(post.Body()))
	}
	return string(post.Body()), nil
}

// GetGroupMemberList
//
//	 @Description: 获取群成员信息
//	 @param groupId 群号
//	 @return string
//		{
//	   "data": [
//	       {
//	           "age": 0,
//	           "area": "",
//	           "card": "XXXX",
//	           "card_changeable": false,
//	           "group_id": 1212121212,
//	           "join_time": 1635232398,
//	           "last_sent_time": 1651401105,
//	           "level": "1",
//	           "nickname": "XXXX",
//	           "role": "owner",
//	           "sex": "male",
//	           "shut_up_timestamp": 0,
//	           "title": "",
//	           "title_expire_time": 0,
//	           "unfriendly": false,
//	           "user_id": 121212121212
//	       },
//	       {
//	           "age": 0,
//	           "area": "",
//	           "card": "",
//	           "card_changeable": false,
//	           "group_id": 1212121212,
//	           "join_time": 1635232488,
//	           "last_sent_time": 1651392760,
//	           "level": "1",
//	           "nickname": "XXXX",
//	           "role": "admin",
//	           "sex": "male",
//	           "shut_up_timestamp": 0,
//	           "title": "",
//	           "title_expire_time": 0,
//	           "unfriendly": false,
//	           "user_id": 12112121212
//	       }
//	   ],
//	   "retcode": 0,
//	   "status": "ok"
//	}
//
//	@return error
func GetGroupMemberList(groupId string) (string, error) {
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

// GetGroupAtAllRemain
//
//	 @Description: 获取群@全体成员 剩余次数
//	 @param groupId 群号
//	 @return string
//		{
//	   "data": {
//	       "can_at_all": true,//是否可以@全体成员
//	       "remain_at_all_count_for_group": 19,//群内所有管理当天剩余 @全体成员 次数
//	       "remain_at_all_count_for_uin": 9 //Bot 当天剩余 @全体成员 次数
//	   },
//	   "retcode": 0,
//	   "status": "ok"
//	}
//
//	@return error
func GetGroupAtAllRemain(groupId string) (string, error) {
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

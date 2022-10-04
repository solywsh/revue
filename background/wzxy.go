package background

import (
	"github.com/solywsh/qqBot-revue/cq"
	"github.com/solywsh/qqBot-revue/service/wzxy"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	MorningCheckEndTime   = "12:00"
	AfternoonCheckEndTime = "18:00"
	EveningCheckEndTime   = "23:59"
)

func wzxyService() {
	for {
		now := time.Now()
		timeNow := now.Format("15:04")
		dateNow := now.Format("2006-01-02")
		wzxyUserMany, _, err := gdb.FindWzxyUserMany(wzxy.UserWzxy{})
		if err != nil {
			log.Println("wzxyService FindWzxyUserMany err:", err)
			time.Sleep(time.Second * 10)
			continue
		}
		monitorWzxyMany, _, err := gdb.FindMonitorWzxyMany(wzxy.MonitorWzxy{})
		if err != nil {
			log.Println("wzxyService FindMonitorWzxyMany err:", err)
			time.Sleep(time.Second * 10)
			return
		}

		// 轮询打卡业务
		for _, userWzxy := range wzxyUserMany {
			if timeNow < MorningCheckEndTime &&
				userWzxy.MorningCheckEnable &&
				userWzxy.MorningCheckTime < timeNow &&
				userWzxy.MorningLastCheckDate < dateNow &&
				userWzxy.JwsessionStatus {
				log.Println("wzxyService 晨检", userWzxy.Name)
				status, msg := userWzxy.CheckOperate(1)
				userId, _ := strconv.Atoi(userWzxy.UserId)
				cpf := cq.PostForm{
					UserId:      userId,
					MessageType: "private",
				}
				if status == 0 {
					cpf.SendMsg("晨检打卡成功")
				} else if msg == "未登录,请重新登录" {
					cpf.SendMsg("晨检打卡失败,可能是jwtsession失效，请尝试wzxy -r 更新jwtsession")
					userWzxy.JwsessionStatus = false
					_, err := gdb.UpdateWzxyUserOne(userWzxy, true)
					if err != nil {
						log.Println("wzxyService UpdateWzxyUserOne err:", err)
						continue
					}
				} else {
					cpf.SendMsg("晨检打卡失败," + msg)
				}
				userWzxy.MorningLastCheckDate = dateNow
				_, err := gdb.UpdateWzxyUserOne(userWzxy, true)
				if err != nil {
					log.Println("wzxyService UpdateWzxyUserOne err:", err)
					continue
				}
			}
			if timeNow < AfternoonCheckEndTime &&
				userWzxy.AfternoonCheckEnable &&
				userWzxy.AfternoonCheckTime < timeNow &&
				userWzxy.AfternoonLastCheckDate < dateNow &&
				userWzxy.JwsessionStatus {
				log.Println("wzxyService 午检", userWzxy.Name)
				status, msg := userWzxy.CheckOperate(2)
				userId, _ := strconv.Atoi(userWzxy.UserId)
				cpf := cq.PostForm{
					UserId:      userId,
					MessageType: "private",
				}
				if status == 0 {
					cpf.SendMsg("午检打卡成功")
				} else if msg == "未登录,请重新登录" {
					cpf.SendMsg("午检打卡失败,可能是jwtsession失效，请尝试wzxy -r 更新jwtsession")
					userWzxy.JwsessionStatus = false
					_, err := gdb.UpdateWzxyUserOne(userWzxy, true)
					if err != nil {
						log.Println("wzxyService UpdateWzxyUserOne err:", err)
						continue
					}
				} else {
					cpf.SendMsg("午检打卡失败," + msg)
				}
				userWzxy.AfternoonLastCheckDate = dateNow
				_, err := gdb.UpdateWzxyUserOne(userWzxy, true)
				if err != nil {
					log.Println("wzxyService UpdateWzxyUserOne err:", err)
					continue
				}
			}
			if timeNow < EveningCheckEndTime &&
				userWzxy.EveningCheckEnable &&
				userWzxy.EveningCheckTime < timeNow &&
				userWzxy.EveningLastCheckDate < dateNow &&
				userWzxy.JwsessionStatus {
				log.Println("wzxyService 晚检", userWzxy.Name)
				status := userWzxy.EveningCheckOperate()
				userId, _ := strconv.Atoi(userWzxy.UserId)
				cpf := cq.PostForm{
					UserId:      userId,
					MessageType: "private",
				}
				if status == 0 {
					cpf.SendMsg("晚检打卡成功")
				} else if status == -1 {
					cpf.SendMsg("晚检打卡失败,可能是网络故障，请尝试wzxy -do 手动打卡")
				} else if status == -2 {
					cpf.SendMsg("晚检打卡失败,可能是jwtsession，请尝试wzxy -r 更新jwtsession")
					userWzxy.JwsessionStatus = false
				} else if status == -3 {
					cpf.SendMsg("晚检打卡失败,可能是不在签到时间范围内，请尝试wzxy -do 手动打卡")
				} else if status == -4 {
					cpf.SendMsg("获取晚检信息失败,未知错误")
				}
				userWzxy.EveningLastCheckDate = dateNow
				_, err := gdb.UpdateWzxyUserOne(userWzxy, true)
				if err != nil {
					log.Println("wzxyService UpdateWzxyUserOne err:", err)
					continue
				}
			}
		}
		// 轮询打卡提醒业务
		for _, monitorWzxy := range monitorWzxyMany {
			var userWzxy wzxy.UserWzxy
			for _, wum := range wzxyUserMany {
				if wum.ID == monitorWzxy.UserWzxyId {
					userWzxy = wum
				}
			}
			if timeNow < MorningCheckEndTime &&
				monitorWzxy.MorningRemindEnable &&
				monitorWzxy.MorningRemindTime < timeNow &&
				monitorWzxy.MorningRemindLastDate < dateNow &&
				userWzxy.JwsessionStatus {
				handleRemindCheckDaily(1, dateNow, monitorWzxy, userWzxy)
			}
			if timeNow < AfternoonCheckEndTime &&
				monitorWzxy.AfternoonRemindEnable &&
				monitorWzxy.AfternoonRemindTime < timeNow &&
				monitorWzxy.AfternoonRemindLastDate < dateNow &&
				userWzxy.JwsessionStatus {
				handleRemindCheckDaily(2, dateNow, monitorWzxy, userWzxy)
			}
			if monitorWzxy.CheckRemindEnable &&
				monitorWzxy.CheckRemindTime < timeNow &&
				monitorWzxy.CheckRemindLastDate < dateNow &&
				userWzxy.JwsessionStatus {
				handleRemindSign(dateNow, monitorWzxy, userWzxy)
			}
		}
		time.Sleep(5 * time.Minute)
	}
}

func handleRemindCheckDaily(seq int, dateNow string, monitorWzxy wzxy.MonitorWzxy, userWzxy wzxy.UserWzxy) {
	var keywords string
	if seq == 1 {
		keywords = "晨检"
	} else {
		keywords = "午检"
	}
	userId, _ := strconv.Atoi(userWzxy.UserId)
	groupId, _ := strconv.Atoi(monitorWzxy.ClassGroupId)
	cpf := cq.PostForm{
		UserId:      userId,
		GroupId:     groupId,
		MessageType: "private", // private group
	}
	log.Println("class name:", monitorWzxy.ClassName,
		"user name:", userWzxy.Name,
		"seq:", seq,
		keywords+"wzxyService 打卡提醒")
	uncheckList, err := userWzxy.GetDailyUncheckList(seq)
	if err != nil {
		if strings.Contains(err.Error(), "未登录") {
			log.Println("class name:", monitorWzxy.ClassName,
				"user name:", userWzxy.Name,
				"seq:", seq,
				"wzxyService GetDailyUncheckList err:", err)
			cpf.SendMsg("获取晨检未打卡列表失败,可能是jwtsession失效，请尝试wzxy -r 更新jwtsession")
			userWzxy.JwsessionStatus = false
			_, err = gdb.UpdateWzxyUserOne(userWzxy, true)
			if err != nil {
				log.Println("class name:", monitorWzxy.ClassName,
					"user name:", userWzxy.Name,
					"seq:", seq,
					"wzxyService UpdateWzxyUserOne err:", err)
				return
			}
		}
		log.Println("class name:", monitorWzxy.ClassName,
			"user name:", userWzxy.Name,
			"seq:", seq,
			"wzxyService FindClassStudentWzxyMany:", err)
		cpf.SendMsg("获取晨检未打卡列表失败")
		return
	}
	if len(uncheckList) == 0 {
		cpf.SendGroupMsg("今天所有人都已经打卡了")
		return
	}
	var msg string
	msg += keywords + "未打卡列表:\n"
	for _, uncheck := range uncheckList {
		many, i, err := gdb.FindClassStudentWzxyMany(wzxy.ClassStudentWzxy{StudentId: uncheck.StudentId})
		if err != nil {
			log.Println("class name:", monitorWzxy.ClassName,
				"user name:", userWzxy.Name,
				"seq:", seq,
				"uncheck name:", uncheck.Name,
				"wzxyService FindClassStudentWzxyMany:", err)
			continue
		} else if i == 1 {
			msg += cq.GetCqCodeAt(many[0].UserId, "") + " "
		} else if i == 0 {
			msg += uncheck.Name + "(未添加至数据库) "
		}
	}
	msg += "\n请尽快打卡"
	cpf.SendGroupMsg(msg)
	if seq == 1 {
		monitorWzxy.MorningRemindLastDate = dateNow
	} else {
		monitorWzxy.AfternoonRemindLastDate = dateNow
	}
	_, err = gdb.UpdateMonitorWzxyOne(monitorWzxy, true)
	if err != nil {
		log.Println("class name:", monitorWzxy.ClassName,
			"user name:", userWzxy.Name,
			"seq:", seq,
			"uncheck name:",
			"wzxyService UpdateWzxyMonitorOne err:", err)
		return
	}
}

func handleRemindSign(dateNow string, monitorWzxy wzxy.MonitorWzxy, userWzxy wzxy.UserWzxy) {
	userId, _ := strconv.Atoi(userWzxy.UserId)
	groupId, _ := strconv.Atoi(monitorWzxy.ClassGroupId)
	cpf := cq.PostForm{
		UserId:      userId,
		GroupId:     groupId,
		MessageType: "private", // private group
	}
	log.Println("class name:", monitorWzxy.ClassName,
		"user name:", userWzxy.Name,
		"wzxyService 签到提醒")
	uncheckList, status := userWzxy.GetUnSignedList()
	if status == -1 {
		cpf.SendMsg("获取未签到列表信息失败,网络错误")
		return
	} else if status == -3 {
		cpf.SendMsg("获取未签到列表失败,不在签到时间范围内")
		return
	} else if status == -4 {
		cpf.SendMsg("获取未签到列表失败,可能是jwtsession失效，请尝试wzxy -r 更新jwtsession")
		userWzxy.JwsessionStatus = false
		_, err := gdb.UpdateWzxyUserOne(userWzxy, true)
		if err != nil {
			log.Println("class name:", monitorWzxy.ClassName,
				"user name:", userWzxy.Name,
				"wzxyService UpdateWzxyUserOne err:", err)
		}
		return
	} else if status == -5 {
		cpf.SendMsg("获取签到列表失败,未知的错误")
		return
	}
	if len(uncheckList) == 0 {
		cpf.SendGroupMsg("今天所有人都已经打卡了")
		return
	}
	var msg string
	msg += "未签到列表:\n"
	for _, uncheck := range uncheckList {
		many, i, err := gdb.FindClassStudentWzxyMany(wzxy.ClassStudentWzxy{StudentId: uncheck.StudentId})
		if err != nil {
			log.Println("class name:", monitorWzxy.ClassName,
				"user name:", userWzxy.Name,
				"uncheck name:", uncheck.Name,
				"wzxyService FindClassStudentWzxyMany:", err)
			continue
		} else if i == 1 {
			msg += cq.GetCqCodeAt(many[0].UserId, "") + " "
		} else if i == 0 {
			msg += uncheck.Name + "(未添加至数据库) "
		}
	}
	msg += "\n请尽快签到"
	cpf.SendGroupMsg(msg)
	monitorWzxy.CheckRemindLastDate = dateNow
	_, err := gdb.UpdateMonitorWzxyOne(monitorWzxy, true)
	if err != nil {
		log.Println("class name:", monitorWzxy.ClassName,
			"user name:", userWzxy.Name,
			"wzxyService UpdateWzxyMonitorOne err:", err)
		return
	}
}

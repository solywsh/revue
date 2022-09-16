package background

import (
	"github.com/solywsh/qqBot-revue/cq"
	"github.com/solywsh/qqBot-revue/service/wzxy"
	"log"
	"strconv"
	"time"
)

const (
	MorningCheckEndTime   = "10:00"
	AfternoonCheckEndTime = "15:00"
	EveningCheckEndTime   = "23:59"
)

func wzxyService() {
	for {
		now := time.Now()
		timeNow := now.Format("15:04")
		dateNow := now.Format("2006-01-02")
		many, _, err := gdb.FindWzxyUserMany(wzxy.UserWzxy{})
		if err != nil {
			log.Println("wzxyService FindWzxyUserMany err:", err)
			continue
		}
		for _, userWzxy := range many {
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
		time.Sleep(5 * time.Minute)
	}
}

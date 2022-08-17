package main

import (
	"context"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/go-resty/resty/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

var (
	MaxClient    = 40 // 最大并发数
	client       = resty.New()
	URL          = "https://api.lolicon.app/setu/v2"                     // 涩图 URL
	DbUrl        = "mongodb://revue:himgrevue@81.68.236.195:27017/revue" // mongo URL
	DbClient     *mongo.Client
	DbCollection *mongo.Collection
	LccCh        = make(chan *Lolicon, 2*MaxClient)
	LimitCh      = make(chan int, MaxClient+1)
	finalCh      = make(chan int) // 结束chanel
	allNumCh     = make(chan int) // 爬取的涩图数量
	missNumCh    = make(chan int) // 重复命中数量
)

type Lolicon struct {
	Error string        `json:"error"`
	Data  []LoliconData `json:"data"`
}

type LoliconData struct {
	Pid        int      `json:"pid"`
	P          int      `json:"p"`
	Uid        int      `json:"uid"`
	Title      string   `json:"title"`
	Author     string   `json:"author"`
	R18        bool     `json:"r18"`
	Width      int      `json:"width"`
	Height     int      `json:"height"`
	Tags       []string `json:"tags"`
	Ext        string   `json:"ext"`
	UploadDate int64    `json:"uploadDate"`
	Urls       struct {
		Original string `json:"original"`
	} `json:"urls"`
}

func init() {
	// Set DbClient options
	clientOptions := options.Client().ApplyURI(DbUrl)
	// Connect to MongoDB
	var err error
	DbClient, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// Check the connection
	err = DbClient.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	DbCollection = DbClient.Database("revue").Collection("himg")
	fmt.Println("Connected to MongoDB!")
}

func getInfo() {
	for i := 0; i < MaxClient; i++ {
		go func() {
			for {
				select {
				case <-LimitCh:
					get, err := client.R().SetQueryParams(map[string]string{
						"r18": "2",
						"num": "20",
					}).Get(URL)
					if err != nil {
						return
					}
					var llc *Lolicon
					err = yaml.Unmarshal(get.Body(), &llc)
					if err != nil {
						fmt.Println("文件解析错误!", err)
						return
					}
					LccCh <- llc
				}
			}
		}()

	}
}

func saveInfo() {
	LimitCh <- 1 // 启动协程
	for i := 0; i < MaxClient; i++ {
		LimitCh <- 1
		go func() {
			for {
				select {
				case llc := <-LccCh:
					for _, datum := range llc.Data {
						var res LoliconData
						err := DbCollection.FindOne(context.TODO(), bson.M{"pid": datum.Pid}).Decode(&res)
						if err != nil {
							_, err := DbCollection.InsertOne(context.TODO(), datum)
							if err != nil {
								fmt.Println("insert error:", err)
								continue
							}
							allNumCh <- 1
						} else {
							missNumCh <- 1
						}
					}
					LimitCh <- 1
				}
			}
		}()
	}

}

func count() {
	var allNum int
	var missNum int
	for {
		select {
		case <-missNumCh:
			allNum++
			missNum++
			res := float64((allNum-missNum)*100) / float64(allNum)
			fmt.Println("data had!", fmt.Sprintf("命中率:%f %%", res))
			if res <= 0.00001 && allNum > 10000 {
				finalCh <- 1 // 发送结束信号
			}
			continue
		case <-allNumCh:
			allNum++
			res := float64((allNum-missNum)*100) / float64(allNum)
			fmt.Println("insert success!!", fmt.Sprintf("命中率:%f %%", res))

			continue
		}
	}
}

func main() {

	client.SetRetryCount(3)
	client.SetRetryWaitTime(time.Second * 10)

	// 关闭数据库
	defer func(DbClient *mongo.Client, ctx context.Context) {
		err := DbClient.Disconnect(ctx)
		if err != nil {
			fmt.Println("Disconnect error:", err)
		}
	}(DbClient, context.TODO())

	go count()    // 启动统计协程
	go getInfo()  // 启动爬取协程
	go saveInfo() // 启动保存协程

	// 阻塞主线程
	for {
		select {
		// 可以给条件结束,
		//如case <-time.After(time.Second * 10):
		//也可以等命中率达到0.000001的时候发送一个chanel结束
		case <-finalCh:
			fmt.Println("结束")
			os.Exit(0)
		}
	}
}

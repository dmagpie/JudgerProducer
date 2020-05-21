package router

import (
	"JudgerProducer/config"
	"JudgerProducer/msgq"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

// Router 路由
var Router *gin.Engine

// 请求格式:
/*
   {
       "id": "主服务器上的提交ID 此项会在回调时提供",
	   "code": "用户的代码", "//": "base64编码",
	   "language": "语言",
       "limits": {
           "cpu": "1000", "//": "CPU时间限制 单位为ms. 用户进程将会在相当于2.25倍CPU时间限制的实际时间后被杀死.",
           "memory": "", "//": "内存使用限制 单位为Byte 实际内存限制会是这里传入的2倍, 请在随后传回的内存用量中自行判断是否MLE",
       },
       "data": [
           "",
           ""
       ], "//": "输入数据. 必须为字符串形式."
   }
*/
type submitData struct {
	ID       string   `form:"id" binding:"required" json:"id"`
	Code     string   `form:"code" binding:"required" json:"code"`
	Language string   `form:"language" binding:"required" json:"language"`
	Limits   *limits  `form:"limits" binding:"required" json:"limits"`
	Data     []string `form:"data" binding:"required" json:"data"`
}

type limits struct {
	CPU    uint64 `form:"cpu" binding:"required" json:"cpu"`
	Memory uint64 `form:"memory" binding:"required" json:"memory"`
}

func init() {
	Router = gin.Default()

	Router.Use(BasicAuth())

	{
		// ping请求
		Router.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"msg": "pong",
			})
		})
	}

	{
		// 提交评测代码以及数据
		// 成功提交返回格式:
		/*
		   {
		       "id": "传入的id",
		       "status": "ok"
		   }
		*/
		// 提交失败的返回格式(如格式错误)
		/*
		   {
		       "id": "传入的id",
		       "status": "invalid-synatx"
		   }
		*/
		Router.POST("/submit", func(c *gin.Context) {
			qname, err := config.GetConfig("RMQ_QNAME")
			if err != nil {
				c.JSON(400, gin.H{
					"errmsg": "config-error",
				})
				return
			}
			data := submitData{}
			err = c.BindJSON(&data)
			if err != nil {
				c.JSON(400, gin.H{
					"errmsg": "form-error",
				})
				return
			}
			channel, err := msgq.Conn.Channel()
			if err != nil {
				c.JSON(400, gin.H{
					"errmsg": "mq-error-1",
				})
				return
			}

			q, err := channel.QueueDeclare(qname, true, false, false, false, nil)
			if err != nil {
				c.JSON(400, gin.H{
					"errmsg": "mq-error-2",
				})
				return
			}

			jdata, err := json.Marshal(data)
			if err != nil {
				c.JSON(400, gin.H{
					"errmsg": "json-error",
				})
				return
			}
			err = channel.Publish("", q.Name, false, false, amqp.Publishing{
				DeliveryMode: amqp.Persistent, //Msg set as persistent
				ContentType:  "application/json",
				Body:         jdata,
			})
			if err != nil {
				c.JSON(400, gin.H{
					"errmsg": "mq-error-3",
				})
				return
			}

			c.JSON(204, gin.H{
				"errmsg": "success",
			})
		})
	}

}

package server

import (
	"fmt"
	"github.com/XIU2/CloudflareSpeedTest/utils"
	"github.com/gin-gonic/gin"
	"strconv"
	"sync"
)

var Service *gin.Engine
var speedData utils.DownloadSpeedSet
var mu sync.Mutex

func init() {
	Service = gin.Default()
}
func SetSpeedData(speed utils.DownloadSpeedSet) {
	mu.Lock()
	speedData = speed
	mu.Unlock()
}
func StartServer() {
	Service.GET("/bestip", func(c *gin.Context) {
		// 检查speedData是否为空或没有数据
		if speedData == nil || len(speedData) == 0 {
			c.JSON(200, gin.H{"error": "no data"})
			return
		}

		// 创建一个slice，用于存放所有speedData记录的信息
		var speedDataList []map[string]interface{}

		// 遍历speedData，将每个记录的信息添加到speedDataList中
		for _, v := range speedData {
			record := map[string]interface{}{
				"ip":            v.PingData.IP.IP,
				"post":          v.PingData.Sended,
				"received":      v.PingData.Received,
				"loss":          v.LossRate,
				"delay":         v.PingData.Delay.Milliseconds(),
				"downloadSpeed": strconv.FormatFloat(v.DownloadSpeed/1024/1024, 'f', 2, 32),
			}
			speedDataList = append(speedDataList, record)
		}

		// 使用c.JSON发送整个speedDataList
		c.JSON(200, speedDataList)
	})
	err := Service.Run(":8080")
	if err != nil {
		fmt.Println("Failed to start server:", err)
	}
}

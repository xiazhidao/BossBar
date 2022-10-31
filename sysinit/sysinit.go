package sysinit

import (
	"BossBar/utils"
	"math/rand"
	"time"

	log "github.com/sirupsen/logrus"
)

func init() {
	//随机种子
	rand.Seed(time.Now().UnixNano())

	utils.InitLogs()

	//初始化缓存
	utils.InitCache()

	log.Info("Initialized is done~")
}

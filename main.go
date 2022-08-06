package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	observerbot "github.com/kaankoken/binance-quick-go/observer-bot"
	log "github.com/sirupsen/logrus"

	"github.com/robfig/cron/v3"
)

func init() {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
}

func main() {
	log.Info("Create new cron")

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	go func() {
		c := cron.New(cron.WithChain(
			cron.Recover(cron.DefaultLogger),
		))

		c.AddFunc("0 */1 * * *", func() {
			log.Info("Running Hourly EMA Calculation")
			observerbot.Run()
		})

		c.AddFunc("1 0 */1 * *", func() {
			log.Info("Running Daily Volume Calculation")
			observerbot.CalculateFiftyDaysOfAverageVolume()
		})

		log.Info("Start cron")
		c.Start()
	}()

	r.Run()
}

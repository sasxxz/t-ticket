package task

import (
	"context"
	"fmt"
	"log"
	"time"
	"udesk/api"
	"udesk/server"

	"github.com/robfig/cron/v3"
)

func CronTask() {

	c := cron.New(cron.WithSeconds())
	_, err := c.AddFunc("0 0 18 * * 1-5", func() {
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			for {
				if time.Now().Format("15:04:05") == "06:00:00" {
					server.ShutdownServer.Shutdown(ctx)
					log.Println("server shutdown..")
					cancel()
					return
				}
			}
		}()
		fmt.Println("cronjob start")
		api.UdeskApi()
		go server.Server()
		go func() {
			for {
				api.ReAgent()
			}
		}()

	})
	if err != nil {
		log.Fatalf("cron job failed,err:%v\n", err)
	}
	c.Start()

}

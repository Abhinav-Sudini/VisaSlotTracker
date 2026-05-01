package main

import (
	"log"
	"time"
)


func start_cron_job(interval int) {
	ticker := time.Tick(time.Second * time.Duration(interval))

	for {
		select {
		case _ = <-ticker:
			SaveSlotsAndNotifyUsers()
			log.Println("[INFO] [done] Saved new slots data")
			log.Println("")	
		case <-cron_stop_chan:
			log.Println("cron job handler stoped")
			return
		}
	}
}

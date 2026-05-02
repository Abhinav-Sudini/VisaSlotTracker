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

func restart_cron_jobs_everyday() {
	ticker := time.Tick(time.Hour * time.Duration(24))

	for {
		select {
		case _ = <-ticker:
			restartCron((24*60*60)/20)
			log.Println("[INFO] restarting cron to 20")
			log.Println("")	
		case <-restarter_stop_chan:
			log.Println("cron job handler stoped")
			return
		}
	}
}

func restartCron(t int) {
	cron_stop_chan<-true
	log.Println("restarting cron handler with int - ",t)
	go start_cron_job(config.ticker_time)
}

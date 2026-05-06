package main

import (
	"log"
	"os"
	"path/filepath"
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
			removeOldData()
			reset_notify()
			restartCron((24 * 60 * 60) / 20)
			log.Println("[INFO] restarting cron to 20")
			log.Println("")
		case <-restarter_stop_chan:
			log.Println("cron job handler stoped")
			return
		}
	}
}

func reset_notify(){
	for i := range n_list{
		n_list[i].valid = true
	}
}

func removeOldData() {
	rmdate := time.Now().AddDate(0, 0, -2).Format("2006-01-02")
	DirLoc := filepath.Join(config.FS_root, rmdate)
	os.RemoveAll(DirLoc)
}

func restartCron(t int) {
	cron_stop_chan <- true
	log.Println("restarting cron handler with int - ", t)
	go start_cron_job(config.ticker_time)
}

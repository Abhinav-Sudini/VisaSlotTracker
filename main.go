package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var cwd, _ = os.Getwd()
var cron_stop_chan = make(chan bool)
var restarter_stop_chan = make(chan bool)

var n_list = []notify_info{
	{
		notify_type:        "mail",
		mail:               "sudiniabhinav@gmail.com",
		visa_location:      "HYDERABAD",
		min_slots_required: 2,
	},
	{
		notify_type:        "mail",
		mail:               "sudiniabhinav@gmail.com",
		visa_location:      "HYDERABAD VAC",
		min_slots_required: 2,
	},
	{
		notify_type:        "mail",
		mail:               "akshithasudini@gmail.com",
		visa_location:      "HYDERABAD",
		min_slots_required: 4,
	},
	{
		notify_type:        "mail",
		mail:               "akshithasudini@gmail.com",
		visa_location:      "HYDERABAD VAC",
		min_slots_required: 4,
	},
	{
		notify_type:        "mail",
		mail:               "abhinav_2301cs03@iitp.ac.in",
		visa_location:      "HYDERABAD",
		min_slots_required: 2,
	},
	{
		notify_type:        "mail",
		mail:               "abhinav_2301cs03@iitp.ac.in",
		visa_location:      "HYDERABAD VAC",
		min_slots_required: 2,
	},
	{
		notify_type:        "mail",
		mail:               "abhinav_2301cs03@iitp.ac.in",
		visa_location:      "*",
		min_slots_required: 10,
	},
	{
		notify_type:        "mail",
		mail:               "sudiniabhinav@gmail.com",
		visa_location:      "*",
		min_slots_required: 10,
	},
}

var config config_struct = config_struct{
	max_tries:   20,
	ticker_time: (24 * 60 * 60) / 20,
	// ticker_time: 10,
	notify_list: n_list,
	FS_root:     filepath.Join(cwd, "img_store"),
	ServerAddr:  ":8000",
	log_file:    "app.log",
	Api_key:     "DJO8V2",
}

func Start_server() {
	http.HandleFunc("/", getHomeHandler)
	http.HandleFunc("/logs/", getLogsHandler)
	http.HandleFunc("/notify_sub/", getNotifyUsersHandler)
	http.HandleFunc("/test_notify/{mail}/", testNotifyHandler)
	http.HandleFunc("/api/info/{day}/{visa_loc}/", infoHandler)
	http.HandleFunc("/api/update_tic_time/{new_time}/", ticTimeHandler)
	http.HandleFunc("/api/update_tic_cnt/{new_time}/", ticCntHandler)

	fshandler := http.FileServer(http.Dir(config.FS_root))
	http.Handle("/fs/", http.StripPrefix("/fs/", fshandler))

	log.Println("starting server at -", config.ServerAddr)
	log.Println("cron job time interval - ", config.ticker_time)
	_ = http.ListenAndServe(config.ServerAddr, nil)
}

func main() {
	file, _ := os.OpenFile(config.log_file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	defer file.Close()
	log.SetOutput(file)
	// log.SetOutput(os.Stdout)

	go restart_cron_jobs_everyday()
	go start_cron_job(config.ticker_time)
	Start_server()
}


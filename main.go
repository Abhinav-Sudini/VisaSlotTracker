package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
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
		mail:               "sudiniabhinav@gmail.com",
		visa_location:      "*",
		min_slots_required: 4,
	},
	{
		notify_type:        "mail",
		mail:               "akshithasudini@gmail.com",
		visa_location:      "*",
		min_slots_required: 6,
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

	fmt.Println(GetMaxTries())
	go restart_cron_jobs_everyday()
	go start_cron_job(config.ticker_time)
	Start_server()
}

func GetMaxTries() int {
	cmd := exec.Command(
		"curl", "-s",
		"https://app.checkvisaslots.com/validate/v3",
		"-X", "GET",
		"-H", "accept: */*",
		"-H", "accept-encoding: gzip, deflate, br, zstd",
		"-H", "accept-language: en-US,en;q=0.9,te;q=0.8,es;q=0.7",
		"-H", "origin: https://checkvisaslots.com",
		"-H", "referer: https://checkvisaslots.com/",
		"-H", `sec-ch-ua: "Google Chrome";v="147", "Not.A/Brand";v="8", "Chromium";v="147"`,
		"-H", "sec-ch-ua-mobile: ?0",
		"-H", `sec-ch-ua-platform: "Windows"`,
		"-H", "sec-fetch-dest: empty",
		"-H", "sec-fetch-mode: cors",
		"-H", "sec-fetch-site: same-site",
		"-H", "user-agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
		"-H", "x-api-key: DJO8V2",
		"--compressed",
	)

	type Activity struct {
		Retrieve int `json:"retrieve"`
		Slots    int `json:"slots"`
		Upload   int `json:"upload"`
	}
	type UserActivityOnly struct {
		Activity Activity `json:"activity"`
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
	}

	var data UserActivityOnly
	err = json.Unmarshal(output, &data)
	if err != nil {
		log.Println("http getter : err in parsing json - ", err)
	}
	// fmt.Println(data)

	out := max(0,20 * data.Activity.Upload - (data.Activity.Retrieve + data.Activity.Slots) + 5)
	return out
}

func getSlotsjson() ([]Slot, error) {
	cmd := exec.Command(
		"curl", "-s",
		"https://app.checkvisaslots.com/retrieve/v1",
		"-X", "GET",
		"-H", "accept: */*",
		"-H", "accept-encoding: gzip, deflate, br, zstd",
		"-H", "accept-language: en-US,en;q=0.9,te;q=0.8,es;q=0.7",
		"-H", "origin: https://checkvisaslots.com",
		"-H", "referer: https://checkvisaslots.com/",
		"-H", `sec-ch-ua: "Google Chrome";v="147", "Not.A/Brand";v="8", "Chromium";v="147"`,
		"-H", "sec-ch-ua-mobile: ?0",
		"-H", `sec-ch-ua-platform: "Windows"`,
		"-H", "sec-fetch-dest: empty",
		"-H", "sec-fetch-mode: cors",
		"-H", "sec-fetch-site: same-site",
		"-H", "user-agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
		"-H", fmt.Sprintf("x-api-key: %v", config.Api_key),
		"--compressed",
	)
	var data []Slot

	output, err := cmd.Output()
	if err != nil {
		log.Println("http getter : err in geting curl req - ", err)
		return data, err
	}

	err = json.Unmarshal(output, &data)
	if err != nil {
		log.Println("http getter : err in parsing json - ", err)
	}
	return data, err
}

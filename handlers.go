package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
)

func getHomeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./index.html")
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	Day := r.PathValue("day")
	location := r.PathValue("visa_loc")

	resp, _ := GetFilesAsJSON(filepath.Join(config.FS_root, Day, location))
	if resp == "" {
		resp = "[]"
	}

	// fmt.Println(resp,Day,location,filepath.Join(cwd,Day,location))
	w.Write([]byte(resp))
}

func getLogsHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./app.log")
}

func getNotifyUsersHandler(w http.ResponseWriter, r *http.Request) {
	out := ""
	for _, n_info := range config.notify_list {
		out += fmt.Sprintf(
			"action=notify mail_to= %s || min_open_slots_required_to_trigger= %v || visa_location=%s || valid %v \n",
			n_info.mail,
			n_info.min_slots_required,
			n_info.visa_location,
			n_info.valid,
		)
	}
	w.Write([]byte(out))
	fmt.Fprintf(w,"ticker time = %v",config.ticker_time)
}

func testNotifyHandler(w http.ResponseWriter, r *http.Request) {
	mail := r.PathValue("mail")
	if mail == ""{
		http.Error(w,"no mail given",http.StatusBadRequest)
		return
	}
	SendMail(mail,1,"HYDERABAD")
}

func ticTimeHandler(w http.ResponseWriter, r *http.Request) {
	new_time := r.PathValue("new_time")
	v,err := strconv.Atoi(new_time)
	if new_time == "" || err != nil{
		http.Error(w,"no mail given",http.StatusBadRequest)
		return
	}
	config.ticker_time = v
	restartCron(config.ticker_time)
}

func ticCntHandler(w http.ResponseWriter, r *http.Request) {
	new_time := r.PathValue("new_time")
	v,err := strconv.Atoi(new_time)
	if new_time == "" || err != nil{
		http.Error(w,"no mail given",http.StatusBadRequest)
		return
	}
	config.max_tries = v
	config.ticker_time = (24*60*60)/v
	fmt.Fprintf(w,"restarting cron with tries %v  and time %v \n",config.max_tries,config.ticker_time)
	restartCron(config.ticker_time)
}

func toggle_valid_users(w http.ResponseWriter, r *http.Request) {
	n_id := r.PathValue("n_id")
	id,err := strconv.Atoi(n_id)
	if n_id == "" || err != nil || id > len(n_list){
		http.Error(w,"not valid id given",http.StatusBadRequest)
		return
	}
	n_list[id].valid = toggle(n_list[id].valid)
	fmt.Fprintf(w,"toggle %v  \n",n_list[id])
}

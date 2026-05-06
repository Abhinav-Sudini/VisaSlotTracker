package main

type config_struct struct {
	max_tries   int
	ticker_time int
	notify_list []notify_info
	FS_root     string
	ServerAddr  string
	log_file    string
	Api_key     string
}

type notify_info struct {
	notify_type        string
	mail               string
	visa_location      string
	min_slots_required int
	valid              bool
}

type Slot struct {
	CreatedOn    string `json:"createdon"`
	ImgURL       string `json:"img_url"`
	Slots        int    `json:"slots"`
	VisaLocation string `json:"visa_location"`
}

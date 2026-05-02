package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)


func SaveSlotsAndNotifyUsers() error {
	slots, err := getSlotsjson()
	if err != nil {
		log.Println("failed to get slots with err = ", err)
		return err
	}
	if len(slots)==0 {
		config.ticker_time = config.ticker_time*2
		restartCron(config.ticker_time)
	}
	saveSlots(slots)

	NotifyUsers(slots)
	return nil
}


func NotifyUsers(slots []Slot) {
	for _,n_info := range(config.notify_list){
		for _,slot := range(slots){
			if n_info.visa_location!="*" && !strings.EqualFold(slot.VisaLocation,n_info.visa_location){
				continue
			}
			if slot.Slots >= n_info.min_slots_required {
				log.Println("")
				log.Println("Sending mail to - ",n_info.mail)
				SendMail(n_info.mail,slot.Slots,slot.VisaLocation)
			}
		}
	}
}


func saveSlots(slots []Slot) {
	saveSlot := func(slot Slot, wg *sync.WaitGroup) error {
		// fmt.Println(slot.Slots,slot.VisaLocation,slot.CreatedOn)
		defer wg.Done()

		t, err := time.Parse(time.RFC1123, slot.CreatedOn)
		if err != nil {
			fmt.Println("time parese err for - ", slot)
			return err
		}

		// Create date dir (YYYY-MM-DD)
		dateDir := t.Format("2006-01-02")


		// Clean visa location (avoid spaces issues)
		locDir := strings.ReplaceAll(slot.VisaLocation, " ", "_")

		// Create full path
		fullDir := filepath.Join(config.FS_root, dateDir, locDir)

		// Create directories
		err = os.MkdirAll(fullDir, os.ModePerm)
		if err != nil {
			fmt.Println("dir create err for slot ", slot)
			return err
		}

		// File name = time (HH-MM-SS)
		fileName := t.Format("15-04-05") + ".png"

		filePath := filepath.Join(fullDir, fileName)

		err = SaveUrlBodyToFile(slot.ImgURL, filePath)
		if err != nil {
			fmt.Println("failed to save img errr -", err)
			return err
		}

		return nil
	}

	var wg sync.WaitGroup
	for _, slot := range slots {
		wg.Add(1)
		err := saveSlot(slot, &wg)
		if err != nil {
			fmt.Println("err in download ",err)
		}
	}
	wg.Wait()
}


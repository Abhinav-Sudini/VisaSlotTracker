package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
)



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

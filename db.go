package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

func GetAllDays() []DayInfo {
	dbFile, err := os.Open("db.json")

	if err != nil {
		log.Fatal(err)
	}
	defer dbFile.Close()

	dbContent, _ := ioutil.ReadAll(dbFile)
	var days []DayInfo

	jerr := json.Unmarshal(dbContent, &days)
	if jerr != nil {
		log.Fatal(jerr)
	}

	return days
}

func GetCurrentDay() DayInfo {
	days := GetAllDays()
	return days[len(days)-1]
}

func UpdateHour(hour int, qty int) HourInfo {
	days := GetAllDays()
	currentDay := &days[len(days)-1]
	var index = 0
	var accumQty = 0

	for i, h := range currentDay.Hours {
		if h.Hour == hour {
			index = i
		}
		accumQty += h.Qty
	}
	accumQty += qty

	currentDay.Hours[index].Qty = qty
	currentDay.Hours[index].Completed = true
	currentDay.Overflow = max(accumQty-currentDay.MaxQty, 0)

	newDBContent, _ := json.Marshal(days)
	err := os.WriteFile("db.json", newDBContent, 0644)

	if err != nil {
		log.Fatal(err)
	}

	return currentDay.Hours[index]
}

func GetOverflow() int {
	day := GetCurrentDay()
	return day.Overflow
}

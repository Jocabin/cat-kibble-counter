package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"
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

	if len(days) == 0 {
		newDay := AddNewDay(time.Now())
		return newDay
	}

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

func AddNewDay(date time.Time) DayInfo {
	days := GetAllDays()
	newDay := DayInfo{
		Name:        date.Weekday().String(),
		BaseQty:     10,
		MaxQty:      60,
		Overflow:    0,
		Date:        date.Format("02/01/2006"),
		MorningMash: false,
		EveningMash: false,
		Hours: []HourInfo{
			{Hour: 6, Qty: 0, Completed: false},
			{Hour: 9, Qty: 0, Completed: false},
			{Hour: 12, Qty: 0, Completed: false},
			{Hour: 15, Qty: 0, Completed: false},
			{Hour: 18, Qty: 0, Completed: false},
			{Hour: 21, Qty: 0, Completed: false},
		},
	}

	days = append(days, newDay)

	newDBContent, _ := json.Marshal(days)
	err := os.WriteFile("db.json", newDBContent, 0644)

	if err != nil {
		log.Fatal(err)
	}

	return newDay
}

func DeleteHour(hour int) HourInfo {
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
	accumQty -= currentDay.Hours[index].Qty

	currentDay.Hours[index].Qty = 0
	currentDay.Hours[index].Completed = false
	currentDay.Overflow = max(accumQty-currentDay.MaxQty, 0)

	newDBContent, _ := json.Marshal(days)
	err := os.WriteFile("db.json", newDBContent, 0644)

	if err != nil {
		log.Fatal(err)
	}

	return currentDay.Hours[index]
}

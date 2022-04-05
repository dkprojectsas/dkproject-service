package utils

import (
	"dk-project-service/entity"
	"time"
)

func CountRangeDate() (string, string, error) {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return "", "", err
	}

	timeLocal := time.Now().In(loc)

	dateRange := int(timeLocal.Weekday())

	if dateRange == 0 {
		start := timeLocal.Add(-6 * time.Hour * 24)
		end := timeLocal.Add(time.Duration(dateRange) * time.Hour * 24)

		startDate := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
		endDate := time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, end.Location())

		startDateStr := startDate.Format(entity.ParseFormat)
		endDateStr := endDate.Format(entity.ParseFormat)

		return startDateStr, endDateStr, nil
	} else {
		start := timeLocal.Add(time.Duration(-dateRange+1) * time.Hour * 24)
		end := timeLocal.Add(time.Duration(7-dateRange) * time.Hour * 24)

		startDate := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
		endDate := time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, end.Location())

		startDateStr := startDate.Format(entity.ParseFormat)
		endDateStr := endDate.Format(entity.ParseFormat)

		return startDateStr, endDateStr, err
	}
}

package myutils

import "time"

func GetNowMonthMaxDay() int64 {
	now := time.Now()
	year := now.Year()
	var month int = int(now.Month())
	return int64(GetMonthMaxDay(year, month))
}

func GetNextMonthMaxDay() int64 {
	now := time.Now()
	year := now.Year()
	var month int = int(now.Month())
	if month == 12 {
		year++
		month = 1
	} else {
		month++
	}
	return int64(GetMonthMaxDay(year, month))
}

func GetMonthMaxDay(year, month int) int {
	switch month {
	case 1:
		return 31
	case 3:
		return 31
	case 5:
		return 31
	case 7:
		return 31
	case 8:
		return 31
	case 10:
		return 31
	case 12:
		return 31
	case 2:
		{
			if year%4 == 0 && year%100 != 0 || year%400 == 0 {
				return 29
			} else {
				return 28
			}
		}
	default:
		{
			return 30
		}
	}
	return 0
}

package utils

import "time"

const (
	Hour  = time.Hour
	Day   = time.Duration(24) * Hour
	Month = time.Duration(30) * Day
	Year  = time.Duration(12) * Month
)

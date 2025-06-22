package util

import "time"

func ParseTime(ts string) (string, error) {
	layout := "2006-01-02 15:04:05 -0700 MST"
	pt, err := time.Parse(layout, ts)
	if err != nil {
		return "", err
	}

	return pt.Format("2006-01-02 15:04"), nil
}

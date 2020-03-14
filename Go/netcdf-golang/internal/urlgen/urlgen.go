package urlgen

import (
	"fmt"
	"time"
)

// NameSet return an []string of names to attempt to pull down
func NameSet(s, e time.Time) []string {

	var urls []string

	for rd := rangeDate(s, e); ; {
		date := rd()
		if date.IsZero() {
			break
		}
		d := fmt.Sprint(date.Format("20060102"))
		for i := 0; i <= 23; i++ {
			t := fmt.Sprintf("%02d", i)
			urls = append(urls, fmt.Sprintf("https://storage.cloud.google.com/national-water-model/nwm.%s/analysis_assim/nwm.t%sz.analysis_assim.channel_rt.tm00.conus.nc", d, t))
		}
	}

	return urls
}

func rangeDate(start, end time.Time) func() time.Time {
	y, m, d := start.Date()
	start = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	y, m, d = end.Date()
	end = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)

	return func() time.Time {
		if start.After(end) {
			return time.Time{}
		}
		date := start
		start = start.AddDate(0, 0, 1)
		return date
	}
}

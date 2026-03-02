package config

type ResetTime struct {
	Hour     int
	Minute   int
	TimeZone string
}

const minutesInDay = 24 * 60

func (r *ResetTime) Add(minutes int) {
	minutes = max(0, minutes)

	if minutes == 0 {
		return
	}

	totalMinutes := (r.Hour * 60) + r.Minute + minutes

	totalMinutes = ((totalMinutes % minutesInDay) + minutesInDay) % minutesInDay

	r.Hour = totalMinutes / 60
	r.Minute = totalMinutes % 60
}

func ResetTimeByAccountType(t string) ResetTime {
	switch t {
	case "endfield", "genshin", "starrail", "honkai", "zzz", "themis":
		return ResetTime{
			Hour:     0,
			Minute:   0,
			TimeZone: "Asia/Shanghai",
		}
	default:
		// default just goes 00:00 at current timezone
		return ResetTime{}
	}
}

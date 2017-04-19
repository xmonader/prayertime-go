package prayertime

import (
	"fmt"
	m "math"
	"time"
)

const (
	radToDeg = 180 / m.Pi
	degToRad = m.Pi / 180
)

type coordinate struct {
	longitude float64
	latitude  float64
	zone      float64
}

type Prayertime struct {
	date       time.Time
	coordinate *coordinate
	Calendar   int
	Season     int
	Mazhab     int
	Fajr       float64
	Shrouk     float64
	Zuhr       float64
	Asr        float64
	Maghrib    float64
	Isha       float64
	dec        float64
}

func removeDuplication(val float64) float64 {
	return float64(int(val) % 360)
}
func (self *Prayertime) equation(alt float64) float64 {
	return radToDeg * (m.Acos((m.Sin(degToRad*(alt)) - m.Sin(degToRad*(self.dec))*m.Sin(degToRad*(self.coordinate.latitude))) / (m.Cos(degToRad*(self.dec)) * m.Cos(degToRad*(self.coordinate.latitude)))))

}

func (self *Prayertime) Calculate() {

	year := self.date.Year()
	month := int(self.date.Month())
	day := float64(self.date.Day())
	longitude := self.coordinate.longitude
	latitude := self.coordinate.latitude
	zone := self.coordinate.zone
	julianDay := -730531.5 + float64(367*year) - float64((year+(month+9)/12)*7/4) + float64(275*month/9) + day
	sunLength := 280.461 + 0.9856474*julianDay
	sunLength = removeDuplication(sunLength)
	middleSun := 357.528 + 0.9856003*julianDay
	middleSun = removeDuplication(middleSun)

	lamda := sunLength + 1.915*m.Sin(degToRad*(middleSun)) + 0.02*m.Sin(degToRad*(2*middleSun))
	lamda = removeDuplication(lamda)

	obliquity := 23.439 - 0.0000004*julianDay

	alpha := radToDeg * (m.Atan(m.Cos(degToRad*(obliquity)) * m.Tan(degToRad*(lamda))))

	if 90 < lamda && lamda < 180 {

		alpha += 180

	} else if 100 < lamda && lamda < 360 {
		alpha += 360
	}
	ST := 100.46 + 0.985647352*julianDay
	ST = removeDuplication(ST)

	self.dec = radToDeg * (m.Asin(m.Sin(degToRad*(obliquity)) * m.Sin(degToRad*(lamda))))

	noon := alpha - ST

	if noon < 0 {
		noon += 360

	}

	UTNoon := noon - longitude
	localNoon := (UTNoon / 15) + zone
	zuhr := localNoon                                // Zuhr Time.
	maghrib := localNoon + self.equation(-0.8333)/15 // Maghrib Time
	shrouk := localNoon - self.equation(-0.8333)/15  // Shrouk Time

	fajrAlt := 0.0
	ishaAlt := 0.0

	if self.Calendar == UmmAlQuraUniversity {
		fajrAlt = -19
	} else if self.Calendar == EgyptianGeneralAuthorityOfSurvey {
		fajrAlt = -19.5
		ishaAlt = -17.5
	} else if self.Calendar == MuslimWorldLeague {
		fajrAlt = -18
		ishaAlt = -17
	} else if self.Calendar == IslamicSocietyOfNorthAmerica {
		fajrAlt = -15
		ishaAlt = -15
	} else if self.Calendar == UnivOfIslamicSciencesKarachi {
		fajrAlt = -18
		ishaAlt = -18
	}
	fajr := localNoon - self.equation(fajrAlt)/15 // Fajr Time
	isha := localNoon + self.equation(ishaAlt)/15 // Isha Time

	if self.Calendar == UmmAlQuraUniversity {
		isha = maghrib + 1.5

	}

	asrAlt := 0.0

	if self.Mazhab == Hanafi {
		asrAlt = 90 - radToDeg*(m.Atan(2+m.Tan(degToRad*(m.Abs(latitude-self.dec)))))

	} else {
		asrAlt = 90 - radToDeg*(m.Atan(1+m.Tan(degToRad*(m.Abs(latitude-self.dec)))))
	}
	asr := localNoon + self.equation(asrAlt)/15 // Asr Time.

	// Add one hour to all times if the season is Summmer.
	if self.Season == Summer {
		fajr++
		shrouk++
		zuhr++
		asr++
		maghrib++
		isha++
	}
	self.Shrouk = shrouk
	self.Fajr = fajr
	self.Zuhr = zuhr
	self.Asr = asr
	self.Maghrib = maghrib
	self.Isha = isha

}

// ToHRTime: Convert a double value (e.g shrouk, fajr,... time calculated in Prayertime struct) to a human readable time
func ToHRTime(val float64, isAM bool) string {
	//"""val: double -> human readable string of format "%I:%M:%S %p" """

	var time string
	var zone string
	var hours int
	var minutes int
	var seconds int

	intval := int(val) // val is double.
	if isAM {
		if (intval%12) > 0 && intval%12 < 12 {
			zone = "AM"
		} else {
			zone = "PM"
		}
	} else {
		zone = "PM"
	}
	if intval > 12 {
		hours = intval % 12
	} else if intval%12 == 12 {
		hours = intval
	} else {
		hours = intval
	}
	val -= m.Floor(val)
	val *= 60
	minutes = int(val)

	val -= m.Floor(val)
	val *= 60
	seconds = int(val)

	time = fmt.Sprintf("%d:%d:%d %s", hours, minutes, seconds, zone)
	return time
}
func New(longitude, latitude, zone float64, year int, month time.Month, day int, calendar int, mazhab int, season int) *Prayertime {
	return &Prayertime{
		coordinate: &coordinate{
			latitude:  latitude,
			longitude: longitude,
			zone:      zone,
		},
		date:     time.Date(year, month, day, 0, 0, 0, 0, time.UTC),
		Calendar: calendar,
		Season:   season,
	}
}

// Show: prints the times for quick access.
func (self *Prayertime) Show() {
	fmt.Println(self.Fajr, self.Zuhr, self.Asr, self.Maghrib, self.Isha)
}

// SimpleReport reports all time to stdout.
func (self *Prayertime) SimpleReport() {
	fmt.Println(ToHRTime(self.Fajr, true))
	fmt.Println(ToHRTime(self.Shrouk, true))
	fmt.Println(ToHRTime(self.Zuhr, true))
	fmt.Println(ToHRTime(self.Asr, false))
	fmt.Println(ToHRTime(self.Maghrib, false))
	fmt.Println(ToHRTime(self.Isha, false))
}

// Convert hrtime (returned by ToHRTime) to a Go Datetime object.
func ToDateTime(hrtime string) (time.Time, error) {
	return time.Parse("%I:%M:%S %p", hrtime)
}

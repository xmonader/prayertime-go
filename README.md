# prayertime-go
Provides an implementation for calculation Muslims prayertimes with simple API


## Example
```golang

package main

import "github.com/xmonader/prayertime-go/prayertime"

func main() {
	pt := prayertime.New(31.2599, 30.0599, 2, 2010, 8, 6, prayertime.CalcEgyptianGeneralAuthorityOfSurvey, prayertime.MazhabDefault, true)
	pt.Calculate()
	pt.SimpleReport()
	pt.Show()
}



```
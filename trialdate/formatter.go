package trialdate

import (
	"fmt"
	"strings"
	"errors"
)


func Formatter(badDate string) func() (string, error) {
	months := make(map[string]string)
	months["January"] = "01"
	months["February"] = "02"
	months["March"] = "03"
	months["April"] = "04"
	months["May"] = "05"
	months["June"] = "06"
	months["July"] = "07"
	months["August"] = "08"
	months["September"] = "09"
	months["October"] = "10"
	months["November"] = "11"
	months["December"] = "12"

	days := make(map[string]string)
	days["1"] = "01"
	days["2"] = "02"
	days["3"] = "03"
	days["4"] = "04"
	days["5"] = "05"
	days["6"] = "06"
	days["7"] = "07"
	days["8"] = "08"
	days["9"] = "09"

	currentDate := badDate
	
	return func() (string, error) {
		array := strings.Split(currentDate, " ")
		
		if len(array) < 3 {
			return "", errors.New("invalid date")
		}

		month := months[array[0]]
		parsedDay := array[1][0:len(array[1]) - 1]
		day := days[parsedDay]
		year := array[2]

		if day == "" {
			day = parsedDay
		}

		return fmt.Sprintf("%s-%s-%s", year, month, day), nil
  }

}
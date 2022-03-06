//Quick script to parse some teams logs and return logs based on either a date or log level, this is because the teams logs are full of rubbish
//and I don't want to spend time cleaning them up. Added some flags so I can one-liner it and pipe it into other things. Basically what I'd do with awk and grep, but in Go so
//I can use it cross-platform - i.e. in Windows:  ./teamslogs.exe -l error -f logs.txt | out-file errors.txt
//Usage: teamslogs -flags -f <file>
//
//Flags:
//-d <date>
//Date to return logs for, format is 03/2022 or 03/03/2022
//
//-l <level>
//Level to return logs for, info, warning, error, event
//
//-f <file>
//File to parse - has to be either a full path or relative to the current directory. This also has to be in a readable format by the io package. If you can cat it, this can read it.
//
//-h
//Prints this help message
//

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

//Struct to hold the log level
type LogLevel struct {
	Info    string
	Warning string
	Error   string
	Event   string
}

//Struct to hold the log message
type LogMessage struct {
	Level   LogLevel
	Date    string
	Message string
}

/*Function to parse the logs based on date and level and return a slice of LogMessage structs. Each log message is a line in the log file, and begins with a date and a level
such as Sun Feb 27 2022 18:00:16 GMT+1100 (Australian Eastern Daylight Time) <15000> -- info -- This is a test message */
func ParseLogs(logFile string) []LogMessage {
	//Open the log file
	file, err := os.Open(logFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	//Create a slice of LogMessage structs
	var logs []LogMessage

	//Create a scanner to read the log file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		//Split the line into an array of strings, splitting at the `--` on each line
		line := strings.Split(scanner.Text(), ` -- `)
		var log LogMessage

		//If the line length is greater than 2
		//Append to the log var
		if len(line) > 2 {
			log.Date = line[0]
			log.Level.Info = line[1]
			log.Message = line[2]
		}

		//Append the LogMessage struct to the slice of LogMessage structs
		logs = append(logs, log)
	}

	//Return the slice of LogMessage structs
	return logs
}

//Return logs based on a specified date
func GetLogsDate(logMessages []LogMessage, date string) []LogMessage {

	cDate := ConvertDate(date)

	//Create a slice to hold the log messages
	var logs []LogMessage

	//Loop through the slice of log messages
	for _, logMessage := range logMessages {
		// fmt.Printf(cDate, logMessage.Date, "\n")
		//Use regexp to check if the date matches the date in the log message, and for each match, append the log message to the logs slice
		if regexp.MustCompile(cDate).MatchString(logMessage.Date) {
			logs = append(logs, logMessage)
		}
	}

	return logs
}

/*Parse date. Date is accepted in the below formats:
05/04/2022
04/2022
and will be converted to 'Mar 4 2022' format to compare against the log file

This is a very quick and dirty way - If I made this to be shared with others I'd probably convert each date to Unix time and use a range
or something.. I don't know, dates suck and switch cases make me sad. */
func ConvertDate(date string) string {

	//Split the date into an array of strings
	dateArray := strings.Split(date, "/")

	//Convert the date to the correct format
	//Day is in the first position of the array
	//Month is in the second position of the array
	//Year is in the third position of the array and should need no conversion as /22 is not accepted.
	var day string
	var month string
	//If the date is in the correct format
	if len(dateArray) == 3 {

		switch dateArray[1] {
		case "01":
			month = "Jan"
		case "02":
			month = "Feb"
		case "03":
			month = "Mar"
		case "04":
			month = "Apr"
		case "05":
			month = "May"
		case "06":
			month = "Jun"
		case "07":
			month = "Jul"
		case "08":
			month = "Aug"
		case "09":
			month = "Sep"
		case "10":
			month = "Oct"
		case "11":
			month = "Nov"
		case "12":
			month = "Dec"
		}
		//If the day is less than 10, add a 0 to the front of the day
		if len(dateArray[0]) == 1 {
			day = "0" + dateArray[0]
		} else {
			day = dateArray[0]
		}
		//From the date array, return a date that is in the format 'Mar 04 2022'
		return month + " " + day + " " + dateArray[2]
	} else if len(dateArray) == 2 {
		switch dateArray[0] {
		case "01":
			month = "Jan"
		case "02":
			month = "Feb"
		case "03":
			month = "Mar"
		case "04":
			month = "Apr"
		case "05":
			month = "May"
		case "06":
			month = "Jun"
		case "07":
			month = "Jul"
		case "08":
			month = "Aug"
		case "09":
			month = "Sep"
		case "10":
			month = "Oct"
		case "11":
			month = "Nov"
		case "12":
			month = "Dec"
		}
		//From the date array, return a date that is in the format 'Mar * 2022'
		return month + " " + ".*" + " " + dateArray[1]
	} else {
		//If the date is not in the correct format, return an empty string
		return ""
	}
}

//Return logs based on a specified level
func GetLogsByLevel(logMessages []LogMessage, level string) []LogMessage {

	//Create a slice to hold the log messages
	var logs []LogMessage

	//Loop through the slice of log messages
	for _, logMessage := range logMessages {
		//If the level matches the level entered by the user
		if logMessage.Level.Info == level || logMessage.Level.Warning == level || logMessage.Level.Error == level || logMessage.Level.Event == level {
			//Append the log message to the slice
			logs = append(logs, logMessage)

		}
	}

	return logs
}

//Usage:
//teamslogs -d <date> -l <level> -f <file>
func main() {
	//Create a flag for the date
	var date string
	//Create a flag for the level
	var level string
	//Create a flag for the file
	var file string

	//Create a flag for the date
	flag.StringVar(&date, "d", "", "Date to return logs for")
	//Create a flag for the level
	flag.StringVar(&level, "l", "", "Level to return logs for")
	//Create a flag for the file
	flag.StringVar(&file, "f", "", "File to parse")

	//Parse the flags
	flag.Parse()

	//If the file flag is empty
	if file == "" {
		//Ask the user for the file
		fmt.Println("Please enter the file to parse: ")
		fmt.Scanln(&file)
	}

	//Parse the logs
	logMessages := ParseLogs(file)

	//If the date and level flags are both not empty, return logs based on both date and level
	if date != "" && level != "" {
		//Return logs based on the date and level
		for _, logMessage := range GetLogsByLevel(GetLogsDate(logMessages, date), level) {
			fmt.Println(logMessage.Level, logMessage.Date, "--", logMessage.Message)
		}
	} else {
		if date == "" {
			//If the level flag is empty
			if level == "" {
				//Return all logs
				fmt.Println(logMessages)
			} else {
				//Return logs based on the level
				for _, logMessage := range GetLogsByLevel(logMessages, level) {
					fmt.Println(logMessage.Date, "--", logMessage.Message)
				}

			}
		} else {
			//If the level flag is empty
			if level == "" {
				//Return logs based on the date
				for _, logMessage := range GetLogsDate(logMessages, date) {
					fmt.Println(logMessage.Date, "--", logMessage.Message)
				}
			}
		}
	}
}

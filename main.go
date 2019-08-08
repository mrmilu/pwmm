package main

import (
	"encoding/json"
	"fmt"
	"github.com/akamensky/argparse"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// structures to parse yaml & json
type config struct {
	Credentials auth
	Events      []event
}

type auth struct {
	Apikey   string
	User     string
	Password string
}

type event struct {
	Name       string
	Startdate  string
	Finishdate string
	disabled   bool
	id         string
}

type pingdomCheck struct {
	Name   string
	Id     int
	Status string
}

// global vars
var checkInterval = 5000 // milliseconds for main loop
var httpClient = &http.Client{}
var credentials auth

// parse yaml, defines global credentials variable and return slice with
// maintenance events. I know, it's unorthodox and dirty.
func parseConfig(path *string) []event {
	source, err := ioutil.ReadFile(*path)
	if err != nil {
		panic(err)
	}
	var parsed config
	err = yaml.Unmarshal(source, &parsed)
	if err != nil {
		panic(err)
	}
	credentials = parsed.Credentials
	return parsed.Events
}

// methods of an event to get dates from the string
func (e *event) getStart() time.Time {
	return parseTime(e.Startdate)
}

func (e *event) getFinish() time.Time {
	return parseTime(e.Finishdate)
}

// convert string time to time datetype.
func parseTime(s string) time.Time {
	// ensure use always the system timezone
	currentZone, _ := time.Now().Zone()
	date, err := time.Parse("02-01-2006 15:04 MST", s+" "+currentZone)
	if err != nil {
		panic(err)
	}
	return date
}

// return (id, paused/unpaused) given a check name.
// this is another dirty part: asumes that the check always be active (could have more states)
func getCheck(name string) (string, bool) {
	// I dont like how do the request :-S
	request, err := http.NewRequest("GET", "https://api.pingdom.com/api/2.1/checks", nil)
	if err != nil {
		log.Fatalln(err)
	}
	request.Header.Set("App-Key", credentials.Apikey)
	request.SetBasicAuth(credentials.User, credentials.Password)
	resp, err := httpClient.Do(request)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	// Create temporal struct to parse response json correctly
	// this part should be outside of this function
	type structure struct {
		Checks []pingdomCheck
		Counts map[string]int
	}
	var s structure

	err = json.Unmarshal(body, &s)
	if err != nil {
		panic(err.Error())
	}

	var ID string
	var disabled bool
	for index := range s.Checks {
		if s.Checks[index].Name == name {
			ID = strconv.Itoa(s.Checks[index].Id)
			if s.Checks[index].Status == "paused" {
				disabled = true
			} else {
				disabled = false
			}
		}
	}
	return ID, disabled
}

// this name is confuse, just do the request to pause/unpause a specific check
func pingdomRequester(id string, action string) {
	var status string

	// to format url string later, very nasty.
	if action == "enable" {
		status = "false"
	} else if action == "disable" {
		status = "true"
	} else {
		panic("wrong param in requester")
	}

	// format the url params like a string, probably the library let do it in another better way
	url := fmt.Sprintf("https://api.pingdom.com/api/2.1/checks/%s?paused=%s", id, status)

	// setup the request, same code as getCheck, again, this should be outside the function because DRY
	request, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	request.Header.Set("App-Key", credentials.Apikey)
	request.SetBasicAuth(credentials.User, credentials.Password)
	resp, err := httpClient.Do(request)
	if err != nil {
		log.Fatalln(err)
	}
	text, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(text))
	if resp.StatusCode != 200 {
		// wonderful control
		log.Fatalln("Problem in request")
		if err != nil {
			panic(err.Error())
		}
		log.Fatalln(string(text))
	}
}

// iterate for events, get datetimes and determine if should be change their pause state
func check(events []event) {
	for e := range events {
		now := time.Now().UnixNano()
		start := events[e].getStart().UnixNano()
		finish := events[e].getFinish().UnixNano()
		if now < start && events[e].disabled {
			pingdomRequester(events[e].id, "enable")
			events[e].disabled = false
			fmt.Printf("%s - Service '%s' enabled", time.Now(), events[e].Name)
		}
		if now > start && now < finish && !events[e].disabled {
			pingdomRequester(events[e].id, "disable")
			events[e].disabled = true
			fmt.Printf("%s - Service '%s' disabled", time.Now(), events[e].Name)
		}
		if now > finish && events[e].disabled {
			pingdomRequester(events[e].id, "enable")
			events[e].disabled = false
			fmt.Printf("%s - Service '%s' enabled", time.Now(), events[e].Name)
		}
	}
}

// parse args, initiate events struct and start infinite loop to check events
func main() {
	parser := argparse.NewParser(os.Args[0], "Manage window maintenance in pingdom automatically")
	configFile := parser.String("f", "file", &argparse.Options{Required: true, Help: "yame file with definitions of maintenance events"})
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}
	events := parseConfig(configFile)
	for e := range events {
		id, disabled := getCheck(events[e].Name)
		events[e].id = id
		events[e].disabled = disabled
	}
	for {
		check(events)
		time.Sleep(time.Duration(checkInterval) * time.Millisecond)
	}
}

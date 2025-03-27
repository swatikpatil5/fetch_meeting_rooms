package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const calendarListURL = "https://www.googleapis.com/calendar/v3/users/me/calendarList"

type Calendar struct {
	ID      string `json:"id"`
	Summary string `json:"summary"`
}

type CalendarList struct {
	Items []Calendar `json:"items"`
}

func main() {
	// Get access token from environment variable
	token := os.Getenv("GOOGLE_ACCESS_TOKEN")
	if token == "" {
		log.Fatal("Missing access token. Set GOOGLE_ACCESS_TOKEN env variable.")
	}

	// Create request
	req, err := http.NewRequest("GET", calendarListURL, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	// Parse JSON response
	var calendarList CalendarList
	if err := json.Unmarshal(body, &calendarList); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	// Print all calendars
	fmt.Println("📅 Google Calendar List:")
	for _, calendar := range calendarList.Items {
		fmt.Printf("- ID: %s, Name: %s\n", calendar.ID, calendar.Summary)
	}

	// Print only meeting rooms
	fmt.Println("\n🏢 Meeting Rooms:")
	meetingRoomsFound := false
	for _, calendar := range calendarList.Items {
		if isMeetingRoom(calendar.ID) {
			fmt.Printf("- ID: %s, Name: %s\n", calendar.ID, calendar.Summary)
			meetingRoomsFound = true
		}
	}

	if !meetingRoomsFound {
		fmt.Println("No meeting rooms found.")
	}
}

// Function to check if a calendar is a meeting room
func isMeetingRoom(calendarID string) bool {
	return strings.Contains(calendarID, "@resource.calendar.google.com")
}

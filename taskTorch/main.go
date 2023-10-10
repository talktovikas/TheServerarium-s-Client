package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// same as server
type Job struct {
	ID        int    `json:"id"`
	Timestamp string `json:"ts"`
	IsDone    bool   `json:"isdone"`
}

func formatTimestamp(timestampStr string) (string, error) {
	// Parse the timestamp string as an integer (milliseconds)
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return "", err
	}
	// Convert milliseconds to a time.Time object
	t := time.Unix(0, timestamp*int64(time.Millisecond))// Need to see some more data here
	// Format the time in the desired format
	formattedTime := t.Format("15:04:05::02:01:2006")
	return formattedTime, nil
}

func completeJob(input string) error {
	// Define the directory path
	directoryPath := "/Users/vikasya/job/" //Just for my test
	// Create the directory if it doesn't exist
	if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
		if err := os.MkdirAll(directoryPath, os.ModePerm); err != nil {
			return err
		}
	}
	// Define the file path
	filePath := directoryPath + input + ".txt" //more more readable filename and making it unique
	// Create or open the file for writing (truncate it if it already exists)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	// Write the input string to the file
	_, err = file.WriteString(input)
	if err != nil {
		return err
	}
	fmt.Printf("File '%s' created and content written successfully.\n", filePath)
	return nil
}

// Execute function which will basically do the job that's all
func execute() http.HandlerFunc {
	fmt.Println("I am about to execute")
	return func(w http.ResponseWriter, r *http.Request) {
		var job Job
		json.NewDecoder(r.Body).Decode(&job)
		formatted, err := formatTimestamp(job.Timestamp)
		if err != nil {
			fmt.Println("Maybe this is not a timestamp", err)
			return
		}
		completeJob(formatted) // What if it fails?
		json.NewEncoder(w).Encode("true")
	}
}

func main() {
	//------------> Server will send a Request to execute a job
	//------------> This client will basically try to execute the job and respose either true or false.
	//------------> True means the job is executed and false means it didn't
	// There might be a clever solution like doing a rpc call, but since I don't know all those, so its
	// like a noob I will try this method.
	// Note to self: Have to ask for a clever solution.
	router := mux.NewRouter()
	router.HandleFunc("/doexecute", execute()).Methods("Post")
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":5299", nil))
}
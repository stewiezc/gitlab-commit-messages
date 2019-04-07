package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type commit struct {
	Message    string `json:"message"`
	Created_at string `json:"created_at"`
}

func main() {
	// discover flags
	debug := flag.Bool("d", false, "some debug output")
	startHash := flag.String("start", "", "starting hash")
	endHash := flag.String("end", "master", "end hash")
	projectId := flag.String("project", "", "Gitlab project id")
	flag.Parse()

	// get Environment variables
	apiKey := os.Getenv("GITLAB_API_KEY")

	// print some info if debug
	if *debug == true {
		fmt.Println("apiKey:", apiKey)
		fmt.Println("projectId:", *projectId)
		fmt.Println("startHash:", *startHash)
		fmt.Println("endHash:", *endHash)
	}

	// get the start date for start commit
	_, created_at, getCommitErr := getSingleCommit(*startHash, *projectId, apiKey)
	if getCommitErr != nil {
		log.Fatal(getCommitErr)
	}
	startDate := created_at
	fmt.Println(startDate)

	// get list of commits since the startDate
}

func getSingleCommit(hash string, projectId string, apiKey string) (string, string, error) {
	client := &http.Client{}
	uri := fmt.Sprintf("https://gitlab.com/api/v4/projects/%v/repository/commits/%v", projectId, hash)
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("PRIVATE-TOKEN", apiKey)
	resp, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	result := commit{}
	jsonErr := json.Unmarshal(body, &result)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return result.Message, result.Created_at, nil
}

func getListCommits() {
	return
}

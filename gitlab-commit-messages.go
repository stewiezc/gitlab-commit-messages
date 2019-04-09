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

	// get the end date for end commit
	_, endCreated_at, endGetCommitErr := getSingleCommit(*endHash, *projectId, apiKey)
	if endGetCommitErr != nil {
		log.Fatal(endGetCommitErr)
	}
	endDate := endCreated_at

	// get list of commits since the startDate

	returnCommits, getListErr := getListCommits(*projectId, apiKey, startDate, endDate)
	if getListErr != nil {
		log.Fatal(getListErr)
	}
	jsonCommits, _ := json.Marshal(returnCommits)
	s := string(jsonCommits)
	fmt.Println(s)
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

	var commit commit
	jsonErr := json.Unmarshal(body, &commit)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return commit.Message, commit.Created_at, nil
}

func getListCommits(projectId string, apiKey string, startDate string, endDate string) ([]commit, error) {
	client := &http.Client{}
	uri := fmt.Sprintf("https://gitlab.com/api/v4/projects/%v/repository/commits?since=%v&until=%v", projectId, startDate, endDate)
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

	var commits []commit
	jsonErr := json.Unmarshal([]byte(body), &commits)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	//for l := range commits {
	//	fmt.Printf("message: %vcreated_at:%v\n", commits[l].Message, commits[l].Created_at)
	//}

	return commits, nil
}

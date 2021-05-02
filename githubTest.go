package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const IssuesURL = ""

type IssuesSearchResult struct {
	TotalCount int `json:"total_count`
	Items      []*Issue
}

type Issue struct {
	Numer     int
	HTMLURL   string `json:"html_url"`
	Title     string
	State     string
	User      *User
	CreatedAt time.Time
	Body      string
}

type User struct {
	Login   string
	HTMLURL string `josn:"html_url"`
}

func SearchIssues(terms []string) (*IssuesSearchResult, error) {
	q := url.QueryEscape(strings.Join(terms, " "))
	resp, err := http.Get(IssuesURL + "?q" + q)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("search query failed: %s", resp.Status)
	}

	var result IssuesSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		return nil, err
	}

	resp.Body.Close()
	return &result, nil
}

func main12() {
	// result, err := SearchIssues([]string{"htt"})
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Printf("%d issues:\n", result.TotalCount)
	// for _, item := range result.Items {
	// 	fmt.Printf("#%-5d %9.9s %.55\n",
	// 		item.Numer, item.User.Login, item.Title)
	// }

	a, _, _ := returnTest()
	println(a)
}

func returnTest() (a, b, c int) {
	return 1, 2, 3
}

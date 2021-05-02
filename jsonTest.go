package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type Movie struct {
	Title  string
	Year   int  `json:"released"`
	Color  bool `json:"color,omitempty"`
	Actors []string
}

var movies = []Movie{
	{Title: "Casablance", Year: 1942, Color: false,
		Actors: []string{"Humphrey Bogart", "ingrid Bergman"}},
	{Title: "Cool Hand Luke", Year: 1967, Color: true,
		Actors: []string{"Paul Newman"}},
	{Title: "Bullitt", Year: 1968, Color: true,
		Actors: []string{"Steve McQueen", "Jacqueline Bisset"}},
}

func main11() {
	data, err := json.MarshalIndent(movies, "", "    ")
	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", data)
	}
	fmt.Printf("%s\n", data)

	var titles []struct{ Title string }
	println("titles address:%d", &titles)
	if err := json.Unmarshal(data, &titles); err != nil {
		log.Fatalf("JSON ummarshaling failed: %s", err)
	}
	fmt.Println(titles)

}

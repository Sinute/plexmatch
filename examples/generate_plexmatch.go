package main

import (
	"fmt"

	"github.com/sinute/plexmatch"
)

func main() {
	p := plexmatch.PlexMatch{
		Title:  "【推しの子】",
		TMDBID: 203737,
		Season: 1,
		Episodes: []plexmatch.Episode{
			{ID: 1, Name: "Oshi no Ko - sp.mkv", Special: true},
			{ID: 1, Name: "Oshi no Ko - 01.mkv", Special: false},
			{ID: 2, Name: "Oshi no Ko - 02.mkv", Special: false},
		},
	}
	fmt.Println(string(plexmatch.Marshal(p)))
	// title: 【推しの子】
	//
	// tmdbid: 203737
	//
	// season: 1
	//
	// ep: SP01: Oshi no Ko - sp.mkv
	// ep: 01: Oshi no Ko - 01.mkv
	// ep: 02: Oshi no Ko - 02.mkv
}

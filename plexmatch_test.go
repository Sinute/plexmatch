package plexmatch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshal(t *testing.T) {
	plexMatch, err := Unmarshal([]byte(`
		title: 【推しの子】
		tmdbid: 203737
		season: 1
		ep: SP01: Oshi no Ko - sp.mkv
		ep: 01: Oshi no Ko - 01.mkv
		ep: 02: Oshi no Ko - 02.mkv
	`))
	assert.Nil(t, err)
	assert.EqualValues(t, PlexMatch{
		Title:  "【推しの子】",
		TMDBID: 203737,
		Season: 1,
		Episodes: []Episode{
			{ID: 1, Name: "Oshi no Ko - sp.mkv", Special: true},
			{ID: 1, Name: "Oshi no Ko - 01.mkv", Special: false},
			{ID: 2, Name: "Oshi no Ko - 02.mkv", Special: false},
		},
	}, plexMatch)
}

func TestMarshal(t *testing.T) {
	got := Marshal(PlexMatch{
		Title:  "【推しの子】",
		TMDBID: 203737,
		Season: 1,
		Episodes: []Episode{
			{ID: 2, Name: "Oshi no Ko - 02.mkv", Special: false},
			{ID: 1, Name: "Oshi no Ko - 01.mkv", Special: false},
			{ID: 2, Name: "Oshi no Ko - sp02.mkv", Special: true},
			{ID: 1, Name: "Oshi no Ko - sp01.mkv", Special: true},
			{ID: 2, Name: "Oshi no Ko - 02.mkv", Special: false},
		},
	})
	want := []byte(`title: 【推しの子】

tmdbid: 203737

season: 1

ep: SP01: Oshi no Ko - sp01.mkv
ep: SP02: Oshi no Ko - sp02.mkv
ep: 01: Oshi no Ko - 01.mkv
ep: 02: Oshi no Ko - 02.mkv
`)
	assert.EqualValues(t, want, got)
}

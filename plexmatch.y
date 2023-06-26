%{
package plexmatch

import(
    "bytes"
    "sort"
)

type PlexMatch struct {
    Title string
    Show string
    TMDBID int
    TVDBID int
    IMDBID int
    Year int
    Season int
    Episodes []Episode
}

type Episode struct {
    ID int
    Name string
    Special bool
}

%}

%union{
    plexMatch PlexMatch
    episode Episode
    episodes []Episode
    strVal string
    intVal int
}

%type <plexMatch> plexMatch
%type <plexMatch> hints
%type <plexMatch> hint
%type <strVal> title
%type <strVal> show
%type <intVal> tmdbID
%type <intVal> tvdbID
%type <intVal> imdbID
%type <intVal> year
%type <intVal> season
%type <episode> episode
%type <episodes> episodes

%token	<strVal>	TitleHint
%token	<strVal>	ShowHint
%token	<strVal>	SeasonHint
%token	<strVal>	YearHint
%token	<strVal>	TMDBHint
%token	<strVal>	TVDBHint
%token	<strVal>	IMDBHint
%token	<strVal>	GUIDHint
%token	<strVal>	EpisodeHint
%token	<strVal>	Special
%token	<strVal>	Str
%token	<intVal>	Int
%token	<float64>	Float

%%

plexMatch: hints
    {
        $$ = $1
        yylex.(*lex).plexMatch = $$
    }

hints: hints hint
    {
        $$ = $1.merge($2)
    }
| hint
    {
        $$ = $1
    }

hint: title
    {
        $$.Title = $1
    }
| show
    {
        $$.Show = $1
    }
| tmdbID
    {
        $$.TMDBID = $1
    }
| tvdbID
    {
        $$.TVDBID = $1
    }
| imdbID
    {
        $$.IMDBID = $1
    }
| year
    {
        $$.Year = $1
    }
| season
    {
        $$.Season = $1
    }
| episodes
    {
        $$.Episodes = $1
    }

title: TitleHint ':' Str
    {
        $$ = $3
    }

show: ShowHint ':' Str
    {
        $$ = $3
    }

tmdbID: TMDBHint ':' Int
    {
        $$ = $3
    }

tvdbID: TVDBHint ':' Int
    {
        $$ = $3
    }

imdbID: IMDBHint ':' Int
    {
        $$ = $3
    }
| IMDBHint ':' 't' 't' Int
    {
        $$ = $5
    }

year: YearHint ':' Int
    {
        $$ = $3
    }

season: SeasonHint ':' Int
    {
        $$ = $3
    }

episodes: episodes episode
	{
	    $$ = append($1, $2)
	}
| episode
    {
        $$ = append($$, $1)
    }

episode: EpisodeHint ':' Int ':' Str
	{
	    $$ = Episode{ID: $3, Name: $5, Special: false}
	}
| EpisodeHint ':' Special Int ':' Str
    {
	    $$ = Episode{ID: $4, Name: $6, Special: true}
	}

%%

func Marshal(p PlexMatch) []byte {
	buf := bytes.Buffer{}
	if p.Title != "" {
		buf.WriteString(__yyfmt__.Sprintf("title: %s\n\n", p.Title))
	}
	if p.Show != "" {
		buf.WriteString(__yyfmt__.Sprintf("show: %s\n\n", p.Show))
	}
	if p.TMDBID != 0 {
		buf.WriteString(__yyfmt__.Sprintf("tmdbid: %d\n\n", p.TMDBID))
	}
	if p.TVDBID != 0 {
		buf.WriteString(__yyfmt__.Sprintf("tvdbid: %d\n\n", p.TVDBID))
	}
	if p.IMDBID != 0 {
		buf.WriteString(__yyfmt__.Sprintf("imdbid: %d\n\n", p.IMDBID))
	}
	if p.Year != 0 {
		buf.WriteString(__yyfmt__.Sprintf("year: %d\n\n", p.Year))
	}
	if p.Season != 0 {
		buf.WriteString(__yyfmt__.Sprintf("season: %d\n\n", p.Season))
	}
	sort.SliceStable(p.Episodes, func(i, j int) bool {
		if p.Episodes[i].Special {
			return true
		}
		if p.Episodes[j].Special {
			return false
		}
		return p.Episodes[i].ID < p.Episodes[j].ID
	})
	ep := Episode{}
	for _, v := range p.Episodes {
		if v.ID == ep.ID && v.Special == ep.Special {
			continue
		}
		ep = v
		sp := ""
		if v.Special {
			sp = "SP"
		}
		buf.WriteString(__yyfmt__.Sprintf("ep: %s%02d: %s\n", sp, v.ID, v.Name))
	}
	return buf.Bytes()
}

func Unmarshal(data []byte) (PlexMatch, error) {
	l := newLex(data)
	_ = yyParse(l)
	return l.plexMatch, l.err
}

func (p PlexMatch) merge(b PlexMatch) PlexMatch {
    a := p
    if a.Title == "" {
        a.Title = b.Title
    }
    if a.Show == "" {
        a.Show = b.Show
    }
    if a.TMDBID == 0 {
        a.TMDBID = b.TMDBID
    }
    if a.TVDBID == 0 {
        a.TVDBID = b.TVDBID
    }
    if a.IMDBID == 0 {
        a.IMDBID = b.IMDBID
    }
    if a.Year == 0 {
        a.Year = b.Year
    }
    if a.Season == 0 {
        a.Season = b.Season
    }
    a.Episodes = append(a.Episodes, b.Episodes...)
    return a
}

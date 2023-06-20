package plexmatch

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

const (
	statusStart   = iota
	statusIntHint // : int
	statusIntVal  // int$

	statusStrHint // : str
	statusStrVal  // str$

	statusEpisodeHint       // : [sp]int : str
	statusEpisodeVal        // [sp]int : str
	statusSpecialEpisodeVal // int : str

	hintKeyTitle    = "title"
	hintKeyShow     = "show"
	hintKeyYear     = "year"
	hintKeySeason   = "season"
	hintKeyTMDB     = "tmdbid"
	hintKeyTVDB     = "tvdbid"
	hintKeyIMDB     = "imdbid"
	hintKeyEpisode1 = "episode"
	hintKeyEpisode2 = "ep"
	hintKeySpecial  = "sp"
)

type lex struct {
	line      int
	input     []byte
	pos       int
	status    int
	plexMatch PlexMatch
	err       error
}

// Lex satisfies yyLexer.
func (l *lex) Lex(lval *yySymType) int {
	return l.scan(lval)
}

// Error satisfies yyLexer.
func (l *lex) Error(s string) {
	l.err = fmt.Errorf("%s, line: %d, pos: %d", s, l.line+1, l.pos)
}

func newLex(input []byte) *lex {
	return &lex{
		input: input,
	}
}

func (l *lex) scan(lval *yySymType) int {
	for b := l.lookNext(1); len(b) != 0; b = l.lookNext(1) {
		switch {
		case b[0] == '\n':
			l.line++
			l.next(1)
			continue
		case unicode.IsSpace(rune(b[0])):
			l.next(1)
			continue
		case l.status == statusStart:
			switch {
			case strings.ToLower(string(l.lookNext(len(hintKeyTitle)))) == hintKeyTitle:
				l.next(len(hintKeyTitle))
				l.status = statusStrHint
				return TitleHint
			case strings.ToLower(string(l.lookNext(len(hintKeyShow)))) == hintKeyShow:
				l.next(len(hintKeyShow))
				l.status = statusStrHint
				return ShowHint
			case strings.ToLower(string(l.lookNext(len(hintKeyYear)))) == hintKeyYear:
				l.next(len(hintKeyYear))
				l.status = statusIntHint
				return YearHint
			case strings.ToLower(string(l.lookNext(len(hintKeySeason)))) == hintKeySeason:
				l.next(len(hintKeySeason))
				l.status = statusIntHint
				return SeasonHint
			case strings.ToLower(string(l.lookNext(len(hintKeyTMDB)))) == hintKeyTMDB:
				l.next(len(hintKeyTMDB))
				l.status = statusIntHint
				return TMDBHint
			case strings.ToLower(string(l.lookNext(len(hintKeyTVDB)))) == hintKeyTVDB:
				l.next(len(hintKeyTVDB))
				l.status = statusIntHint
				return TVDBHint
			case strings.ToLower(string(l.lookNext(len(hintKeyIMDB)))) == hintKeyIMDB:
				l.next(len(hintKeyIMDB))
				l.status = statusIntHint
				return IMDBHint
			case strings.ToLower(string(l.lookNext(len(hintKeyEpisode1)))) == hintKeyEpisode1:
				l.next(len(hintKeyEpisode1))
				l.status = statusEpisodeHint
				return EpisodeHint
			case strings.ToLower(string(l.lookNext(len(hintKeyEpisode2)))) == hintKeyEpisode2:
				l.next(len(hintKeyEpisode2))
				l.status = statusEpisodeHint
				return EpisodeHint
			default:
				return int(b[0])
			}
		case l.status == statusStrHint:
			l.next(1)
			l.status = statusStrVal
			return int(b[0])
		case l.status == statusStrVal:
			lval.strVal = l.scanStringUntil('\n')
			l.status = statusStart
			return Str
		case l.status == statusIntHint:
			l.next(1)
			l.status = statusIntVal
			return int(b[0])
		case l.status == statusIntVal:
			i, ok := l.scanInt()
			if !ok {
				l.next(1)
				return int(b[0])
			}
			lval.intVal = i
			l.status = statusStart
			return Int
		case l.status == statusEpisodeHint:
			l.next(1)
			l.status = statusEpisodeVal
			return int(b[0])
		case l.status == statusEpisodeVal:
			if strings.ToLower(string(l.lookNext(len(hintKeySpecial)))) == hintKeySpecial {
				l.next(len(hintKeySpecial))
				l.status = statusSpecialEpisodeVal
				return Special
			}
			i, ok := l.scanInt()
			if !ok {
				l.next(1)
				return int(b[0])
			}
			lval.intVal = i
			l.status = statusStrHint
			return Int
		case l.status == statusSpecialEpisodeVal:
			i, ok := l.scanInt()
			if !ok {
				l.next(1)
				return int(b[0])
			}
			lval.intVal = i
			l.status = statusStrHint
			return Int
		default:
			return int(b[0])
		}
	}
	return 0
}

func (l *lex) scanInt() (int, bool) {
	buf := bytes.NewBuffer(nil)
	for {
		b := l.lookNext(1)
		if len(b) == 0 || !unicode.IsDigit(rune(b[0])) {
			break
		}
		l.next(1)
		buf.WriteByte(b[0])
	}
	if buf.Len() == 0 {
		return 0, false
	}
	val, err := strconv.ParseInt(buf.String(), 10, 64)
	if err != nil {
		return 0, false
	}
	return int(val), true
}

func (l *lex) scanStringUntil(brk byte) string {
	buf := bytes.NewBuffer(nil)
	for {
		b := l.lookNext(1)
		if len(b) == 0 || b[0] == brk {
			break
		}
		l.next(1)
		buf.WriteByte(b[0])
	}
	return buf.String()
}

func (l *lex) next(i int) {
	if l.pos >= len(l.input) || l.pos == -1 {
		l.pos = -1
	}
	l.pos += i
}

func (l *lex) lookNext(i int) []byte {
	if l.pos >= len(l.input) || l.pos == -1 {
		l.pos = -1
		return []byte{}
	}
	return l.input[l.pos : l.pos+i]
}

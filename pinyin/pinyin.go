package pinyin

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	STYLE_0 = 0 // zhong guo
	STYLE_1 = 1 // zhong1 guo2
	STYLE_2 = 2 // zhōng guó
)

type Word struct {
	code   []int32 // unicode
	pinyin string
}
type Words struct {
	w []Word
}

type Dict struct {
	spaceDelimiter bool
	keepNonChinese bool

	Punctuations       []int32
	PunctuationsOutput []string
	DB                 []Words
	Surnames           map[int32]*Words
}

func pinyinStyle1To2(finals map[string]string, s string) string {
	l := len(s)

	if l >= 4 {
		f := s[l-4:]
		if v, ok := finals[f]; ok {
			return strings.Replace(s, f, v, 1)
		}
	}

	if l >= 3 {
		f := s[l-3:]
		if v, ok := finals[f]; ok {
			return strings.Replace(s, f, v, 1)
		}
	}
	if l >= 2 {
		f := s[l-2:]
		if v, ok := finals[f]; ok {
			return strings.Replace(s, f, v, 1)
		}
	}
	return s
}

func (d *Dict) InitDict(style int) {
	var p []int32
	var pOutput []string

	for i := 0; i < len(PUNCTUATIONS); i += 2 {
		r, _ := utf8.DecodeRuneInString(PUNCTUATIONS[i])
		p = append(p, int32(r))
		pOutput = append(pOutput, PUNCTUATIONS[i+1])
	}
	d.Punctuations = p
	d.PunctuationsOutput = pOutput

	finals := make(map[string]string, 32)
	for i := 0; i < len(FINALS); i = i + 2 {
		finals[FINALS[i]] = FINALS[i+1]
	}
	style1to2 := func(s string) string {
		return pinyinStyle1To2(finals, s)
	}
	re := regexp.MustCompile(`\w+[1234]`)
	deleteTone := func(r rune) rune {
		switch r {
		case '1', '2', '3', '4':
			return -1
		}
		return r
	}
	newWord := func(w *Word, s, py string) {
		for len(s) > 0 {
			r, size := utf8.DecodeRuneInString(s)
			s = s[size:]
			w.code = append(w.code, int32(r))
		}

		if d.spaceDelimiter {
			py = strings.ReplaceAll(py, "\t", " ")
		} else {
			py = strings.ReplaceAll(py, "\t", "")
		}
		if style == STYLE_1 {
		} else if style == STYLE_0 {
			py = strings.Map(deleteTone, py)
		} else {
			py = re.ReplaceAllStringFunc(py, style1to2)
		}
		w.pinyin = py
	}

	ws := make([]Words, 40900, 40900)
	for i := 0; i < len(HAN_ZI); i = i + 2 {
		var w Word
		newWord(&w, HAN_ZI[i], HAN_ZI[i+1])
		code := w.code[0]
		ws[code].w = append(ws[code].w, w)
	}
	d.DB = ws

	sn := make(map[int32]*Words, 32)
	for i := 0; i < len(SURNAMES); i = i + 2 {
		var w Word
		newWord(&w, SURNAMES[i], SURNAMES[i+1])

		code := w.code[0]
		ws, ok := sn[code]
		if !ok {
			ws = &Words{}
			sn[code] = ws
		}
		ws.w = append(ws.w, w)
	}
	d.Surnames = sn
}

func NewDict(spaceDelimiter bool, keepNonChinese bool, style int) *Dict {
	d := &Dict{
		spaceDelimiter: spaceDelimiter,
		keepNonChinese: keepNonChinese,
	}

	d.InitDict(style)
	return d
}

func matchWord(w *Word, s []int32) ([]int32, bool) {
	if len(s) < len(w.code) {
		return s, false
	}
	for i := 0; i < len(w.code); i++ {
		if w.code[i] != s[i] {
			return s, false
		}
	}
	s = s[len(w.code):]
	return s, true
}

func (d *Dict) Translate(s string, surnamesMode bool) string {
	var pinyin strings.Builder
	var sr []int32

	for len(s) > 0 {
		r, size := utf8.DecodeRuneInString(s)
		s = s[size:]
		sr = append(sr, r)
	}

	if surnamesMode {
		code := sr[0]
		if ws, ok := d.Surnames[code]; ok {
			for i := 0; i < len(ws.w); i++ {
				ok := false
				sr, ok = matchWord(&ws.w[i], sr)
				if ok {
					pinyin.WriteString(ws.w[i].pinyin)
					break
				}
			}
		}
	}

	lastIsChinese := true
	for len(sr) > 0 {
		r := sr[0]

		if r >= int32(len(d.DB)) || len(d.DB[r].w) == 0 {
			for i := 0; i < len(d.Punctuations); i += 2 {
				if d.Punctuations[i] == r {
					pinyin.WriteString(d.PunctuationsOutput[i])
					goto next_char
				}
			}
			if d.keepNonChinese {
				if d.spaceDelimiter && lastIsChinese {
					pinyin.WriteString(" ")
				}
				pinyin.WriteRune(r)
			}

		next_char:
			lastIsChinese = false
			sr = sr[1:]
			continue
		}
		lastIsChinese = true

		ws := d.DB[r]
		for i := 0; i < len(ws.w); i++ {
			ok := false
			sr, ok = matchWord(&ws.w[i], sr)
			if ok {
				pinyin.WriteString(ws.w[i].pinyin)
				break
			}
		}
	}

	return strings.TrimSpace(pinyin.String())
}

func (d *Dict) Pinyin(s string) string {
	return d.Translate(s, false)
}

func (d *Dict) RenMing(s string) string {
	return d.Translate(s, true)
}

func IsChinese(str string) bool {
	for _, v := range str {
		if unicode.Is(unicode.Han, v) {
			return true
		}
	}
	return false
}

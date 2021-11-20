package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode"

	"eng-words/dict"

	"github.com/aaaton/golem/v4"
	"github.com/aaaton/golem/v4/dicts/en"
)

type config struct {
	dir        string
	file       *regexp.Regexp
	minWordLen int
	sortMode   sortMode
}

func defaultConfig() config {
	return config{
		dir:        "./sample",
		file:       regexp.MustCompile(`(?i)\.(txt|md)$`),
		minWordLen: 3,
		sortMode:   sortModeIndex,
		// sortMode: sortModeFrequency,
	}
}

type entry struct {
	word  string
	count int
	index int
}

type sortMode string

const (
	sortModeIndex     = sortMode("sortModeIndex")
	sortModeFrequency = sortMode("sortModeFrequency")
)

func main() {
	cfg := defaultConfig()

	lemmatizer, err := golem.New(en.New())
	if err != nil {
		panic(err)
	}

	dict, err := dict.FromEJDict()
	if err != nil {
		fmt.Print(err)
	}

	engWordMap := getEngWords(cfg, lemmatizer, dict)
	engWords := toEntry(&engWordMap)
	if cfg.sortMode == sortModeFrequency {
		sortInDescByCount(&engWords)
	} else if cfg.sortMode == sortModeIndex {
		sortByIndex(&engWords)
	}

	for _, v := range engWords {
		fmt.Printf("%s\t%s\n", v.word, dict.Get((v.word)))
	}
}

func toEntry(wordMap *map[string]entry) []entry {
	res := []entry{}
	for _, v := range *wordMap {
		res = append(res, v)
	}
	return res
}

func sortInDescByCount(words *[]entry) {
	sort.Slice(*words,
		func(i, j int) bool {
			return (*words)[i].count > (*words)[j].count

		})
}

func sortByIndex(words *[]entry) {
	sort.Slice(*words,
		func(i, j int) bool {
			return (*words)[i].index < (*words)[j].index

		})
}

func getEngWords(cfg config, lemmatizer *golem.Lemmatizer, dict *dict.Dict) map[string]entry {
	res := map[string]entry{}

	filepath.WalkDir(cfg.dir,
		func(path string, d fs.DirEntry, err error) error {
			if d.IsDir() {
				return err
			}

			if !cfg.file.MatchString(path) {
				return err
			}

			fp, err := os.Open(path)
			if err != nil {
				return err
			}

			index := 0
			lines := readLines(fp)
			for _, line := range lines {
				words := readWords(line)
				for _, word := range words {
					lemWord := lemmatizer.Lemma(word)
					if len(lemWord) < cfg.minWordLen {
						continue
					}
					if !dict.InDict(lemWord) {
						continue
					}
					e, ok := res[lemWord]
					if ok {
						e.count++
						res[lemWord] = e
					} else {
						newEntry := entry{lemWord, 0, index}
						res[lemWord] = newEntry
						index++
					}
				}
			}

			return err
		})

	return res
}

func readLines(f *os.File) []string {
	input := bufio.NewScanner(f)
	res := make([]string, 0, 2000)
	for input.Scan() {
		res = append(res, input.Text())
	}
	return res
}

func readWords(input string) []string {
	res := []string{}
	inputWithEnd := input + "."
	runes := make([]rune, 0, 10)
	for _, r := range inputWithEnd {
		if !isAlphabet(r) {
			if len(runes) > 0 {
				res = append(res, strings.ToLower(string(runes)))
			}
			runes = []rune{}
			continue
		}
		runes = append(runes, r)
	}
	return res
}

func isAlphabet(r rune) bool {
	if unicode.IsLower(r) {
		return true
	}
	if unicode.IsUpper(r) {
		return true
	}
	return false
}

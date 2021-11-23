package dict

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// Definitions is List Definition
type Definitions []string

// Dict is Dictionary
type Dict struct {
	m map[string]Definitions
}

// FromEJDict makes Dict from EJDict
func FromEJDict() (*Dict, error) {
	const path = "./dict/EJDict/src"
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, errors.Wrap(err, "read EJDict")
	}

	d := make(map[string]Definitions, 10000)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		f, err := os.Open(filepath.Join(path, file.Name()))
		if err != nil {
			return nil, errors.Wrap(err, "read EJDict file")
		}

		input := bufio.NewScanner(f)
		for input.Scan() {
			item := strings.Split(input.Text(), "\t")
			defs := d[item[0]]
			defs = append(defs, item[1])
			d[item[0]] = defs
		}

	}

	return &Dict{
			m: d,
		},
		nil
}

// Get returns definition
func (d *Dict) Get(word string) Definitions {
	return d.m[word]
}

// InDict check word is in dictionary
func (d *Dict) InDict(word string) bool {
	_, ok := d.m[word]
	return ok
}

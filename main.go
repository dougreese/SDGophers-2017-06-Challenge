package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/emirpasic/gods/trees/avltree"
)

const FILENAME = "./worldcitiespop.txt"

var records [][]string

type LineReader struct {
	buf *bufio.Reader
}

func NewLineReader(reader io.Reader) *LineReader {
	l := &LineReader{
		buf: bufio.NewReader(reader),
	}

	return l
}

func (l *LineReader) ReadLine(prompt string) (string, error) {
	if prompt != "" {
		fmt.Print(prompt, " ")
	}

	line, err := l.buf.ReadBytes('\n')
	if err != nil {
		return "", err
	}

	line = bytes.TrimRight(line, "\n")
	if len(line) > 0 {
		if line[len(line)-1] == 13 { //'\r'
			line = bytes.TrimRight(line, "\r")
		}
	}
	return string(line), err
}

func main() {
	f, err := os.Open(FILENAME)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fmt.Println("Reading data...")
	r := csv.NewReader(f)
	r.LazyQuotes = true
	records, err = r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	cols := records[0]
	log.Printf("%#v\n", cols)

	for i, v := range cols {
		fmt.Printf("%d: %s\n", i, v)
	}

	lr := NewLineReader(os.Stdin)
	colstr, err := lr.ReadLine("Search which column? (enter number)")
	if err != nil {
		log.Fatal(err)
	}
	col, _ := strconv.Atoi(colstr)
	fmt.Printf("Searching %d - %s\n", col, cols[col])

	fmt.Println("Indexing data...")
	tree := avltree.NewWithStringComparator()

	for i, val := range records[1:] {
		key := strings.ToLower(val[col])

		if value, found := tree.Get(key); found {
			ai := value.([]int)
			ai = append(ai, i+1)
			// fmt.Println("Found Adding:", key, len(ai))
			tree.Put(key, ai)
			continue
		}
		// fmt.Println("Adding:", key, i)
		tree.Put(key, []int{i + 1})

	}

	// fmt.Println(tree)

	for {
		start, _ := lr.ReadLine("Search start value?")
		if start == "" {
			fmt.Println("Come on, enter something")
			continue
		}
		end, _ := lr.ReadLine(fmt.Sprintf("Search end value (%s)?", start))
		if end == "" {
			end = start
		}
		// start := "erez"
		// end := "la"

		nodeStart, _ := tree.Floor(strings.ToLower(start))
		nodeEnd, _ := tree.Ceiling(strings.ToLower(end))
		// fmt.Printf("%+v\n", nodeStart)
		// fmt.Printf("%+v\n", nodeEnd)

		fmt.Println("======= Results =======")
		for i := nodeStart; i != nodeEnd.Next(); i = i.Next() {
			// fmt.Printf("%+v\n", i)
			vals := i.Value.([]int)
			for _, val := range vals {
				fmt.Printf("%s: %#v\n", i, records[val])
			}
		}
	}

}

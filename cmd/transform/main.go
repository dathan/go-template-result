package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"text/template"
)

// generic type of any array of dictionaries a flat (csv)
type GenericArrayMap []map[string]string

// the initial input is a csv of header == key to row value per key row

func main() {

	filePtr := flag.String("filename", "./tf.tmpl", "./filename")
	flag.Parse()

	strData := readStdIn()
	todos, err := convertCSVToTemplatevar(strData)
	if err != nil {
		panic(err)
	}

	baseName := path.Base(*filePtr)
	t := template.Must(template.New(baseName).ParseFiles("./" + *filePtr))
	err = t.Execute(os.Stdout, todos)
	if err != nil {
		panic(err)
	}

}

func readStdIn() string {

	info, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	if info.Mode()&os.ModeCharDevice != 0 || info.Size() <= 0 {
		fmt.Println("The command is intended to work with pipes.")
		fmt.Println("Usage: cat file | $0")
		return ""
	}

	reader := bufio.NewReader(os.Stdin)
	var ret string = ""
	for {
		input, _, err := reader.ReadLine()
		if err != nil && err == io.EOF {
			break
		}
		ret += fmt.Sprintf("%s\n", string(input))
	}
	return ret

}

func readCsv(input string) (GenericArrayMap, error) {
	r := csv.NewReader(strings.NewReader(input))
	//r.Comma = ','
	//r.Comment = '#'

	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	var todos GenericArrayMap = GenericArrayMap{}
	for c, data := range records {
		var lmap map[string]string = make(map[string]string)
		// TODO need an array of inputs and anb array of csv item that are strings
		lmap["Name"] = fmt.Sprintf("%03d", c)
		lmap["IP"] = data[0]
		todos = append(todos, lmap)
	}

	return todos, nil
}

func convertCSVToTemplatevar(input string) (GenericArrayMap, error) {

	return readCsv(input)

}

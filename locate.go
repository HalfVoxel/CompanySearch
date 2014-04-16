package main

import "fmt"
import "encoding/csv"
import "os"
import "strconv"
import "bufio"
import "log"
import "flag"

var seenBefore map[int]bool = map[int]bool{}

func Parse(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(file)
	r.FieldsPerRecord = 9
	r.TrailingComma = true
	r.LazyQuotes = true
	lines, err := r.ReadAll()

	file.Close()
	if err != nil {
		log.Fatalf("error in file %v, %v", path, err)
	}

	for _, line := range lines {
		index, _ := strconv.Atoi(line[2])

		if seenBefore[index] {
			continue
		} else {
			seenBefore[index] = true
		}

		fmt.Println(index)
	}
}

func Find(path string, needleIndex int) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(file)
	r.FieldsPerRecord = 9
	r.TrailingComma = true
	r.LazyQuotes = true
	lines, err := r.ReadAll()

	file.Close()
	if err != nil {
		log.Fatalf("error in file %v, %v", path, err)
	}

	for _, line := range lines {
		index, _ := strconv.Atoi(line[2])

		if index == needleIndex {
			fmt.Printf("\"%s\",\"%s\",\"%s\", \"%s\"\n", line[1], line[4], line[5], line[8])
		}
	}
}

func main() {
	const aleph = "ABCDEFGHIJKLMNOPQRSTUVXYZÅÄÖ"
	var generate = flag.Bool("generate", false, "generate cache file")

	flag.Parse()

	if *generate {

		for _, c := range aleph {
			for i := 0; i < 5000; i += 100 {

				path := fmt.Sprintf("output/%c_%d.csv", c, i)
				fmt.Println(path)
				Parse(path)
				fmt.Println("-1")
			}
		}

	} else {

		file, err := os.Open("lookupcache")
		if err != nil {
			log.Fatal(err)
		}

		lookup := map[int]string{}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			path := scanner.Text()
			for scanner.Scan() {
				s := scanner.Text()
				index, _ := strconv.Atoi(s)
				if index == -1 {
					break
				}
				lookup[index] = path
			}
		}

		scanner = bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			idx, _ := strconv.Atoi(scanner.Text())
			//fmt.Println(lookup[idx])
			path := lookup[idx]
			if len(path) > 0 {
				Find(lookup[idx], idx)
			}
		}
	}
}

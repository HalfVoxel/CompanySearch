package main

import "fmt"
import "os"
import "encoding/csv"
import "log"
import "bufio"
import "unicode/utf8"
import "strconv"
import "sort"
import "flag"
import "unicode"
import "runtime/pprof"
import _ "net/http/pprof"
import _ "net/http"
import "strings"

type Node struct {
	ch			rune
	lookup    []*Node
	links     []Info
	lastQuery int
	mins      [16]uint8
}

type Info struct {
	idx int
	//name string
}

var counter int = 0
var seenBefore map[int]bool = map[int]bool{}

func (node *Node) Insert(s string, index int, idx int ) {

	if index >= len(s) {
		info := Info{idx}
		node.links = append(node.links, info)

		if counter%10000 == 0 {
			fmt.Print("*")
		}

		counter++
		return
	}

	v, len := utf8.DecodeRuneInString(s[index:])
	v = unicode.ToLower(v)

	var next *Node = nil
	for _, other := range node.lookup {
		if other.ch == v {
			next = other
			break
		}
	}

	//next, ok := node.lookup[v]
	if next == nil {
		next = new(Node)
		next.ch = v
		//next.val = v
		next.lookup = make([]*Node,0)
		node.lookup = append(node.lookup,next)
		//fmt.Printf("Inserting %c from %s\n", next.ch, s)
	}

	//fmt.Print(string(v))
	next.Insert(s, index+len, idx)
}

type Result struct {
	result string
	error  int
	idx int
}

type ByError []Result

func (a ByError) Len() int           { return len(a) }
func (a ByError) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByError) Less(i, j int) bool { return a[i].error < a[j].error }

func (node *Node) Search(s string, index int, trace string, skipped int, queryIndex int, result *[]Result, best *int, addedWhitespace bool ) int {

	if node.lastQuery != queryIndex {
		node.lastQuery = queryIndex
		//node.mins = map[int]int{index: skipped}
	} else {
		prevErr := node.mins[index]
		//fmt.Printf("Been here before? %v %d < %d", ok2, prevErr, skipped)
		if int(prevErr) <= skipped+1 && prevErr != 0 {
			return -1
		}

		node.mins[index] = uint8(skipped+1)
	}

	if skipped > *best+6 {
		return -1
	}

	var v rune
	length := 0

	if index >= len(s) {
		for _, v := range node.links {
			//if len(node.links) > 0 {
			*result = append(*result, Result{trace, skipped, v.idx})

			if skipped < *best {
				*best = skipped
			}
			//}
			//fmt.Printf("Found %s with error %d\n", trace, skipped)
		}

		
		//return 0
	} else {

		
		v, length = utf8.DecodeRuneInString(s[index:])

		//fmt.Println(string(v))

		//next, ok := node.lookup[v]
		/*if ok {
			next.Search(s, index+length, trace+string(v), skipped+0, queryIndex, result, best)
		}*/
	}

	//fmt.Printf("Searching %s with error %d now at %c current node %c\n", trace, skipped, v, node.ch)

	if skipped < 16 {

		for _, next := range node.lookup {
			k := next.ch

			if k == v {
				next.Search(s, index+length, trace/*+string(k)*/, skipped+0, queryIndex, result, best, addedWhitespace)
			}
		}

		for _, next := range node.lookup {
			k := next.ch

			if k != v {
				if index == 0 {
					// Insert
					next.Search(s, index, trace+string(k), skipped+2, queryIndex, result, best, addedWhitespace)
				} else if index == len(s) {
					if k == ' ' {
						next.Search(s, index, trace+string(k), skipped+0, queryIndex, result, best, true)
					} else {
						if addedWhitespace {
							next.Search(s, index, trace+string(k), skipped+ 1, queryIndex, result, best, addedWhitespace)
						} else {
							next.Search(s, index, trace+string(k), skipped+ 2, queryIndex, result, best, addedWhitespace)
						}
					}
				} else {
					// Insert
					next.Search(s, index, trace+string(k), skipped+3, queryIndex, result, best, addedWhitespace)
				}
				// Replace
				next.Search(s, index+length, trace+string(k), skipped+4, queryIndex, result, best, addedWhitespace)
			}
		}

		// Skip
		node.Search(s, index+length, trace, skipped+2, queryIndex, result, best, addedWhitespace)
	}

	return -1
}

func parse(root *Node, path string) {
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

		root.Insert(line[1], 0, index)
		/*for i, _ := range line[1] {
			fmt.Printf (line[1][i])//"%v",runeValue)
		}*/
	}
}

func main() {
	const aleph = "ABCDEFGHIJKLMNOPQRSTUVXYZÅÄÖ"

	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	/*go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()*/

	node := new(Node)
	//node.lookup = map[rune]*Node{}
	node.lookup = make([]*Node,0)

	for _, c := range aleph {
		for i := 0; i < 5000; i += 100 {
			path := fmt.Sprintf("output/%c_%d.csv", c, i)
			parse(node, path)
		}
	}

	//fmt.Println("")

	bio := bufio.NewReader(os.Stdin)
	cnt := 0
	for {
		cnt++

		needle, _, err := bio.ReadLine()
		sneedle := strings.ToLower(string(needle))
		if err != nil {
			return
			//panic(err)
		}

		if sneedle == "x" {
			return
		}

		//fmt.Printf("Searching for %s\n",needle)
		var res []Result
		best := 1000
		node.Search(sneedle, 0, "", 0, cnt, &res, &best, false)

		sort.Sort(ByError(res))
		fmt.Println(len(res))
		for i, v := range res {
			if i > 10 {
				break
			}
			fmt.Printf("%d\n", v.idx)
			fmt.Printf("Found Result: %d: %s\n", v.error, v.result)
		}
		//fmt.Println("Done")
	}
}

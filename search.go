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
import "runtime/pprof"
import _ "net/http/pprof"
import "net/http"

type Node struct {
	lookup map[rune]*Node
	links []Info
	lastQuery int
	mins map[int]int
}

type Info struct {
	//name string
}

var counter int = 0
var seenBefore map[int]bool = map[int]bool{}

func (node *Node) Insert ( s string, index int ) {

	if index >= len(s) {
		info := Info{}
		node.links = append(node.links, info)

		if counter % 10000 == 0 {
			fmt.Print("*")
		}

		counter++
		return
	}

	v, len := utf8.DecodeRuneInString(s[index:])

	next, ok := node.lookup[v]
	if !ok {
		next = new(Node)
		//next.val = v
		next.lookup = map[rune]*Node {}
		node.lookup[v] = next
	}

	//fmt.Print(string(v))
	next.Insert (s,index+len)	
}

type Result struct {
	result string
	error int
}

type ByError []Result

func (a ByError) Len() int           { return len(a) }
func (a ByError) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByError) Less(i, j int) bool { return a[i].error < a[j].error }

func (node *Node) Search ( s string, index int, trace string, skipped int, queryIndex int, result *[]Result ) (int) {
	
	if node.lastQuery != queryIndex {
		node.lastQuery = queryIndex
		node.mins = map[int]int {index: skipped}
	} else {
		prevErr, ok2 :=  node.mins[index]
		//fmt.Printf("Been here before? %v %d < %d", ok2, prevErr, skipped)
		if ok2 && prevErr <= skipped {
			return -1
		}

		node.mins[index] = skipped
	}

	var v rune
	length := 0

	if index >= len(s) {
		/*for i, _ := range node.links {
			_ = i
			fmt.Printf("Found %s with error %d #%d\n", trace, skipped,i)
		}*/

		if len(node.links) > 0 {
			*result = append(*result, Result{trace,skipped})
		}

		//return 0
	} else {

		v, length = utf8.DecodeRuneInString(s[index:])

		//fmt.Print(string(v))

		next, ok := node.lookup[v]
		if ok {
			next.Search (s, index+length, trace +string(v), skipped+0, queryIndex, result)
		}
	}

	if skipped < 16 {
		// Skip
		node.Search (s, index+length, trace, skipped+2, queryIndex, result)

		for k, other := range node.lookup {
			if k != v {
				// Insert
				other.Search (s, index, trace +string(k), skipped+1, queryIndex, result)
				// Replace
				other.Search (s, index+length, trace +string(k), skipped+3, queryIndex, result)
			}
		}
	}

	return -1
}

func parse (root *Node, path string ) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal (err)
	}

	r := csv.NewReader (file)
	r.FieldsPerRecord = 9
	r.TrailingComma = true
	r.LazyQuotes = true
	lines, err := r.ReadAll ()

	file.Close ()
	if err != nil {
		log.Fatalf ("error in file %v, %v", path, err)
	}

	for _, line := range lines {
		index, _ := strconv.Atoi (line[2])
		
		if seenBefore[index] {
			continue
		} else {
			seenBefore[index] = true
		}

		root.Insert (line[1], 0)
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

    go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	node := new(Node)
	node.lookup = map[rune]*Node {}

	for _,c := range aleph {
		for i := 0; i < 5000; i += 100 {
			path := fmt.Sprintf("output/%c_%d.csv",c,i)
			parse (node, path)
		}
	}

	fmt.Println("")

	bio := bufio.NewReader(os.Stdin)
	cnt := 0
	for {
		cnt++

		needle, _, err := bio.ReadLine()

		if err != nil {
			panic(err)
		}

		if string(needle) == "x" {
			return
		}

		fmt.Printf("Searching for %s\n",needle)
		var res []Result
		node.Search(string(needle),0,"",0,cnt, &res)

		sort.Sort(ByError(res))

		for i, v := range res {
			if i > 10 {
				break
			}
			fmt.Printf("Found Result: %d: %s\n", v.error, v.result)
		}
		fmt.Println("Done")
	}
}
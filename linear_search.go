package main

import "fmt"
import "os"
import "encoding/csv"
import "log"
//import "encoding/json"
import "bufio"
//import "unicode/utf8"
import "strconv"
import "sort"
import "flag"
//import "unicode"
import "runtime/pprof"
import _ "net/http/pprof"
import _ "net/http"
import "strings"
import "./locate"
//import "time"
import "io"

var counter int = 0
var seenBefore map[int]bool = map[int]bool{}
var seenBeforeName map[string]bool = map[string]bool{}

type result struct {
	line []string
	error int
}

type ByError []result

func (a ByError) Len() int           { return len(a) }
func (a ByError) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByError) Less(i, j int) bool { return a[i].error < a[j].error }

func parse(path string, needles []string, results *[][]result) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(file)
	r.FieldsPerRecord = 9
	r.TrailingComma = true
	r.LazyQuotes = true

	if err != nil {
		log.Fatalf("error in file %v, %v", path, err)
	}

	for {
		line, err := r.Read ()

		if err == io.EOF {
			break;
		}

		if err != nil {
			panic (err);
		}

		index, _ := strconv.Atoi(line[2])

		if seenBefore[index] {
			continue
		} else {
			seenBefore[index] = true
		}

		name := strings.ToLower(line[1]);

		if seenBeforeName[name] {
			continue
		} else {
			seenBeforeName[name] = true
		}
		
		//root.Insert(line[1], 0, index)

		// 36,"A & A Ekholm Arkitektkontor AB","2837253","046-30 74 08","Lund","","/Lund/A__A_Ekholm_Arkitektkontor_AB/2837253","","http://maps.google.se/maps?f=q&source=s_q&hl=sv&geocode=&q=Sp%c3%a5rsn%c3%b6gatan+37%2c+22652%2c+LUND%2c+Sweden"

		
		url := line[5];
		//fmt.Printf ("%d : %s\n", string_distance(name,needles[0]), name );//"%v",runeValue)
		
		for needleIdx, needle := range(needles) {
			error := string_distance(name,needle)

			// URLs give bonus points
			if len(url) > 0 {
				error -= 1;
			}

			res := result{line:line, error: error};

			if len((*results)[needleIdx]) < 10 {
				(*results)[needleIdx] = append((*results)[needleIdx], res);
			} else {
				worstId := 0;
				worstErr := 0;
				for i, v := range (*results)[needleIdx] {
					if v.error > worstErr {
						worstErr = v.error;
						worstId = i;
					}
				}

				if error < worstErr {
					(*results)[needleIdx][worstId] = res;
				}
			}
			//for _, _ = range line {
				
			//}
		}
	}

	file.Close()
}

var distArr [][]int = nil;

func string_distance ( a, b string ) int {

	var cost int
	if distArr == nil || len(distArr) < len(a)+1 || len(distArr[0]) < len(b)+1 {
		distArr = make([][]int, len(a)+1)
		for i := 0; i < len(distArr); i++ {
			distArr[i] = make([]int, len(b)+1)
		}
	}
 
	for i := 0; i <= len(a); i++ {
		distArr[i][0] = i*2;
	}
 
	for i := 0; i <= len(b); i++ {
		distArr[0][i] = i*5;
	}
	
	for i := 1; i <= len(a); i++ {
		for j := 1; j <= len(b); j++ {
			if a[i-1] == b[j-1] {
				cost = 0
			} else {
				cost = 4 // Replace
			}
 			
			min1 := distArr[i-1][j] + 2 // Insert
			/*if a[i-1] == ' ' {
				min1 = 0;
			}*/

			min2 := distArr[i][j-1] + 5 // Erase
			min3 := distArr[i-1][j-1] + cost
			if min2 < min1 {
				min1 = min2;
			}
			if min3 < min1 {
				min1 = min3;
			}
			distArr[i][j] = min1;

			//distArr[i][j] = int(math.Min(math.Min(float64(min1), float64(min2)), float64(min3)))
		}
	}
 
	return distArr[len(a)][len(b)]
}

type SearchResult struct {
	Results []locate.LocatorResult
}

func read ( reader *bufio.Reader, ch chan []byte ) {
	for true {
		s, _, err := reader.ReadLine()
		if err != nil {
			panic(err)
		}

		ch <- s
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

	bio := bufio.NewReader(os.Stdin)
	
	needles := make([]string, 0)

	for {
		s, ioerr := bio.ReadString('\n');
		if ioerr != nil && ioerr != io.EOF {
			panic(ioerr)
		}
		//fmt.Printf("Read: %s\n",s)

		if len(s) > 0 {
			needles = append(needles,strings.ToLower(s))
		}

		if ioerr != nil {
			break;
		}
	}

	results := make([][]result, len(needles));
	for i := range needles {
		results[i] = make([]result, 0);
	}

	for _, c := range aleph {
		for i := 0; i < 5000; i += 100 {
			path := fmt.Sprintf("output/%c_%d.csv", c, i)
			parse(path, needles, &results)
		}
	}

	//fmt.Println("\nResults:\n");

		
	fmt.Println("[");
	for _, arr := range(results) {
		sort.Sort(ByError(arr));
		//fmt.Printf("%d : %s\n", v.error, v.name)
		//fmt.Println(v.line)

		fmt.Println("[");
		for _, v := range arr {
			// 36,"A & A Ekholm Arkitektkontor AB","2837253","046-30 74 08","Lund","","/Lund/A__A_Ekholm_Arkitektkontor_AB/2837253","","http://maps.google.se/maps?f=q&source=s_q&hl=sv&geocode=&q=Sp%c3%a5rsn%c3%b6gatan+37%2c+22652%2c+LUND%2c+Sweden"
			// 62,"A & Be Hotell & Vandrarhem AB","1752119","08-660 21 00","Stockholm","http://www.abehotel.com","/Stockholm/A__Be_Hotell__Vandrarhem_AB/1752119","mailto:info@abehotel.com","http://maps.google.se/maps?f=q&source=s_q&hl=sv&geocode=&q=Surbrunnsgatan+57%2c+11327%2c+STOCKHOLM%2c+Sweden"
			fmt.Printf("{ error: %d, name: '%s', city: '%s', url: '%s' }, \n", v.error, v.line[1], v.line[4], v.line[5]);
		}
		fmt.Println("],");
	}
	fmt.Println("]");
}

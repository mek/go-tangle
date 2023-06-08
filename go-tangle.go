package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type OpenFileParams struct {
	Fname string
	Mode  int
	Perm  os.FileMode
}

func Err(args ...interface{}) {
	for _, arg := range args {
		fmt.Fprint(os.Stderr, arg)
	}
	fmt.Fprintln(os.Stderr)
}

func WithOpenFile(params OpenFileParams, block func(*os.File)) {
	if params.Fname == "" {
		Err("no file name given to WithOpenFile, exiting")
		os.Exit(1)
	}
	if params.Mode == 0 {
		params.Mode = os.O_RDONLY
	}
	if params.Perm == 0 {
		params.Perm = 0400
	}
	fp, err := os.OpenFile(params.Fname, params.Mode, params.Perm)
	if err != nil {
		Err("Could not open file ", params.Fname, " error: ", err)
	}
	defer fp.Close()

	block(fp)
}

func AddArrayValue(arr map[string]string, key string, value string) {
	if _, exists := arr[key]; exists {
		arr[key] += "\n" + value
	} else {
		arr[key] = value
	}
}

func ExpandChunks(arr map[string]string, chunk string, indent string) {
	if _, exists := arr[chunk]; !exists {
		Err("could not find chunk ", chunk)
		Err("Available Chunks are:")
		for c := range arr {
			Err(c)
		}
		os.Exit(1)
	}

	lines := strings.Split(arr[chunk], "\n")
	for _, line := range lines {
		if re := regexp.MustCompile(`^(\s*)(<<)(.*)(>>)\s*$`); re.MatchString(line) {
			submatches := re.FindStringSubmatch(line)
			newIndent := submatches[1]
			newChunk := submatches[3]
			ExpandChunks(arr, newChunk, newIndent)
		} else {
			fmt.Print(indent)
			fmt.Println(line)
		}
	}
}

func main() {
	var requestedChunk string
	chunks := make(map[string]string)

	flag.Usage = func() {
		Err(os.Args[0], " -R <chuck to extract> filename")
	}
	flag.StringVar(&requestedChunk, "R", "*", "Expect Chunk")
	flag.Parse()

	if len(flag.Args()) != 1 {
		Err("only allowed one filename, you have %d\n", len(flag.Args()))
		flag.Usage()
		os.Exit(1)
	}
	fileName := flag.Arg(0)

	curFile := OpenFileParams{Fname: fileName}

	lineno := 0
	WithOpenFile(curFile, func(fp *os.File) {
		scanner := bufio.NewScanner(fp)
		inChunk := false
		var chunk string

		for scanner.Scan() {
			line := scanner.Text()
			lineno++
			if !inChunk {
				if re := regexp.MustCompile(`(^<<)(.*)(>>=$)`); re.MatchString(line) {
					submatches := re.FindStringSubmatch(line)
					chunk = submatches[2]
					inChunk = true
					AddArrayValue(chunks, chunk, fmt.Sprintf("# %s lineno %d", fileName, lineno))
					continue
				}
			} else {
				if re := regexp.MustCompile(`^(@).*(% def)?$`); re.MatchString(line) {
					inChunk = false
					chunk = ""
					continue
				}
				AddArrayValue(chunks, chunk, line)
			}
		}
		if err := scanner.Err(); err != nil {
			Err(err)
		}
	})

	ExpandChunks(chunks, requestedChunk, "")

}

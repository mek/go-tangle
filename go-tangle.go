// A Go version of mek's tangle, based on noweb
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// A structure for open files
type OpenFileParams struct {
	Fname string      // file name to open
	Mode  int         // os.O_READONLY, etc. (read, write, append, etc.)
	Perm  os.FileMode // The file mode for the file (0400, 0644, etc.)
}

// output to standard error
// Will go though each arg and write, without a new line, to standard error.
// Will write an new line after the args are written.
func Err(args ...interface{}) {
	for _, arg := range args {
		fmt.Fprint(os.Stderr, arg)
	}
	fmt.Fprintln(os.Stderr)
}

// Opens a file and runs a function.
// Will open the file (params) given by OpenFileParams and if
// successful will use the open file pointer (fp) to call the
// function (block).
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
		os.Exit(15)
	}
	defer fp.Close()
	fileInfo, err := fp.Stat()
	if err != nil {
		Err("Error getting file info ", err)
		os.Exit(15)
	}
	if fileInfo.IsDir() {
		Err(params.Fname, " is a file directory, exiting")
		os.Exit(15)
	}
	block(fp)

}

// Add string value to a map of strings
// If there is already a value in the string, append a newline
// to the end of the current string and then add the value.
func AddArrayValue(arr map[string]string, key string, value string) {
	if _, exists := arr[key]; exists {
		arr[key] += "\n" + value
	} else {
		arr[key] = value
	}
}

// Check if the chunk exists, if it does not err, exit with error.
// If it does exists, print out the value for the chuck, recusily
// checking to see if other chunks are in the value.
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

// main function
func main() {
	var requestedChunk string
	chunks := make(map[string]string)

	flag.Usage = func() {
		Err(os.Args[0], " -R <chuck to extract> filename")
	}
	flag.StringVar(&requestedChunk, "R", "*", "Expect Chunk")
	flag.Parse()

	if len(flag.Args()) != 1 {
		Err(fmt.Sprintf("only allowed one filename, you have %d", len(flag.Args())))
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

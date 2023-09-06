A Go version of mek's tangle, based on noweb

FUNCTIONS

func AddArrayValue(arr map[string]string, key string, value string)
    Add string value to a map of strings If there is already a value in the
    string, append a newline to the end of the current string and then add the
    value.

func Err(args ...interface{})
    output to standard error Will go though each arg and write, without a new
    line, to standard error. Will write an new line after the args are written.

func ExpandChunks(arr map[string]string, chunk string, indent string)
    Check if the chunk exists, if it does not err, exit with error. If it does
    exists, print out the value for the chuck, recusily checking to see if other
    chunks are in the value.

func WithOpenFile(params OpenFileParams, block func(*os.File))
    Opens a file and runs a function. Will open the file (params) given by
    OpenFileParams and if successful will use the open file pointer (fp) to call
    the function (block).


TYPES

type OpenFileParams struct {
	Fname string      // file name to open
	Mode  int         // os.O_READONLY, etc. (read, write, append, etc.)
	Perm  os.FileMode // The file mode for the file (0400, 0644, etc.)
}
    A structure for open files


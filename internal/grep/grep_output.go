package grep

// TODO: should any of these fields be pointers?
type GrepOutput struct {
	output   string // output data as an array of bytes
	filename string
	numLines int
}

// Return the grepOutput object serialized
func SerializeGrepOutput() {

}

func DeserializeGrepOutput() {

}

package transfer

type Metadata struct {
	Conn  string
	Fname string
	Host  string
	Port  string
	Size  int
}

type FileMetadataHeader struct {
	FnameSize uint64
	Fname     string
	Fsize     uint64
}

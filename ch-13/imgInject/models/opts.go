package models

//CmdLineOpts represents the cli arguments
type CmdLineOpts struct {
	Input    string
	Output   string
	Meta     bool
	Suppress bool
	Offset   string
	Inject   bool
	Payload  string
	Type     string
	Encode   bool
	Decode   bool
	AESencode bool
	AESdecode bool
	Key      string
	Specific string
	Jpeg	 bool
	Png      bool
	MultiInject string
}

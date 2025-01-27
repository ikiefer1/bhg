package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/blackhat-go/bhg/ch-13/imgInject/models"
	"github.com/blackhat-go/bhg/ch-13/imgInject/pnglib"
	"github.com/blackhat-go/bhg/ch-13/imgInject/utils"
	"github.com/spf13/pflag"
)

var (
	flags = pflag.FlagSet{SortFlags: false}
	opts  models.CmdLineOpts
	png   pnglib.MetaChunk
)

func init() {
	flags.StringVarP(&opts.Input, "input", "i", "", "Path to the original image file")
	flags.StringVarP(&opts.Output, "output", "o", "", "Path to output the new image file")
	flags.BoolVarP(&opts.Meta, "meta", "m", false, "Display the actual image meta details")
	flags.BoolVarP(&opts.Suppress, "suppress", "s", false, "Suppress the chunk hex data (can be large)")
	flags.StringVar(&opts.Offset, "offset", "", "The offset location to initiate data injection")
	flags.BoolVar(&opts.Inject, "inject", false, "Enable this to inject data at the offset location specified")
	flags.StringVar(&opts.Payload, "payload", "", "Payload is data that will be read as a byte stream")
	flags.StringVar(&opts.Type, "type", "rNDm", "Type is the name of the Chunk header to inject")
	flags.StringVar(&opts.Key, "key", "", "The enryption key for payload")
	flags.BoolVar(&opts.Encode, "encode", false, "XOR encode the payload")
	flags.BoolVar(&opts.Decode, "decode", false, "XOR decode the payload")
	flags.BoolVar(&opts.AESencode, "AESencode", false, "AES encode the payload")
	flags.BoolVar(&opts.AESdecode, "AESdecode", false, "AES decode the payload")
	flags.StringVar(&opts.Specific, "specific", "", "Enable this to edit existing gAMA")
	flags.BoolVar(&opts.Jpeg, "jpeg",false, "Enable if jpeg")
	flags.BoolVar(&opts.Png, "png",false, "Enable if png")
	flags.StringVar(&opts.MultiInject, "multi_inject", "", "This is the ancillary chunks that you would like to inject into")
	flags.Lookup("type").NoOptDefVal = "rNDm"
	flags.Usage = usage
	flags.Parse(os.Args[1:])

	if flags.NFlag() == 0 {
		flags.PrintDefaults()
		os.Exit(1)
	}
	if opts.Input == "" {
		log.Fatal("Fatal: The --input flag is required")
	}
	if opts.Offset != "" {
		byteOffset, _ := strconv.ParseInt(opts.Offset, 0, 64)
		opts.Offset = strconv.FormatInt(byteOffset, 10)
	}
	if opts.Suppress && (opts.Meta == false) {
		log.Fatal("Fatal: The --meta flag is required when using --suppress")
	}
	if opts.Meta && (opts.Offset != "") {
		log.Fatal("Fatal: The --meta flag is mutually exclusive with --offset")
	}
	if opts.Inject && (opts.Offset == "") {
		log.Fatal("Fatal: The --offset flag is required when using --inject")
	}
	if opts.Inject && (opts.Payload == "") {
		log.Fatal("Fatal: The --payload flag is required when using --inject")
	}
	if opts.Inject && opts.Key == "" {
		fmt.Println("Warning: No key provided. Payload will not be encrypted")
	}
	if opts.Encode && opts.Key == "" {
		log.Fatal("Fatal: The --encode flag requires a --key value")
	}
	if opts.Specific != "" && opts.Payload ==""{
		log.Fatal("Fatal: The --specific flag requires a --payload value")
	}
	if !opts.Jpeg && !opts.Png {
		log.Fatal("Fatal: The --jpg or --png file type must be designated")
	}
	if opts.MultiInject !="" && opts.Payload ==""{
		log.Fatal("Fatal: The --multi-inject flag requires a --payload value")
	}
	// if opts.MultiInject !="" && opts.Offset ==""{
	// 	log.Fatal("Fatal: The --multi-inject flag requires a --offset value")
	// }
}

func usage() {
	fmt.Fprintf(os.Stderr, "Example Usage: %s --png -i in.png -o out.png --inject --offset 0x85258 --payload 1234\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Example Encode Usage: %s -i in.png -o encode.png --inject --offset 0x85258 --payload 1234 --encode --key secret\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Example Decode Usage: %s -i encode.png -o decode.png --offset 0x85258 --decode --key secret\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Example AESencode Usage: --jpeg  -i encode.jpg -o workingTest.jpg --inject --offset 0x85258 --payload secret --AESencode --key dbooty\n")
	fmt.Fprintf(os.Stderr, "Example AESdecode Usage: --jpeg  -i workingTest.jpg -o workingTestOutput.jpg --offset 0x85258 --AESdecode --key dbooty\n")
	//example multi_inject --png -i images/gamaog.png -o multi.png --multi_inject gAMA,sBIT --offset 90B --payload 1234
	fmt.Fprintf(os.Stderr, "Flags: %s {OPTION]...\n", os.Args[0])
	flags.PrintDefaults()
	os.Exit(0)
}

func main() {
	dat, err := os.Open(opts.Input)
	defer dat.Close()
	bReader, err := utils.PreProcessImage(dat)
	if err != nil {
		log.Fatal(err)
	}
	if opts.Jpeg{
		png.ProcessImageJpeg(bReader,&opts)
	}else if opts.Png{
		png.ProcessImage(bReader, &opts)
	}
	
}

package utils

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	//"encoding/hex"

	"github.com/blackhat-go/bhg/ch-13/imgInject/models"
)

//WriteData writes new data to offset
func WriteData(r *bytes.Reader, c *models.CmdLineOpts, b []byte) {
	offset, err := strconv.ParseInt(c.Offset, 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	w, err := os.OpenFile(c.Output, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal("Fatal: Problem writing to the output file!")
	}
	r.Seek(0, 0)

	var buff = make([]byte, offset)
	r.Read(buff)
	w.Write(buff)
	w.Write(b)
	if c.Decode {
		r.Seek(int64(len(b)), 1) // right bitshift to overwrite encode chunk
	}
	_, err = io.Copy(w, r)
	if err == nil {
		fmt.Printf("Success: %s created\n", c.Output)
	}
}

func WriteDataSpecific(r *bytes.Reader, c *models.CmdLineOpts, b []byte, offset int64) {
	// offset, err := strconv.ParseInt(os, 10, 64)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	w, err := os.OpenFile(c.Output, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal("Fatal: Problem writing to the output file!")
	}
	r.Seek(0, 0)

	var buff = make([]byte, offset)
	r.Read(buff)
	w.Write(buff)
	w.Write(b)
	r.Seek(int64(len(b)), 1)//skip the existing data
	_, err = io.Copy(w, r)
	if err == nil {
		fmt.Printf("Success: %s created\n", c.Output)
	}
	
}

func WriteDataMulti(r *bytes.Reader, c *models.CmdLineOpts, b [][]byte, offsets []int64, trailerChunk bool) {
	//offset, err := strconv.ParseInt(c.Offset, 10, 64)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("Offset Value: %v\n", offset)
	w, err := os.OpenFile(c.Output, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal("Fatal: Problem writing to the output file!")
	}
	r.Seek(0, 0)
	var tmpOffset int64
	tmpOffset = 0
	//for each ancillary chunk read buff, write buff, and write b[i]
	for i, each := range offsets{
		numTmp:= each - tmpOffset
		
		fmt.Printf("NumTmp: %v\n", numTmp)
		var buff = make([]byte, numTmp)
		fmt.Printf("Offset: %v\n", each)
		r.Read(buff)
		w.Write(buff)
		w.Write(b[i])
		fmt.Printf("Bytes[i]: %v\n", b[i])
		if trailerChunk && i == len(offsets)-1{

		}else{
		
			r.Seek(int64(len(b[i])), 1)//skip the existing data
		}
		tmpOffset = each + int64(len(b[i]))
	}
	
	_, err = io.Copy(w, r)
	if err == nil {
	fmt.Printf("Success: %s created\n", c.Output)
	}

}

func WriteDataJpeg(r *bytes.Reader, c *models.CmdLineOpts, b []byte) {

	w, err := os.OpenFile(c.Output, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal("Fatal: Problem writing to the output file!")
	}
	r.Seek(0, 0)

	var buff = make([]byte, 0)
	
	 for {
	 	ind, _ := r.ReadByte()
		buff = append(buff,ind)
		
		if string(ind)=="Ù"{
			r.Seek(-2,1)//because each read automatically moves it up one so seek(-1,1) would just be Ù
			ind, _ :=r.ReadByte()
			if string(ind)=="ÿ" {
				r.Seek(1,1)//only move one because it will automatically move an additional one
				break
			}
			r.Seek(1,1)
		}
	}
	w.Write(buff)
	w.Write(b)
	//r.Seek(int64(len(b)), 1)

	// if c.Decode {
	// 	r.Seek(int64(len(b)), 1) // right bitshift to overwrite encode chunk
	// }
	// _, err = io.Copy(w, r)
	// if err == nil {
	// 	fmt.Printf("Success: %s created\n", c.Output)
	// }
}


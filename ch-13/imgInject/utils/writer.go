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
	//offset, err := strconv.ParseInt(os, 10, 64)
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
	r.Seek(int64(len(b)), 1)
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


package pnglib

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"log"
	"strconv"
	"strings"

	"github.com/blackhat-go/bhg/ch-13/imgInject/models"
	"github.com/blackhat-go/bhg/ch-13/imgInject/utils"
)

const (
	endChunkType = "IEND"
	fillerAtEnd = "ThisIsToThrowOffTheScentOfTheData"
)

//Header holds the first byte (aka magic byte)
type Header struct {
	Header uint64 //  0:8
}


//Chunk represents a data byte chunk
type Chunk struct {
	Size uint32
	Type uint32
	Data []byte
	CRC  uint32
}

//MetaChunk inherits a Chunk struct
type MetaChunk struct {
	Chk    Chunk
	Offset int64
}

// var ancillaryChunks = map[string]int{
//     "gAMA": 4,
//     "sBIT": 5,
//     "bkGD": 10,
//     "hIST": 50,
//     "tRNS": 10
//     "pHYs": 500,
//     "tIME": 1000,
// 	"tEXt": 10,
// 	"zTXt": 10,
// }
var ancillaryChunks [9]string = [9]string{"gAMA", "sBIT", "bkGD", "hIST", "tRNS","pHYs","tIME","tEXt", "zTXt" }

//ProcessImage is the wrapper to parse PNG bytes
func (mc *MetaChunk) ProcessImage(b *bytes.Reader, c *models.CmdLineOpts) {
	mc.validate(b)
	if (c.Offset != "") && (c.Encode == false && c.Decode == false) {
		var m MetaChunk
		m.Chk.Data = []byte(c.Payload)
		m.Chk.Type = m.strToInt(c.Type)
		m.Chk.Size = m.createChunkSize()
		m.Chk.CRC = m.createChunkCRC()
		bm := m.marshalData()
		bmb := bm.Bytes()
		fmt.Printf("Payload Original: % X\n", []byte(c.Payload))
		fmt.Printf("Payload: % X\n", m.Chk.Data)
		utils.WriteData(b, c, bmb)
	}
	if (c.Offset != "") && c.Encode {
		var m MetaChunk
		m.Chk.Data = utils.XorEncode([]byte(c.Payload), c.Key)
		m.Chk.Type = m.strToInt(c.Type)
		m.Chk.Size = m.createChunkSize()
		m.Chk.CRC = m.createChunkCRC()
		bm := m.marshalData()
		bmb := bm.Bytes()
		fmt.Printf("Payload Original: % X\n", []byte(c.Payload))
		fmt.Printf("Payload Encode: % X\n", m.Chk.Data)
		utils.WriteData(b, c, bmb)
	}
	if (c.Offset != "") && c.Decode {
		var m MetaChunk
		offset, _ := strconv.ParseInt(c.Offset, 10, 64)
		b.Seek(offset, 0)
		m.readChunk(b)
		origData := m.Chk.Data
		m.Chk.Data = utils.XorDecode(m.Chk.Data, c.Key)
		m.Chk.CRC = m.createChunkCRC()
		bm := m.marshalData()
		bmb := bm.Bytes()
		fmt.Printf("Payload Original: % X\n", origData)
		fmt.Printf("Payload Decode: % X\n", m.Chk.Data)
		utils.WriteData(b, c, bmb)
	}
	if c.Meta {
		count := 1 //Start at 1 because 0 is reserved for magic byte
		var chunkType string
		for chunkType != endChunkType {
			mc.getOffset(b)
			mc.readChunk(b)
			fmt.Println("---- Chunk # " + strconv.Itoa(count) + " ----")
			fmt.Printf("Chunk Offset: %#02x\n", mc.Offset)
			fmt.Printf("Chunk Length: %s bytes\n", strconv.Itoa(int(mc.Chk.Size)))
			fmt.Printf("Chunk Type: %s\n", mc.chunkTypeToString())
			fmt.Printf("Chunk Importance: %s\n", mc.checkCritType())
			if c.Suppress == false {
				fmt.Printf("Chunk Data: %#x\n", mc.Chk.Data)
			} else if c.Suppress {
				fmt.Printf("Chunk Data: %s\n", "Suppressed")
			}
			fmt.Printf("Chunk CRC: %x\n", mc.Chk.CRC)
			chunkType = mc.chunkTypeToString()
			count++
		}
	}
	if c.Specific != "" {
		var chunkType string
		anChunk := false
		for chunkType != endChunkType {
			mc.getOffset(b)
			mc.readChunk(b)

			if mc.chunkTypeToString() == c.Specific{
				anChunk = true
				var m MetaChunk
				m.Chk.Data = []byte(c.Payload)
				m.Chk.Type = mc.Chk.Type
				m.Chk.Size = m.createChunkSize()
				if m.Chk.Size > mc.Chk.Size {
					fmt.Printf("The payload size is too large for the chosen ancillary chunk")
					return
				}
				m.Chk.CRC = mc.Chk.CRC
				bm := m.marshalData()
				bmb := bm.Bytes()
				// fmt.Printf("Payload Original: % X\n", []byte(c.Payload))
				// fmt.Printf("Payload: % X\n", m.Chk.Data)
				utils.WriteDataSpecific(b, c, bmb, mc.Offset)
			}
		}
		if !anChunk{
			fmt.Printf("Your chosen an ancillary chunk DOES NOT EXIST for this png file\n")
		}
	}
	//gAMA 4 bytes
	// 67 41 4d 41
}

func (mc *MetaChunk) ProcessImageJpeg(b *bytes.Reader, c *models.CmdLineOpts) {
	mc.validateJpeg(b)
	fmt.Printf("c.Decode: %v", c.Decode)
	if c.Encode == false && c.Decode == false{
		var m MetaChunk
		m.Chk.Data = []byte(c.Payload + fillerAtEnd)
		bm := m.marshalDataJpeg()
		bmb := bm.Bytes()
		fmt.Printf("Payload Original: % X\n", []byte(c.Payload))
		fmt.Printf("Payload: % X\n", m.Chk.Data)
		utils.WriteDataJpeg(b, c, bmb)
	}else if c.Encode == true {
		var m MetaChunk
		m.Chk.Data = utils.XorEncode([]byte(c.Payload+fillerAtEnd), c.Key)
		bm := m.marshalDataJpeg()
		bmb := bm.Bytes()
		fmt.Printf("Payload Original: % X\n", []byte(c.Payload))
		fmt.Printf("Payload Encode: % X\n", m.Chk.Data)
		utils.WriteDataJpeg(b, c, bmb)
	}else if c.Decode == true {
		var m MetaChunk
		//find end file marker
		var offset uint32
		offset=2;//because 2 bytes were already read in the validate
		for {
			ind, _ := b.ReadByte()
			offset++
			if string(ind)=="Ù"{
				b.Seek(-2,1)//because each read automatically moves it up one so seek(-1,1) would just be Ù
				ind, _ :=b.ReadByte()
				if string(ind)=="ÿ" {
					b.Seek(1,1)//only move one because it will automatically move an additional one
					break
				}
				b.Seek(1,1)
			}
		}
		//m.readChunk(b)
		imageSize := b.Size()
		m.readChunkBytes(b, uint32(imageSize)-offset)
		fmt.Printf("Made it Here after FOR Loop1\n")
		origData := m.Chk.Data
		fmt.Printf("Made it Here after FOR Loop2\n")
		m.Chk.Data = utils.XorDecode(m.Chk.Data, c.Key)
		fmt.Printf("Made it Here after FOR Loop3\n")
		bm := m.marshalDataJpeg()
		fmt.Printf("Made it Here after FOR Loop4\n")
		bmb := bm.Bytes()
		fmt.Printf("Payload Original: % X\n", origData)
		fmt.Printf("Payload Decode: % X\n", m.Chk.Data)
		utils.WriteDataJpeg(b, c, bmb)
	}
	
	
	//gAMA 4 bytes
	// 67 41 4d 41
}

func (mc *MetaChunk) marshalData() *bytes.Buffer {
	bytesMSB := new(bytes.Buffer)
	if err := binary.Write(bytesMSB, binary.BigEndian, mc.Chk.Size); err != nil {
		log.Fatal(err)
	}
	if err := binary.Write(bytesMSB, binary.BigEndian, mc.Chk.Type); err != nil {
		log.Fatal(err)
	}
	if err := binary.Write(bytesMSB, binary.BigEndian, mc.Chk.Data); err != nil {
		log.Fatal(err)
	}
	if err := binary.Write(bytesMSB, binary.BigEndian, mc.Chk.CRC); err != nil {
		log.Fatal(err)
	}

	return bytesMSB
}

func (mc *MetaChunk) marshalDataJpeg() *bytes.Buffer {
	bytesMSB := new(bytes.Buffer)

	if err := binary.Write(bytesMSB, binary.BigEndian, mc.Chk.Data); err != nil {
		log.Fatal(err)
	}
	
	return bytesMSB
}

func (mc *MetaChunk) readChunk(b *bytes.Reader) {
	mc.readChunkSize(b)
	mc.readChunkType(b)
	mc.readChunkBytes(b, mc.Chk.Size)
	mc.readChunkCRC(b)
}

func (mc *MetaChunk) readChunkSize(b *bytes.Reader) {
	if err := binary.Read(b, binary.BigEndian, &mc.Chk.Size); err != nil {
		log.Fatal(err)
	}
}

func (mc *MetaChunk) readChunkType(b *bytes.Reader) {
	if err := binary.Read(b, binary.BigEndian, &mc.Chk.Type); err != nil {
		log.Fatal(err)
	}
}

func (mc *MetaChunk) readChunkBytes(b *bytes.Reader, cLen uint32) {
	mc.Chk.Data = make([]byte, cLen)
	if err := binary.Read(b, binary.BigEndian, &mc.Chk.Data); err != nil {
		log.Fatal(err)
	}
}

func (mc *MetaChunk) readChunkCRC(b *bytes.Reader) {
	if err := binary.Read(b, binary.BigEndian, &mc.Chk.CRC); err != nil {
		log.Fatal(err)
	}
}

func (mc *MetaChunk) getOffset(b *bytes.Reader) {
	offset, _ := b.Seek(0, 1)
	mc.Offset = offset
}

func (mc *MetaChunk) chunkTypeToString() string {
	h := fmt.Sprintf("%x", mc.Chk.Type)
	decoded, _ := hex.DecodeString(h)
	result := fmt.Sprintf("%s", decoded)
	return result
}

func (mc *MetaChunk) checkCritType() string {
	fChar := string([]rune(mc.chunkTypeToString())[0])
	if fChar == strings.ToUpper(fChar) {
		return "Critical"
	}
	return "Ancillary"
}

func (mc *MetaChunk) validate(b *bytes.Reader) {
	var header Header

	if err := binary.Read(b, binary.BigEndian, &header.Header); err != nil {
		log.Fatal(err)
	}

	bArr := make([]byte, 8)
	binary.BigEndian.PutUint64(bArr, header.Header)

	if string(bArr[1:4]) != "PNG" {
		log.Fatal("Provided file is not a valid PNG format")
	} else {
		fmt.Println("Valid PNG so let us continue!")
	}
}

func (mc *MetaChunk) validateJpeg(b *bytes.Reader) {
	//var header Header
	var header [2]byte 

	if err := binary.Read(b, binary.BigEndian, &header); err != nil {
		log.Fatal(err)
	}

	//bArr := make([]byte, 2)
	//binary.BigEndian.PutUint64(bArr, header)
	if string(header[0]) != "ÿ" || string(header[1]) != "Ø"{
		log.Fatal("Provided file is not a valid JPEG format")
	} else {
		fmt.Println("Valid JPEG so let us continue!")
	}
}


func (mc *MetaChunk) createChunkSize() uint32 {
	return uint32(len(mc.Chk.Data))
}

func (mc *MetaChunk) createChunkCRC() uint32 {
	bytesMSB := new(bytes.Buffer)
	if err := binary.Write(bytesMSB, binary.BigEndian, mc.Chk.Type); err != nil {
		log.Fatal(err)
	}
	if err := binary.Write(bytesMSB, binary.BigEndian, mc.Chk.Data); err != nil {
		log.Fatal(err)
	}
	return crc32.ChecksumIEEE(bytesMSB.Bytes())
}

func (mc *MetaChunk) strToInt(s string) uint32 {
	t := []byte(s)
	return binary.BigEndian.Uint32(t)
}

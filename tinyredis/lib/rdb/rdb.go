package rdb

import (
	"fmt"
	"io"
)

// RDBLoad -
func RDBLoad(filename string) {
	var rdb rioI

	rdb = newRio(filename)
	defer rdb.releaseRio() // recover panic in this

	rdbLoadRio(rdb)
}

func rdbLoadRio(rdb rioI) {
	rdbSegmentPrint()
	rdbConvertHeaderPrint()

	// segment:: magic string
	rdbSegmentPrint()
	_, rdbver := rdb.rdbLoadMagicString()
	rdb.rioSetRDBVersion(rdbver)

	// segment:: data
	rdbSegmentPrint()
	for {
		// read type of operation
		optype := rdb.rdbLoadType()

		// find a segment handler and excute it
		if handler, exist := opcodeHandlerMap[optype]; exist {
			err := handler(rdb)
			if err != nil { // include io.EOF
				if err == io.EOF {
					break
				}
				panic(fmt.Sprintln("[rdbLoadRio] opcode:", optype, " => handler() err:", err))
			}
			continue
		}

		// no segment handler found, then read rdbtype data
		rdbTypeCommonHandler(rdb, int(optype))
	}

	// segment:: crc 64 checksum
	rdbSegmentPrint()
	rdb.rdbLoadCRC64Checksum(rdbver)

	rdbSegmentPrint()
}

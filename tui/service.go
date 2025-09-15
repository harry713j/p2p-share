package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
)

// dynamic/private port range (49152–65535)
// TODO: I need to call the server for creating the session details

func sendFile(filepath string, port string) error {
	// open the file
	file, err := os.Open(filepath)

	if err != nil {
		return err
	}

	defer file.Close()

	// get the file info
	info, err := file.Stat()

	if err != nil {
		return err
	}

	metadata := Metadata{
		FileName: info.Name(),
		FileSize: info.Size(),
	}

	// send metadata first to the reciever
	listener, err := net.Listen("tcp", ":"+port) // create a socket

	if err != nil {
		return err
	}

	fmt.Printf("Waiting for reciever on port %v...\n", port)

	conn, err := listener.Accept() // listen for connection

	if err != nil {
		return err
	}

	defer conn.Close()

	// send the metadata
	metaBytes, _ := json.Marshal(metadata)
	metaLen := int32(len(metaBytes))

	binary.Write(conn, binary.BigEndian, metaLen) // 4bytes by big-endian
	conn.Write(metaBytes)

	// send the file in chunks
	buf := make([]byte, 4096) // send 4KB chunks

	for {
		n, err := file.Read(buf)

		if err == io.EOF {
			break
		}

		conn.Write(buf[:n])

	}

	return nil
}

func receiveFile(addr, port string) error {
	conn, err := net.Dial("tcp", addr+":"+port)

	if err != nil {
		return err
	}

	defer conn.Close()

	var metaLen int32
	binary.Read(conn, binary.BigEndian, &metaLen) // read the 4kbyte

	metaBytes := make([]byte, metaLen)
	io.ReadFull(conn, metaBytes) // read exactly 4KB

	var metadata Metadata
	json.Unmarshal(metaBytes, &metadata)

	// create output file
	outFIle, err := os.Create(metadata.FileName)

	if err != nil {
		return err
	}

	defer outFIle.Close()

	// read 4KB at a time
	buf := make([]byte, 4096)
	var received int64

	if received < metadata.FileSize {
		n, _ := conn.Read(buf)

		outFIle.Write(buf[:n])
		received += int64(n)
		fmt.Printf("\rProgress: %.2f%%", float64(received)/float64(metadata.FileSize)*100)
	}

	return nil
}

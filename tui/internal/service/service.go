package service

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/harry713j/p2p-share/tui/internal/config"
	"github.com/harry713j/p2p-share/tui/internal/util"
)

type Metadata struct {
	FileName  string    `json:"file_name"`
	FileSize  int64     `json:"file_size"`
	IP        string    `json:"ip"`
	Port      string    `json:"port"`
	Code      string    `json:"code"`
	ExpiresAt time.Time `json:"expires_at"`
	Checksum  []byte    `json:"checksum"`
}

type RegisterResp struct {
	Message string        `json:"message"`
	Code    string        `json:"code"`
	Timeout time.Duration `json:"timeout"`
}

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

	u := util.NewUtility()
	code := u.GetRandomCode(6)

	checksum, err := u.GenerateChecksum(file)
	if err != nil {
		return err
	}

	localIp, err := u.GetLocalIP()

	if err != nil {
		return err
	}

	metadata := Metadata{
		FileName:  info.Name(),
		FileSize:  info.Size(),
		Port:      port,
		IP:        localIp.String(),
		ExpiresAt: time.Now().Add(5 * time.Minute),
		Code:      code,
		Checksum:  checksum,
	}

	data, err := json.Marshal(metadata)

	if err != nil {
		return err
	}

	// send to the server
	resp, err := http.Post(config.ServerURL+"/register", "application/json", bytes.NewBuffer(data))

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	var srvResp RegisterResp

	err = json.Unmarshal(respData, &srvResp)

	if err != nil {
		return err
	}

	// send metadata first to the reciever
	listener, err := net.Listen("tcp", ":"+port) // create a socket

	if err != nil {
		return err
	}

	fmt.Printf("Code: %v\n", srvResp.Code)
	fmt.Printf("Timeout Duration: %v\n", srvResp.Timeout)
	fmt.Printf("Waiting for reciever on port %v...\n", port)

	conn, err := listener.Accept() // listen for connection

	if err != nil {
		return err
	}

	defer conn.Close()

	// send the file data in chunks

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

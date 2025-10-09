package service

import (
	"bytes"
	"encoding/json"
	"errors"
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
	Checksum  string    `json:"checksum"`
}

type RegisterResp struct {
	Message string        `json:"message"`
	Code    string        `json:"code"`
	Timeout time.Duration `json:"timeout"`
}

type QueryResp struct {
	FileName string `json:"file_name"`
	FileSize int64  `json:"file_size"`
	IP       string `json:"ip"`
	Port     string `json:"port"`
	Checksum string `json:"checksum"`
}

func SendFile(filepath string, port string) error {
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
	buff := make([]byte, 4096) // 4KB

	for {
		n, err := file.Read(buff)

		if err != nil {
			if err == io.EOF {
				break
			}

			return err
		}

		numOfBytesWritten, err := conn.Write(buff[:n])

		if err != nil {
			return err
		}

		fmt.Printf("%d bytes written to peer", numOfBytesWritten)
	}

	return nil
}

func ReceiveFile(code string) error {
	if len(code) != 6 {
		return errors.New("invalid code")
	}

	// get the file metadata from server
	resp, err := http.Get(config.ServerURL + "/session?code=" + code)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var queryResp QueryResp
	err = json.Unmarshal(respData, &queryResp)

	if err != nil {
		return err
	}

	fmt.Printf("Downlaod Size: %.2fKB", float64(queryResp.FileSize))
	fmt.Println("Download started")
	err = download(queryResp.IP, queryResp.Port, queryResp.FileName, queryResp.Checksum)

	if err != nil {
		return err
	}

	return nil
}

func download(addr, port, fileName, remoteChecksum string) error {
	conn, err := net.Dial("tcp", net.JoinHostPort(addr, port))
	if err != nil {
		return err
	}

	defer conn.Close()

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}

	defer file.Close()

	buff := make([]byte, 4096)

	for {
		n, err := conn.Read(buff)
		if err != nil {
			if err == io.EOF {
				break
			}

			return err
		}

		numOfBytesWritten, err := file.Write(buff[:n])
		if err != nil {
			return err
		}

		fmt.Printf("%d bytes written to file", numOfBytesWritten)
	}
	//verify the check sum
	u := util.NewUtility()
	localChecksum, err := u.GenerateChecksum(file)
	if err != nil {
		return err
	}

	if localChecksum != remoteChecksum {
		// remove the file
		return fmt.Errorf("❌ Checksum mismatch, file corrupted")
	}

	return nil
}

package util

import (
	"crypto/sha256"
	"io"
	"math/rand"
	"net"
	"os"
)

type Utility struct{}

func NewUtility() *Utility {
	return &Utility{}
}

// Generate random strings consisting of number and lower-case alphabates
// of minimum size 6
func (u *Utility) GetRandomCode(size int) string {
	if size < 6 {
		size = 6
	}

	alpha := "abcdefghijklmnopqrstuvwxyz"
	nums := "0123456789"
	mixed := alpha + nums
	shuffled := []rune(mixed)
	u.shuffle(shuffled)
	code := make([]rune, size)

	for i := 0; i < size; i++ {
		rIdx := rand.Intn(len(shuffled))
		code[i] = shuffled[rIdx]
	}

	return string(code)
}

// Generate Hash from the file value
func (u *Utility) GenerateChecksum(file *os.File) ([]byte, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return nil, err
	}
	checksum := hash.Sum(nil)

	return checksum, nil
}

func (u Utility) shuffle(arr []rune) {
	for idx := range arr {
		i := rand.Intn(len(arr))
		u.swap(arr, idx, i)
	}
}

func (u Utility) swap(arr []rune, s, e int) {
	arr[s], arr[e] = arr[e], arr[s]
}

// Generate dynamic/private port range (49152–65535)
func (u Utility) GetDynamicPort() int {
	return 49152 + rand.Intn(65535-49152)
}

// Get the Local IP address
func (u Utility) GetLocalIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")

	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP, nil
}

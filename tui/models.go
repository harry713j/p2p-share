package main

type Checksum string

type Metadata struct {
	FileName string `json:"fileName"`
	FileSize int64  `json:"fileSize"`
	// Checksum Checksum `json:"checksum"` // for checking data integrity
	// Compression string 	`json:"compression"` // compression type using zip, gzip or tar
}

// type Key string

type Session struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
	// Key  Key    `json:"key"` // public key for encrypt or decrypt the file (currently lets make it nil)
}

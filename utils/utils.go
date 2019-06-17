package utils

// utils contains utility functions that are used in packages
import (
	"encoding/hex"
	"math/rand"
	"os/user"
	"strconv"
	"time"

	"golang.org/x/crypto/sha3"
)

// Timestamp gets the human readable timestamp
func Timestamp() string {
	return time.Now().Format(time.RFC850)
}

// Unix gets the unix timestamp
func Unix() int64 {
	return time.Now().Unix()
}

// I64toS converts an int64 to a string
func I64toS(a int64) string {
	return strconv.FormatInt(a, 10) // s == "97" (decimal)
}

// ItoB converts an integer to a byte
func ItoB(a int) []byte {
	// need to convert int to a byte array for indexing
	string1 := strconv.Itoa(a)
	return []byte(string1)
}

// ItoS converts an integer to a string
func ItoS(a int) string {
	aStr := strconv.Itoa(a)
	return aStr
}

// BToI converts a byte to an integer
func BToI(a []byte) int {
	x, _ := strconv.Atoi(string(a))
	return x
}

// FtoS converts a float to a string
func FtoS(a float64) string {
	return strconv.FormatFloat(a, 'f', 6, 64)
	// return fmt.Sprintf("%f", a) is also possible, but slower due to the Sprintf
}

// StoF converts a string to a float
func StoF(a string) float64 {
	x, _ := strconv.ParseFloat(a, 32)
	// ignore this error since we hopefully call this in the right place
	return x
}

// StoFWithCheck converts a string to a float with checks
func StoFWithCheck(a string) (float64, error) {
	return strconv.ParseFloat(a, 32)
}

// StoI converts a string to an int
func StoI(a string) int {
	// convert string to int
	aInt, _ := strconv.Atoi(a)
	return aInt
}

// StoICheck converts a string to an int with checks
func StoICheck(a string) (int, error) {
	// convert string to int
	return strconv.Atoi(a)
}

// SHA3hash gets the SHA3-512 hash of the passed string
func SHA3hash(inputString string) string {
	byteString := sha3.Sum512([]byte(inputString))
	return hex.EncodeToString(byteString[:])
	// so now we have a SHA3hash that we can use to assign unique ids to our assets
}

// GetHomeDir gets the home directory of the user
func GetHomeDir() (string, error) {
	usr, err := user.Current()
	return usr.HomeDir, err
}

// GetRandomString gets a random string of length _n_
func GetRandomString(n int) string {
	// random string implementation courtesy: icza
	// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
	const (
		letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)

	var src = rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

package ipfs

// this package contains the ipfs interacting parts
// when we are adding a file to ipfs, we either could use the javascript handler
// to call the ipfs api and then use the hash ourselves to decrypt it. Or we need to
// process a pdf file (ie build an xref table) and then convert that into an ipfs file
import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"strings"

	consts "github.com/YaleOpenLab/openx/consts"
	utils "github.com/YaleOpenLab/openx/utils"
	shell "github.com/ipfs/go-ipfs-api"
)

// RetrieveShell retrieves the ipfs shell for use by other functions
// the path must be set to the rpc port used by the local / remote host
func RetrieveShell() *shell.Shell {
	// this is the api endpoint of the ipfs daemon
	return shell.NewShell("localhost:5001")
}

// AddStringToIpfs stores the given s tring in ipfs and returns
// the hash of the string
func AddStringToIpfs(a string) (string, error) {
	sh := RetrieveShell()
	hash, err := sh.Add(strings.NewReader(a)) // input must be an io.Reader
	if err != nil {
		log.Println("Error while adding string to ipfs", err)
		return "", err
	}
	return hash, nil
}

// GetFileFromIpfs gets back the contents of an ipfs hash and stores them
// in the required extension format. This has to match with the extension
// format that the original file had or else one would not be able to view
// the file
func GetFileFromIpfs(hash string, extension string) error {
	// extension can be pdf, txt, ppt and others
	sh := RetrieveShell()
	// generate a random fileName and then return the file to the user
	fileName := utils.GetRandomString(globals.IpfsFileLength) + "." + extension
	return sh.Get(hash, fileName)
}

// GetStringFromIpfs gets back the contents of an ipfs hash as a string
func GetStringFromIpfs(hash string) (string, error) {
	sh := RetrieveShell()
	// since ipfs doesn't provide a method to read the string directly, we create a
	// random fiel at tmp/, decrypt contents to that fiel and then read the file
	// contents from there
	tmpFileDir := "/tmp/" + utils.GetRandomString(globals.IpfsFileLength) // using the same length here for consistency
	sh.Get(hash, tmpFileDir)
	data, err := ioutil.ReadFile(tmpFileDir)
	if err != nil {
		log.Println("Error while reading file", err)
		return "", err
	}
	os.Remove(tmpFileDir)
	return string(data), nil
}

// ReadfromFile reads a pdf and returns the datastream
func ReadfromFile(filepath string) ([]byte, error) {
	return ioutil.ReadFile(filepath)
}

// IpfsHashFile returns the ipfs hash of a file
func IpfsHashFile(filepath string) (string, error) {
	var dummy string
	dataStream, err := ReadfromFile(filepath)
	if err != nil {
		log.Println("Error while reading from file", err)
		return dummy, err
	}
	// need to get the ifps hash of this data stream and return hash
	reader := bytes.NewReader(dataStream)
	sh := RetrieveShell()
	hash, err := sh.Add(reader)
	if err != nil {
		log.Println("Error while adding string to ipfs", err)
		return dummy, err
	}
	return hash, nil
}

// IpfsHashData hashes a byte string
func IpfsHashData(data []byte) (string, error) {
	var dummy string
	reader := bytes.NewReader(data)
	sh := RetrieveShell()
	hash, err := sh.Add(reader)
	if err != nil {
		log.Println("Error while adding string to ipfs", err)
		return dummy, err
	}
	return hash, nil
}

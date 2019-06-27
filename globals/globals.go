package globals

import (
	"os"
)

var (
	PrivateKey         = ""
	PrivateKeyPassword = ""
	DbDir = os.Getenv("HOME") + "/.openclimate/database"
	HomeDir = os.Getenv("HOME") + "/.openx"
)

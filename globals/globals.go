package globals

import (
	"os"
)

var (
	PrivateKey         = ""
	PrivateKeyPassword = ""
	HomeDir            = os.Getenv("HOME") + "/.openclimate"
	DbDir              = HomeDir + "/database"
	DbPath             = DbDir + "/openclimate.db"
)

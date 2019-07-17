package globals

import (
	"os"
)

var (
	HomeDir = os.Getenv("HOME") + "/.openclimate"
	DbDir   = HomeDir + "/database"
	DbPath  = DbDir + "/openclimate.db"
)

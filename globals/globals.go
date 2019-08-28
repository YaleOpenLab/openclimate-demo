package globals

import (
	"os"
)

var (
	HomeDir   = os.Getenv("HOME") + "/.openclimate"
	DbDir     = HomeDir + "/database"
	DbPath    = DbDir + "/openclimate.db"
	StDataDir = "staticdata/json_data"
	DefaultRpcPort = 8001
)

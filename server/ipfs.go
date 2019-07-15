package server

import (

	// "net/http"

	// erpc "github.com/Varunram/essentials/rpc"
	// ipfs "github.com/Varunram/essentials/ipfs"

)

type IpfsData struct {

	/* Meta-data */

	UserIndex		int // the index of the user the data is associated with
	// shows whether further verification is needed before commiting data to blockchain,
	// whether that's from using an oracle or checked by a third-party
	Verified		bool
	// determines how the data achieved its verification rating (or why it isn't verified)
	Methodology 	string 
	// gives us source of data, e.g. indirect self-reporting, direct IoT data, etc.
	DataSource		string 
	
	Data 			[]byte

}



package server

import (

	ipfs "github.com/Varunram/essentials/ipfs"

)

func ReportAndCommitData() {
	http.HandleFunc("/user/report", func(w http.ResponseWriter, r *http.Request){
		err := erpc.CheckGet(w, r)
		if err != nil {
			return
		}
	})
}
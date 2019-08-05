package database

type Actor interface {
	SetID(id int)
	GetID() int
	AddPledge(pledge Pledge)

}


type BucketItem interface {
	SetID(id int)
}


type RepData struct {
	// pledge, emissions, mitigation, adaption, etc.
	ReportType string
	Year       int
	IpfsHash   string
}
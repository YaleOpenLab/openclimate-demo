package database

// type Actor interface {
	
// }

type RepData struct {
	// pledge, emissions, mitigation, adaption, etc.
	ReportType string
	Year       int
	IpfsHash   string
}

/***********************/
/* PLEDGE DATA STRUCTS */
/***********************/

type Pledge struct {
	// * emissions reductions
	// * mitigation actions (energy efficiency, renewables, etc.)
	// * adaptation actions
	PledgeType string
	BaseYear   int
	TargetYear int
	Goal       int
	// is this goal determined by a regulator, or voluntarily
	// adopted by the climate actor?
	Regulatory bool
}

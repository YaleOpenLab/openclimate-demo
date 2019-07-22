package oracle

import ()

/*

	Defines the structure of all the data we want to commit to IPFS.

	*	The oracle should process inputs from self-reported data, external databases,
		static files uploaded to the server, and all other data.
	*	The oracle should produce clean, verified data in the format defined
		by the following structs.

*/

/***********************/
/* PLEDGE DATA STRUCTS */
/***********************/

type Pledges struct {

	// Meta-data
	UserID     int
	EntityType string

	// Info on specific pledges
	Pledges []PledgeData
}

type PledgeData struct {
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

/**************************/
/* EMISSIONS DATA STRUCTS */
/**************************/

type Emissions struct {
	// Meta-data
	UserID     int
	EntityType string
	Year       int
	// Emissions data (by asset)
	// Country children: regions
	// Region children: companies & cities
	// Company children: assets
	ByChild []EmissionsChild
}

type EmissionsChild struct {
	ChildID      int
	ChildName    string
	ScopeICO2e   float64
	ScopeIICO2e  float64
	ScopeIIICO2e float64

	// Where is the report and its data from?
	// (options: internally conducted report, consulting group, etc.)
	Source string

	// what methodology was used in the reporting and
	// verification of the emissions data?
	Methodology string

	// // "verified" represents if the data is sufficiently reviewed
	// // and confirmed/corroborated (from oracle, third-party auditor, etc)
	// Verified string
}

/***************************/
/* MITIGATION DATA STRUCTS */
/***************************/

type Mitigation struct {
	// Meta-data
	UserID     int
	EntityType string
	Year       int

	// Emissions data (by asset)
	// Country children: regions
	// Region children: companies & cities
	// Company children: assets
	ByChild []MitigationChild
}

type MitigationChild struct {
	ChildID      int
	ChildName    string
	CarbonOffset float64
	EnergySaved  float64
	EnergyGen    float64

	// Options:
	// - Renewable energy
	// - Energy efficiency
	// - Agriculture, Forestry & Other
	// - Carbon sequestrations
	Type string

	// Options:
	// - Reduction of GHG sources
	// - Enhancement of GHG sinks
	// - Both
	Category string


	// Where is the report and its data from?
	// (options: internally conducted report, consulting group, etc.)
	Source string

	// what methodology was used in the reporting and
	// verification of the mitigation data?
	Methodology string
}

/***************************/
/* ADAPTATION DATA STRUCTS */
/***************************/

type Adaptation struct {
}

type AdaptationChild struct {
}

/**********************/
/* EARTH DATA STRUCTS */
/**********************/

type Earth struct {

	Source string
	
	AtmosCO2 float64
	TropOzone float64 // tropospheric ozone concentration
	StratOzone float64 // stratospheric ozone concentration

	GlobalTemp float64
	ArcticIceMin float64
	IceSheets float64
	SeaLevelRise float64

	LandUse float64

}


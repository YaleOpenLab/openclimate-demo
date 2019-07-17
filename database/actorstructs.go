package database

type RepData struct {
	// pledge, emissions, mitigation, adaption, etc.
	ReportType string
	Year       int
	IpfsHash   string
}

// Data aggregated from children; ex: The U.S. user will aggregate
// the emissions of all its states/regions to get its sum total
// of emissions for the whole country.
type AggEmiData struct {
	ScopeICO2e   float64
	ScopeIICO2e  float64
	ScopeIIICO2e float64
}

type AggMitData struct {
}

type AggAdptData struct {
}

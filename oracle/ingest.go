package oracle

import (
)

/*

	Ingests data from self-reporting on the platform, external databases,
	and static data. Functions will clean the data then send to oracle
	for verification and filtering.

*/



/* Actor Data Ingesters */

// Data ingesters should read from a file (or the DB, 
// depending on where the rpc handlers put the data)

func IngestEmissions() error {
	return nil
}

func IngestPledges() error {
	return nil
}

func IngestMitigation() error {
	return nil
}

func IngestAdaptation() error {
	return nil
}


/* Earth Data Ingesters */

func IngestGlobalEmissions() error {
	return nil
}
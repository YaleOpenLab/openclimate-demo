package database

import ()

type LandingPageInfo struct {
	DataPoints []struct {
		CompanyName      string
		InstallationName string
		Address          string
		Emission         string
		MoreInfoLink     string
	}
}

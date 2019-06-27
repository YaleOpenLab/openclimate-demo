package database

import ()

type Company struct {
	Dashboard struct {
		Name                   string
		ScopeMetrics           []string // emissions, mitigation, adaptation
		DirectEmissions        float64
		DirectEmissionsLink    string
		MitigationOutcomes     float64
		MitigationOutcomesLink string
		WindSolar              float64
		WindSolarLink          string
		Adaptation             float64
		Pledges                [][]string // ["Net Emissions"]["Carbon Neutral by 2025"]
		YourReportingProfile   struct {
			PublicDisclosure [][]string // ["Accounting visibility"]["Aggregate Accounting"]
		}
		WeightedScore         int // 2 out of 3 stars
		ClimateAccountability struct {
			Direct      float64
			Indirect    float64
			Untrackable float64
		}
	}
	View struct {
		Earth struct {
			CWarming        string
			CO2PPM          string
			GtCO2Left       string
			GtCo2Year       string
			EarthStatusLink string
		}
		NationState struct {
			Country   string
			NDCPledge string
		}
		Subnational struct {
			Country string
			State   string
		}
		ClimateActionAssets struct {
			Mitigation struct {
				RenewableEnergy bool
			}
			Class1 []struct {
				Name                    string
				Type                    string
				ScopeMetrics            []string
				Capacity                string
				MitigationOutcomes      string
				GWhYer                  string
				CertificateAssetsToDate string
				ReportingDevice         string
				ActiveIssues            []string
				MRVMethodology          string
				MRVMethodologyLink      string
				BlueProgressBar         float64
				StarRating              int
				Certificates            struct {
					CertificateId string
					Type          string
					Unit          string
					Start         string
					End           string
					Status        string
				}
				AccessTerminal     string
				AccessTerminalLink string
			}
			EditNestedScopes string
		}
	}
	Review struct {
		CarbonBalance struct {
			MtCO2Year   string
			Emissions   string
			Reductions  string
			LastUpdated string
			ReviewLink  string

			RandomBlock struct {
				CertID string
				Type   string
				Unit   string
				Status string
			}
			BuyTradeLink string
		}

		CLimateReports struct {
			Name                 string
			Scope                string
			Date                 string
			Verified             string
			DownloadLink         string
			SeeAllLink           string
			NewClimateReportLink string
			ExportData           string
		}

		IssuesAndDeposits struct {
			All struct {
				IssueName    string
				IssueId      string
				IssueDate    string
				IssueCreator string
				Tags         []string
				Author       string
				Labels       []string
				Assets       []string
				Started      bool
				Assignee     string
			}
		}

		Manage struct {
			ClimateActions struct {
				StartNewLink      string
				AddExistingLink   string
				BulkIntegrateLink string
			}

			ManageCAP struct {
				Type            []string
				Action          string
				Name            string
				Quantity        string
				Region          string
				Note            string
				MRVProcess      string
				BlueProgressBar float64
				Rating          int
			}
		}
	}
}

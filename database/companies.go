package database

import (
	"encoding/json"
	utils "github.com/Varunram/essentials/utils"
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

type Company struct {
	Index     int
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

// RetrieveUser retrieves a particular User indexed by key from the database
func RetrieveCompany(key int) (Company, error) {
	var company Company
	db, err := OpenDB()
	if err != nil {
		return company, errors.Wrap(err, "error while opening database")
	}
	defer db.Close()
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(UserBucket)
		x := b.Get(utils.ItoB(key))
		if x == nil {
			return errors.New("retrieved user nil, quitting!")
		}
		return json.Unmarshal(x, &company)
	})
	return company, err
}

// ValidateUser validates a particular user
func RetrieveCompanyByName(name string) (Company, error) {
	var company Company
	temp, err := RetrieveAllUsers()
	if err != nil {
		return company, errors.Wrap(err, "error while retrieving all users from database")
	}
	limit := len(temp) + 1
	db, err := OpenDB()
	if err != nil {
		return company, errors.Wrap(err, "could not open db, quitting!")
	}
	defer db.Close()
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(UserBucket)
		for i := 1; i < limit; i++ {
			var rCompany Company
			x := b.Get(utils.ItoB(i))
			err := json.Unmarshal(x, &rCompany)
			if err != nil {
				return errors.Wrap(err, "could not unmarshal json, quitting!")
			}
			// check name
			if rCompany.Dashboard.Name == name {
				company = rCompany
				return nil
			}
		}
		return errors.New("Not Found")
	})
	return company, err
}

// RetrieveAllUsers gets a list of all User in the database
func RetrieveAllCompanies() ([]Company, error) {
	var arr []Company
	db, err := OpenDB()
	if err != nil {
		return arr, errors.Wrap(err, "Error while opening database")
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(UserBucket)
		for i := 1; ; i++ {
			var rCompany Company
			x := b.Get(utils.ItoB(i))
			if x == nil {
				return nil
			}
			err := json.Unmarshal(x, &rCompany)
			if err != nil {
				return errors.Wrap(err, "Error while unmarshalling json")
			}
			arr = append(arr, rCompany)
		}
	})
	return arr, err
}

// Save inserts a passed User object into the database
func (a *Company) Save() error {
	db, err := OpenDB()
	if err != nil {
		return errors.Wrap(err, "Error while opening database")
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(UserBucket)
		encoded, err := json.Marshal(a)
		if err != nil {
			return errors.Wrap(err, "Error while marshaling json")
		}
		return b.Put([]byte(utils.ItoB(a.Index)), encoded)
	})
	return err
}

func NewCompany(name string) (Company, error) {
	var x Company

	companies, err := RetrieveAllCompanies()
	if err != nil {
		return x, errors.Wrap(err, "could not retrieve all companies, quitting")
	}

	if len(companies) == 0 {
		x.Index = 1
	} else {
		x.Index = len(companies) + 1
	}

	x.Dashboard.Name = name
	return x, x.Save()
}

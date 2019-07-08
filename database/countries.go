package database


type CountryPrelim struct {
	Nation			string
	Year			int
	Total			int
	SolidFuel		float64
	LiquidFuel		float64
	GasFuel			float64
	Cement			int
	GasFlaring		float64
	PerCapita		float64
	Bunkers			int
}

type CountryFinal struct {
	Year []struct {
		Total		int
		SolidFuel	float64
		LiquidFuel	float64
		GasFuel		float64
		Cement		int
		GasFlaring	float64
		PerCapita	float64
		Bunkers		int
	}

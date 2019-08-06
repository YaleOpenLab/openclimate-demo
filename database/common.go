package database

import (
	globals "github.com/YaleOpenLab/openclimate/globals"
)

type Actor interface {
	RetrievePledges() ([]Pledge, error)
	// AddPledge(pledge Pledge)
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

// Puts asset object in assets bucket. Called by NewAsset
func (a *Asset) Save() error {
	return Save(globals.DbPath, AssetBucket, a)
}

// Saves city object in cities bucket. Called by NewCity
func (city *City) Save() error {
	return Save(globals.DbPath, CityBucket, city)
}

// Saves country object in countries bucket. Called by NewCountry
func (country *Country) Save() error {
	return Save(globals.DbPath, CountryBucket, country)
}

func (o *Oversight) Save() error {
	return Save(globals.DbPath, OversightBucket, o)
}

func (p *Pledge) Save() error {
	return Save(globals.DbPath, PledgeBucket, p)
}

// Saves region object in regions bucket. Called by NewRegion
func (region *Region) Save() error {
	return Save(globals.DbPath, RegionBucket, region)
}

func (cr *ConnectRequest) Save() error {
	return Save(globals.DbPath, RequestBucket, cr)
}

// Saves state object in states bucket. Called by NewState
func (state *State) Save() error {
	return Save(globals.DbPath, StateBucket, state)
}

// Save inserts a passed User object into the database
func (u *User) Save() error {
	return Save(globals.DbPath, UserBucket, u)
}

// Saves company object in companies bucket. Called by NewCompany
func (c *Company) Save() error {
	return Save(globals.DbPath, CompanyBucket, c)
}

func (c *Company) SetID(id int) {
	c.Index = id
}

func (a *Asset) SetID(id int) {
	a.Index = id
}

func (c *City) SetID(id int) {
	c.Index = id
}

func (c *Country) SetID(id int) {
	c.Index = id
}

func (o *Oversight) SetID(id int) {
	o.Index = id
}

func (p *Pledge) SetID(id int) {
	p.ID = id
}

func (r *Region) SetID(id int) {
	r.Index = id
}

func (cr *ConnectRequest) SetID(id int) {
	cr.Index = id
}

func (s *State) SetID(id int) {
	s.Index = id
}

func (u *User) SetID(id int) {
	u.Index = id
}

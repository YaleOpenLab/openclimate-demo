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
	GetID() int
	Save() error
}

type RepData struct {
	// pledge, emissions, mitigation, adaption, etc.
	ReportType string
	Year       int
	IpfsHash   string
}

// Puts asset object in assets bucket. Called by NewAsset
func (x *Asset) Save() error {
	return Save(globals.DbPath, AssetBucket, x)
}

// Saves city object in cities bucket. Called by NewCity
func (x *City) Save() error {
	return Save(globals.DbPath, CityBucket, x)
}

// Saves country object in countries bucket. Called by NewCountry
func (x *Country) Save() error {
	return Save(globals.DbPath, CountryBucket, x)
}

func (x *Oversight) Save() error {
	return Save(globals.DbPath, OversightBucket, x)
}

func (x *Pledge) Save() error {
	return Save(globals.DbPath, PledgeBucket, x)
}

// Saves region object in regions bucket. Called by NewRegion
func (x *Region) Save() error {
	return Save(globals.DbPath, RegionBucket, x)
}

func (x *ConnectRequest) Save() error {
	return Save(globals.DbPath, RequestBucket, x)
}

// Saves state object in states bucket. Called by NewState
func (x *State) Save() error {
	return Save(globals.DbPath, StateBucket, x)
}

// Save inserts a passed User object into the database
func (x *User) Save() error {
	return Save(globals.DbPath, UserBucket, x)
}

// Saves company object in companies bucket. Called by NewCompany
func (x *Company) Save() error {
	return Save(globals.DbPath, CompanyBucket, x)
}

/* 
	SetID() is a common method between all structs that qualify as
	bucket items that allow them to implement the BucketItem interface.
	SetID() is a simple setter method that allows the updating of the
	bucket item's ID. The function's only use should be in the Save()
	function; otherwise, the ID should not be modified.
*/ 

func (x *Company) SetID(id int) {
	x.Index = id
}

func (x *Asset) SetID(id int) {
	x.Index = id
}

func (x *City) SetID(id int) {
	x.Index = id
}

func (x *Country) SetID(id int) {
	x.Index = id
}

func (x *Oversight) SetID(id int) {
	x.Index = id
}

func (x *Pledge) SetID(id int) {
	x.ID = id
}

func (x *Region) SetID(id int) {
	x.Index = id
}

func (x *ConnectRequest) SetID(id int) {
	x.Index = id
}

func (x *State) SetID(id int) {
	x.Index = id
}

func (x *User) SetID(id int) {
	x.Index = id
}

func (x *Company) GetID() int {
	return x.Index
}

func (x *Asset) GetID() int {
	return x.Index
}

func (x *City) GetID() int {
	return x.Index
}

func (x *Country) GetID() int {
	return x.Index
}

func (x *Oversight) GetID() int {
	return x.Index
}

func (x *Pledge) GetID() int {
	return x.ID
}

func (x *Region) GetID() int {
	return x.Index
}

func (x *ConnectRequest) GetID() int {
	return x.Index
}

func (x *State) GetID() int {
	return x.Index
}

func (x *User) GetID() int {
	return x.Index
}

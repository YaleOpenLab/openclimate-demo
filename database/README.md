## database

Contains database handlers, defines structs for storing data in the database, and functions for tracking entities that interact with the platform (users, companies, national/sub-national actors, etc.).

### Folder structure

 - cities.go: Contains functions to create, save, or retrieve cities from db
 - companies.go: Contains functions to create, save, or retrieve companies from db
 - countries.go: Contains functions to create, save, or retrieve countries from db
 - database_test.go: Tests the creation and retrieval of users.
 - db.go: defines boltDB buckets and functions to handle DB
 - landing.go:
 - populate.go: Populates the local test database with static test data.
 - randomdata.go: 
 - regions.go: Contains functions to create, save, or retrieve regions from db
 - users.go:


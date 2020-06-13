# brief-api-gen
 brief-api-gen is tool for generating CRUD by database table. It helps not to spend a lot of time for writing standard code.

### Install

```sh
$ go get -u github.com/nir007/brief-api-gen
```

### Usage

brief-api-gen -connection="postgres://localhost:5432/dbname?user=guest&password=guest&sslmode=disable" -table="statuses" -en="Status" -enp="Statuses"

Will generate 3 files:
1. api_statuses.go - http handlers CRUD;
2. dto_status.go - structs for CRUD;
3. orm_status.go - functions to work with database.


### Profit

* Standard template, naming, logging, errors wrapping;
* Time saving. We can go to drink beer and fuck girls instead of writing stupid sql-queries;
* Clear code, Swagger comments;
* Easy to modify;
* Functions: filtering, search for all columns and paging;
* Validation of encoming params;
* Fast extend, join other entities etc.

# Go Echo MicroService

## Table of Contents

* [Requirements](#requirements)
* [Installation](#installation)
* [Migration](#migration)
* [Seeder](#seeder)
* [Usage](#usage)
* [Versioning](#versioning)
* [Authors](#authors)
* [Contributors](#contributors)
* [Official Documentation for Go Language](#official-documentation-for-go-language)
* [License](#license)

## Requirements

To use Go Echo MicroService, you must ensure that you meet the following requirements:
- [Go](https://golang.org/) >= 1.22

## Installation

To use Go Echo MicroService, you must follow the steps below:
- Clone a Repository
```git
# git clone https://github.com/MrAndreID/goechoms.git
```
- Get Dependancies
```go
# go mod download
# go mod tidy
```
- Create .env file from .env.example (Linux)
```go
# cp .env.example .env
```
- Configuring .env file

## Migration

To Run Migration for Go Echo MicroService, you must ensure that you meet the following requirements:
- Run Migration for Go Echo MicroService with Drop All Tables
```go
# go run applications/databases/migrations/main.go --migrate=fresh
```
- Run Migration for Go Echo MicroService
```go
# go run applications/databases/migrations/main.go --migrate=default
```

## Seeder

To Run Seeder for Go Echo MicroService, you must ensure that you meet the following requirements:
- Run Seeder for Go Echo MicroService
```go
# go run applications/databases/seeders/main.go --seed=default
```

## Usage

To Use Go Echo MicroService, you must ensure that you meet the following requirements:
- Run Go Echo MicroService
```go
# go run main.go
```

## Versioning

I use [SemVer](https://semver.org/) for versioning. For the versions available, see the tags on this repository. 

## Authors

- **Andrea Adam** - [MrAndreID](https://github.com/MrAndreID)

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.
Please make sure to update tests as appropriate.

## Official Documentation for Go Language

Documentation for Go Language can be found on the [Go Package website](https://pkg.go.dev/).

## License

Go Echo MicroService is released under the [MIT License](https://opensource.org/licenses/MIT). See the `LICENSE` file for more information.

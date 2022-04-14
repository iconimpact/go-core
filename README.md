# Go-Core

A collection of valuable packages for Golang projects. They are self-sustained and can be imported to any project.
To learn more about each package please open the subfolder and see the README.

<br>

## How to Install

```bash
go get github.com/iconimpact/go-core
```

<br>

## How to Contribute

### Development Setup

Firstly, clone the repo. And then go into the repo root folder, install dependencies and run tests.

```shell

go mod download
go test ./... -v -count=1

```

<br>

### General Contributor Guidelines

This project is maintained by [Sebastian Kreutzberger (@skreutzberger)](https://github.com/skreutzberger)

Please create a pull request which fulfills the following conditions:

1. is of general interest and to be used in more than one project
2. includes a README file which documents the purpose, each public function with an example and how to run and test it
2. contains inline code comments which make it clear what the logic does (the more comments the better)
3. has a unit test code coverage of at least 90%
4. does not violate any patents or copyrights
5. has the Apache 2.0 copyright statement at the top of each file with icommobile GmbH as copyright holder

<br>


## License
Go-Core is licensed under the Apache 2.0 License (see file `LICENSE`).

<br>

## Copyright
© 2020 iconmobile GmbH

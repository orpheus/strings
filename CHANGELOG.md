# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### [2022-07-12]
#### Added
- update string name api
- update string order api

### [2022-07-11]
#### Updated
- delete thread now deletes all associated strings first
- fetch strings by thread in order

#### Fixed
- string FindAll wasn't using the right function to get the query parameter for `thread` id

### [2022-07-09]
#### Added
- createOne/deleteById/findAll strings

### [2022-07-06]
#### Added
- createOne/deleteById/findAll threads

[unreleased]: https://github.com/olivierlacan/keep-a-changelog/compare/v1.1.0...HEAD
[1.1.0]: https://github.com/olivierlacan/keep-a-changelog/compare/v1.0.0...v1.1.0

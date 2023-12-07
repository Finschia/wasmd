# Changelog

## [Unreleased](https://github.com/Finschia/wasmd/compare/v0.2.0...HEAD)

### Features

### Improvements

### Bug Fixes

### Breaking Changes

### Build, CI

### Document Updates
* [\#113](https://github.com/Finschia/wasmd/pull/113) add compatibility of wasmd v0.2


## [v0.2.0](https://github.com/Finschia/wasmd/releases/tag/v0.2.0) - 2023.10.19

### Features
* [\#61](https://github.com/Finschia/wasmd/pull/61) bumpup ibc-go to v4
* [\#84](https://github.com/Finschia/wasmd/pull/84) add AcceptlistStaragteQuerier

### Improvements
* [\#63](https://github.com/Finschia/wasmd/pull/63) add event checking to TestStoreCode
* [\#65](https://github.com/Finschia/wasmd/pull/65) add test cases for empty request in each function
* [\#66](https://github.com/Finschia/wasmd/pull/66) add test cases for invalid pagination key in some functions
* [\#64](https://github.com/Finschia/wasmd/pull/64) test: add test cases to confirm output for PinnedCodes
* [\#70](https://github.com/Finschia/wasmd/pull/70) add event checking to TestInstantiateContract
* [\#73](https://github.com/Finschia/wasmd/pull/73) test: add the check for expPaginationTotal
* [\#72](https://github.com/Finschia/wasmd/pull/72) add pagination next key test in ContractHistory
* [\#75](https://github.com/Finschia/wasmd/pull/75) test: add the test case for InactiveContract
* [\#74](https://github.com/Finschia/wasmd/pull/74) add event checking to TestInstantiateContract2
* [\#78](https://github.com/Finschia/wasmd/pull/78) add the check for TestMigrateContract
* [\#69](https://github.com/Finschia/wasmd/pull/69) refactor: refactor test cases for Params
* [\#71](https://github.com/Finschia/wasmd/pull/71) add test cases in ContractsByCode
* [\#82](https://github.com/Finschia/wasmd/pull/82) add test case to QueryInactiveContracts
* [\#87](https://github.com/Finschia/wasmd/pull/87) add an integration test for TestExecuteContract
* [\#79](https://github.com/Finschia/wasmd/pull/79) add an integration test for TestStoreAndInstantiateContract
* [\#88](https://github.com/Finschia/wasmd/pull/88) add the test case for invalid address
* [\#76](https://github.com/Finschia/wasmd/pull/76) add an integration test for ClearAdmin
* [\#68](https://github.com/Finschia/wasmd/pull/68) add an integration test for UpdateAdmin
* [\#99](https://github.com/Finschia/wasmd/pull/99) format test files
* [\#98](https://github.com/Finschia/wasmd/pull/98) refactor TestStoreAndInstantiateContract
* [\#100](https://github.com/Finschia/wasmd/pull/100) refactor tests for cosmwasm/wasm/v1/tx.proto other than TestClearAdmin

### Bug Fixes
* [\#62](https://github.com/Finschia/wasmd/pull/62) fill ContractHistory querier result's Updated field
* [\#52](https://github.com/Finschia/wasmd/pull/52) fix cli_test error of wasmplus and add cli_test ci
* [\#89](https://github.com/Finschia/wasmd/pull/89) fill ContractInfo querier result's Updated field
* [\#90](https://github.com/Finschia/wasmd/pull/90) delete output in TestQueryContractsByCode
* [\#95](https://github.com/Finschia/wasmd/pull/95) add a test case to verify ContractInfo gets correct IBCPortID
* [\#101](https://github.com/Finschia/wasmd/pull/101) move helper method out of generated file

### Build, CI
* [\#104](https://github.com/Finschia/wasmd/pull/104) change depending wasmvm to v1.1.1-0.11.4-rc1
* [\#60](https://github.com/Finschia/wasmd/pull/60) Update golang version to 1.20
* [\#105](https://github.com/Finschia/wasmd/pull/105) change depending finschia-sdk to v0.48.0-rc2 and backport #81
* [\#109](https://github.com/Finschia/wasmd/pull/109) version bump of ostracon, finschia-sdk and wasmvm

### Document Updates
* [\#54](https://github.com/Finschia/wasmd/pull/54) add documentation about errors (codespace and codes)
* [\#92](https://github.com/Finschia/wasmd/pull/92) modify links in x/wasmplus README.md
* [\#103](https://github.com/Finschia/wasmd/pull/103) add clarifications to the order of list response in query
* [\#107](https://github.com/Finschia/wasmd/pull/107) add comment about ordering of addresses in wasmplus


## [v0.1.4](https://github.com/Finschia/wasmd/releases/tag/v0.1.4) - 2023.05.22

### Features
* [\#46](https://github.com/Finschia/wasmd/pull/46) add admin-related events

### Improvements
* [\#43](https://github.com/Finschia/wasmd/pull/43) delete unnecessary test

### Bug Fixes
* [\#35](https://github.com/Finschia/wasmd/pull/35) stop wrap twice the response of handling non-plus wasm message in plus handler
* [\#77](https://github.com/Finschia/wasmd/pull/77) use ctx cache in msg server integration test

### Document Updates
* [\#44](https://github.com/Finschia/wasmd/pull/44) update notice


## [v0.1.3](https://github.com/Finschia/wasmd/releases/tag/v0.1.3) - 2023.04.19

### Build, CI
* [\#30](https://github.com/Finschia/wasmd/pull/30) replace line repositories with finschia repositories


## [v0.1.2](https://github.com/Finschia/wasmd/releases/tag/v0.1.2) - 2023.04.10

### Features
* [\#21](https://github.com/Finschia/wasmd/pull/21) bump up Finschia/ibc-go v3.3.2


## [v0.1.0](https://github.com/Finschia/wasmd/releases/tag/v0.1.0) - 2023.03.28

### Features
* [\#9](https://github.com/Finschia/wasmd/pull/9) apply the changes of finschia-sdk and ostracon proto

### Improvements
* [\#1](https://github.com/Finschia/wasmd/pull/1) apply all changes of `x/wasm` in finschia-sdk until [finschia-sdk@3bdcb6ffe01c81615bedb777ca0e039cc46ef00c](https://github.com/Finschia/finschia-sdk/tree/3bdcb6ffe01c81615bedb777ca0e039cc46ef00c)
* [\#5](https://github.com/Finschia/wasmd/pull/5) bump up wasmd v0.29.1
* [\#7](https://github.com/Finschia/wasmd/pull/7) separate custom features in `x/wasm` into `x/wasmplus` module
* [\#8](https://github.com/Finschia/wasmd/pull/8) Bump Finschia/finschia-sdk to a7557b1d10
* [\#10](https://github.com/Finschia/wasmd/pull/10) update wasmvm version

### Bug Fixes
* [\#12](https://github.com/Finschia/wasmd/pull/12) fix not to register wrong codec in `x/wasmplus`
* [\#14](https://github.com/Finschia/wasmd/pull/14) fix the cmd error that does not recognize wasmvm library version

### Breaking Changes

### Build, CI

### Document Updates
* [\#2](https://github.com/Finschia/wasmd/pull/2) add wasm events description


## [cosmwasm/wasmd v0.27.0](https://github.com/CosmWasm/wasmd/blob/v0.27.0/CHANGELOG.md) (2022-05-19)
Initial wasmd is based on the cosmwasm/wasmd v0.27.0

* cosmwasm/wasmd [v0.27.0](https://github.com/CosmWasm/wasmd/releases/tag/v0.27.0)

Please refer [CHANGELOG_OF_COSMWASM_WASMD_v0.27.0](https://github.com/CosmWasm/wasmd/blob/v0.27.0/CHANGELOG.md)

# Changelog

## [Unreleased](https://github.com/Finschia/wasmd/compare/v0.1.4...HEAD)

### Features

### Improvements
* [\#36](https://github.com/Finschia/wasmd/pull/36) separate `x/wasm` into `x/wasmplus` module of dynamiclink

### Bug Fixes
* [\#57](https://github.com/Finschia/wasmd/pull/57) fix dynamic link APIs do not panic with invalid bech32
* [\#35](https://github.com/Finschia/wasmd/pull/35) stop wrap twice the response of handling non-plus wasm message in plus handler

### Breaking Changes
* [\#41](https://github.com/Finschia/wasmd/pull/41) add `cosmwasmAPIGenerator` to keeper
* [\#26](https://github.com/Finschia/wasmd/pull/26) implement CallCallablePoint and ValidateDynamicLinkInterface to cosmwasmAPI
* [\#29](https://github.com/Finschia/wasmd/pull/29) remove getContractEnv from cosmwasmAPI

### Build, CI

### Document Updates


## [v0.1.4](https://github.com/Finschia/wasmd/releases/tag/v0.1.4) - 2023.05.22

### Features
* [\#46](https://github.com/Finschia/wasmd/pull/46) add admin-related events

### Improvements
* [\#43](https://github.com/Finschia/wasmd/pull/43) delete unnecessary test 

### Bug Fixes
* [\#35](https://github.com/Finschia/wasmd/pull/35) stop wrap twice the response of handling non-plus wasm message in plus handler

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
* [\#18](https://github.com/Finschia/wasmd/pull/18) apply the wasm module of finschia-sdk(dynamic_link branch) until [finschia-sdk@911e8b47774f142d70d5c696722b0291e39e0c0c](https://github.com/Finschia/finschia-sdk/tree/911e8b47774f142d70d5c696722b0291e39e0c0c)

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

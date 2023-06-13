package types

import (
	storetypes "github.com/Finschia/finschia-sdk/store/types"
	wasmvm "github.com/Finschia/wasmvm"
)

var _ wasmvm.KVStore = (*WasmStore)(nil)

// WasmStore is a wrapper struct of `KVStore`
// It translates from cosmos KVStore to wasmvm-defined KVStore.
// The spec of interface `Iterator` is a bit different so we cannot use cosmos KVStore directly.
type WasmStore struct {
	storetypes.KVStore
}

// NewWasmStore creates a instance of WasmStore
func NewWasmStore(kvStore storetypes.KVStore) WasmStore {
	return WasmStore{kvStore}
}

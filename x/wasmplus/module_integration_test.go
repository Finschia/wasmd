package wasmplus_test

// todo: We need to check this test.
// func TestAppPlusModuleMigrations(t *testing.T) {
// 	wasmApp := appplus.Setup(false)
// 	ctx := wasmApp.BaseApp.NewContext(false, tmproto.Header{})
// 	upgradeHandler := func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
// 		return wasmApp.ModuleManager().RunMigrations(ctx, wasmApp.ModuleConfigurator(), fromVM)
// 	}
// 	fromVM := wasmApp.UpgradeKeeper.GetModuleVersionMap(ctx)
// 	fromVM[wasm.ModuleName] = 1 // start with initial version
// 	upgradeHandler(ctx, upgradetypes.Plan{Name: "testing"}, fromVM)
// 	// when
// 	gotVM, err := wasmApp.ModuleManager().RunMigrations(ctx, wasmApp.ModuleConfigurator(), fromVM)
// 	// then
// 	require.NoError(t, err)
// 	assert.Equal(t, uint64(1), gotVM[wasm.ModuleName])
// }

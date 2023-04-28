package keeper

func (k Keeper) GetGasRegister() GasRegister {
	return k.gasRegister
}

func (k Keeper) GetCosmwasmAPIGenerator() CosmwasmAPIGenerator {
	return k.cosmwasmAPIGenerator
}

func (k *Keeper) SetCosmwasmAPIGenerator(generator CosmwasmAPIGenerator) {
	k.cosmwasmAPIGenerator = generator
}

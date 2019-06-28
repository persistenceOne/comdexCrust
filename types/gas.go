package types

// Gas measured by the SDK
type Gas = int64

// ErrorOutOfGas : Error thrown when out of gas
type ErrorOutOfGas struct {
	Descriptor string
}

// GasMeter interface to track gas consumption
type GasMeter interface {
	GasConsumed() Gas
	ConsumeGas(amount Gas, descriptor string)
}

type basicGasMeter struct {
	limit    Gas
	consumed Gas
}

// NewGasMeter : returns New basicGasMeter
func NewGasMeter(limit Gas) GasMeter {
	return &basicGasMeter{
		limit:    limit,
		consumed: 0,
	}
}

func (g *basicGasMeter) GasConsumed() Gas {
	return g.consumed
}

func (g *basicGasMeter) ConsumeGas(amount Gas, descriptor string) {
	g.consumed += amount
	if g.consumed > g.limit {
		panic(ErrorOutOfGas{descriptor})
	}
}

type infiniteGasMeter struct {
	consumed Gas
}

// NewInfiniteGasMeter : returns infiniteGasMeter
func NewInfiniteGasMeter() GasMeter {
	return &infiniteGasMeter{
		consumed: 0,
	}
}

func (g *infiniteGasMeter) GasConsumed() Gas {
	return g.consumed
}

func (g *infiniteGasMeter) ConsumeGas(amount Gas, descriptor string) {
	g.consumed += amount
}

// GasConfig defines gas cost for each operation on KVStores
type GasConfig struct {
	HasCost          Gas
	ReadCostFlat     Gas
	ReadCostPerByte  Gas
	WriteCostFlat    Gas
	WriteCostPerByte Gas
	KeyCostFlat      Gas
	ValueCostFlat    Gas
	ValueCostPerByte Gas
}

var (
	cachedDefaultGasConfig   = DefaultGasConfig()
	cachedTransientGasConfig = TransientGasConfig()
)

// DefaultGasConfig : Default gas config for KVStores
func DefaultGasConfig() GasConfig {
	return GasConfig{
		HasCost:          0,
		ReadCostFlat:     0,
		ReadCostPerByte:  0,
		WriteCostFlat:    0,
		WriteCostPerByte: 0,
		KeyCostFlat:      0,
		ValueCostFlat:    0,
		ValueCostPerByte: 0,
	}
}

// TransientGasConfig : Default gas config for TransientStores
func TransientGasConfig() GasConfig {
	// TODO: define gasconfig for transient stores
	return DefaultGasConfig()
}

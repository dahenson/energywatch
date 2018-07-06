package energyapis

type EnergyAPI interface {
	GetName() string
	PushInstantaneousDemand(watts float64) error
	PushCurrentSummationDelivered(watthours float64) error
}

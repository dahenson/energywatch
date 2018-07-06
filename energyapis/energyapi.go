package energyapis

type EnergyAPI interface {
	PushInstantaneousDemand(watts float64) error
	PushCurrentSummationDelivered(watthours float64) error
}

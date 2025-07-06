package constant

type ActivateRegion int

const (
	ApSouthEast1 ActivateRegion = iota
	ApEast1
)

var ActivateRegionMap = map[ActivateRegion]string{
	ApSouthEast1: "ap-southeast-1",
	ApEast1:      "ap-east-1",
}

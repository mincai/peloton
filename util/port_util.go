package util

import (
	"sort"

	mesos "github.com/uber/peloton/.gen/mesos/v1"
)

// ExtractPortSet is helper function to extract available port set
// from a Mesos resource.
func ExtractPortSet(resource *mesos.Resource) map[uint32]bool {
	res := make(map[uint32]bool)

	if resource.GetName() != "ports" {
		return res
	}

	for _, r := range resource.GetRanges().GetRange() {
		// Remember that end is inclusive
		for i := r.GetBegin(); i <= r.GetEnd(); i++ {
			res[uint32(i)] = true
		}
	}

	return res
}

// GetPortsSetFromResources is helper function to extract ports resources.
func GetPortsSetFromResources(resources []*mesos.Resource) map[uint32]bool {
	res := make(map[uint32]bool)
	for _, rs := range resources {
		portSet := ExtractPortSet(rs)
		for port := range portSet {
			res[port] = true
		}
	}
	return res
}

// GetPortsNumFromOfferMap is helper function to get number of available ports
// from given id to offer map.
func GetPortsNumFromOfferMap(offerMap map[string]*mesos.Offer) uint32 {
	numPorts := 0
	for _, offer := range offerMap {
		numPorts += len(GetPortsSetFromResources(offer.GetResources()))
	}
	return uint32(numPorts)
}

// CreatePortRanges create Mesos Ranges type from given port set.
func CreatePortRanges(portSet map[uint32]bool) *mesos.Value_Ranges {
	var sorted []int
	for p, ok := range portSet {
		if ok {
			sorted = append(sorted, int(p))
		}
	}
	sort.Ints(sorted)

	res := mesos.Value_Ranges{
		Range: []*mesos.Value_Range{},
	}
	for _, p := range sorted {
		tmp := uint64(p)
		res.Range = append(
			res.Range,
			&mesos.Value_Range{Begin: &tmp, End: &tmp},
		)
	}
	return &res
}

// CreatePortResources create a list of Mesos resources suitable for launching
// from a map from port number to role name.
func CreatePortResources(portSet map[uint32]string) []*mesos.Resource {
	resources := []*mesos.Resource{}
	for port, role := range portSet {
		tmp := uint64(port)
		rs := NewMesosResourceBuilder().
			WithName("ports").
			WithType(mesos.Value_RANGES).
			WithRole(role).
			WithRanges(&mesos.Value_Ranges{
				Range: []*mesos.Value_Range{
					{Begin: &tmp, End: &tmp},
				},
			}).
			Build()
		resources = append(resources, rs)
	}
	return resources
}

package strategy

type garbageCollector struct {
	state *map[key]map[question]state
}

func newGarbageCollector(state *map[key]map[question]state) *garbageCollector {
	return &garbageCollector{
		state: state,
	}
}

func (gc *garbageCollector) collect(){

}
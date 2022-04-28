package strategy

import "time"

const (
	// garbageCollectionDefaultTTLSeconds defines the period after which the resource is removed
	garbageCollectionDefaultTTLSeconds = 600
	// garbageCollectionPeriodSeconds defines the period when garbage collection is triggered
	garbageCollectionPeriodSeconds = 5
)

// garbageCollector clear the state of dead records
type garbageCollector struct {
	state *map[key]map[question]state
}

func newGarbageCollector(state *map[key]map[question]state) *garbageCollector {
	return &garbageCollector{
		state: state,
	}
}

func (gc *garbageCollector) collect(){

	for k, qm := range *gc.state {
		for q, s := range qm {
			// remove death states for death questions
			if s.timestamp.Before(time.Now().Add(-garbageCollectionDefaultTTLSeconds * time.Second)){
				delete(qm, q)
			}
		}
		// remove subnet key, if contains 0 items
		if len(qm) == 0 {
			delete(*gc.state,k)
		}
	}
}
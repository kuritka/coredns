package strategy

import (
	"fmt"
	"testing"
	"time"
)

// po tom co updatnu record, tak musim zmenit timestamp => timestamp check!
// handlovat, kdyz udelam A a AAAA request
func TestStatefulGCCleaning(t *testing.T) {
	flattenTests := []stateFlatten{
		{"10.20.30.40", "test.example.com.", time.Now().Add(time.Hour * -5), []string{"10.10.10.10"}},
		{"10.20.30.40", "alpha.example.com.", time.Now().Add(time.Minute * -5), []string{"10.10.10.10", "20.20.20.20"}},
	}
	fstate := buildState(flattenTests)
	tests := []struct{
		name string
		ttlSeconds         int
		state *mstate
	}{
		{"clean on empty", 5, new(mstate)},
		{"clean all records", 5, &fstate},
		{"nil state", 5, nil},
	}
	for _ ,test := range tests {
		t.Run(test.name, func(t *testing.T) {
			newGarbageCollector(test.state, test.ttlSeconds).collect()
			if (test.state != nil) &&  len(*test.state) != 0 {
				t.Fatalf("Expected empty state but have %v records", len(*test.state))
			}
		})
	}
	// clean on empty
	// clean all records
}


func TestStatefulGCCleaningLive(t *testing.T) {
	// add records
	// remove records live
}


func TestStatefulGCRemoveItem(t *testing.T) {
	flattenTests := []stateFlatten {
		{"10.20.30.40", "test.example.com.", time.Now().Add(time.Hour*-5), []string{"10.10.10.10"}},
		{"10.20.30.40", "alpha.example.com.", time.Now().Add(time.Minute*-5), []string{"10.10.10.10", "20.20.20.20"}},
		{"10.20.30.40", "beta.example.com.", time.Now().Add(time.Second*-5), []string{}},
		{"10.20.30.40", "beta.example.com.", time.Now().Add(time.Second*-1), []string{"11.111.111.111", "222.222.222.333"}},
		{"11.11.11.11", "gc.test.com.", time.Now(), []string{"10.10.10.10"}},
	}

	tests := []struct{
		state  []stateFlatten
		ttlSeconds         int
		survivalRowIndexes map[int]bool
	}{
		{flattenTests, 3600*5+1, map[int]bool{0:true, 1:true, 2:true, 3:true,4:true}},
		{flattenTests, 60*5, map[int]bool{0:false, 1:false, 2:true, 3:true,4:true}},
		{flattenTests, 60*5+1, map[int]bool{0:false, 1:true, 2:true, 3:true,4:true}},
		{flattenTests, 2, map[int]bool{0:false, 1:false, 2:true, 3:true,4:true}},
		{flattenTests, 1, map[int]bool{0:false, 1:false, 2:false, 3:false,4:true}},
		{flattenTests, 0, map[int]bool{0:false, 1:false, 2:false, 3:false,4:false}},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("Delete records older than %v seconds", test.ttlSeconds), func(t *testing.T) {
			s := buildState(test.state)
			newGarbageCollector(&s, test.ttlSeconds).collect()

			for i, v := range flattenTests {
				// check if state for key x question exists
				exists := s.exists(key(v.key),question(v.question))

				if  test.survivalRowIndexes[i] != exists {
					t.Fatalf("Inconsistent state. Check if %v should be there or not ",v)
				}
			}
		})
	}
}

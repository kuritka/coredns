// Package roundrobin is responsible for roundrobin between DNS records. It contains several strategies
// - consistent (default) - all records are provided with the same probability
// - random - the random selection of records
// - weighted - records are provided according to a set weight.
package roundrobin

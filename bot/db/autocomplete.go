package db

import (
	"sort"
	"strings"

	"github.com/Nv7-Github/sevcord"
)

const maxResults = 25

func (d *Data) Autocomplete(query string) []sevcord.Choice {
	type searchResult struct {
		priority int
		uses     int
		id       string
		name     string
	}
	results := make([]searchResult, 0)
	d.RLock()
	for _, prog := range d.Programs {
		if strings.EqualFold(prog.Name, query) {
			results = append(results, searchResult{0, prog.Uses, prog.ID, prog.Name})
		} else if strings.HasPrefix(strings.ToLower(prog.Name), query) {
			results = append(results, searchResult{1, prog.Uses, prog.ID, prog.Name})
		} else if strings.Contains(strings.ToLower(prog.Name), query) {
			results = append(results, searchResult{2, prog.Uses, prog.ID, prog.Name})
		}
		if len(results) > 1000 {
			break
		}
	}
	d.RUnlock()

	// sort by uses
	sort.Slice(results, func(i, j int) bool {
		return results[i].uses > results[j].uses
	})

	// sort by priority
	sort.Slice(results, func(i, j int) bool {
		return results[i].priority < results[j].priority
	})

	// return top 25
	if len(results) > maxResults {
		results = results[:maxResults]
	}

	// sort by name
	sort.Slice(results, func(i, j int) bool {
		return results[i].name < results[j].name
	})

	// Return
	out := make([]sevcord.Choice, len(results))
	for i, res := range results {
		out[i] = sevcord.Choice{
			Name:  res.name,
			Value: res.id,
		}
	}
	return out
}

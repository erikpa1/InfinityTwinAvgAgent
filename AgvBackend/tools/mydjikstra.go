package tools

import (
	"turtle/core/lgr"

	"github.com/RyanCarrier/dijkstra/v2"
)

type PathResult struct {
	Path     []string
	Distance float64
	NotFound bool
}

func NewPathResult() PathResult {
	tmp := PathResult{}
	tmp.Path = make([]string, 0)
	return tmp
}

type Graph struct {
	djikstraGraph *dijkstra.Graph
	incIndex      int
	indexToUid    map[int]string
	uidToIndex    map[string]int
}

func NewGraph() *Graph {
	tmp := &Graph{}
	tmp.djikstraGraph = &dijkstra.Graph{}
	tmp.indexToUid = map[int]string{}
	tmp.uidToIndex = map[string]int{}
	tmp.incIndex = 0

	return tmp
}

func (self *Graph) AddEmptyVertex(uid string) {
	self.uidToIndex[uid] = self.incIndex
	self.indexToUid[self.incIndex] = uid
	self.djikstraGraph.AddEmptyVertex(self.incIndex)
	self.incIndex += 1
}

func (self *Graph) AddArc(a, b string, distance float64) {

	a_i, a_ok := self.uidToIndex[a]
	b_i, b_ok := self.uidToIndex[b]

	if a_ok && b_ok {
		//Obojstranny connection sa robi o uroven NAD!!!
		err1 := self.djikstraGraph.AddArc(a_i, b_i, uint64(distance*100))
		if err1 != nil {
			lgr.Error("Unable to find connection between", err1)
		}
	} else {
		if a_ok == false {
			lgr.Error("Unable to convert A", a)
		} else {
			lgr.Error("Unable to convert B", b)
		}

	}
}
func (self *Graph) Shortest(a, b string) PathResult {

	if a == b {
		lgr.Error("The same routes searching: [", a, "] , [", b, "]")
	}

	a_i, a_ok := self.uidToIndex[a]
	b_i, b_ok := self.uidToIndex[b]

	if a_ok && b_ok {
		shortest, err := self.djikstraGraph.Shortest(a_i, b_i)

		if err == nil {
			var uidPath = make([]string, len(shortest.Path))
			for i, val := range shortest.Path {
				uidPath[i] = self.indexToUid[val]
			}
			return PathResult{
				Distance: float64(shortest.Distance) / 100,
				Path:     uidPath,
			}
		} else {
			lgr.Error("Error between: [", a, "] or [", b, "]")
			lgr.Error(err.Error())
		}
	} else {
		lgr.Error("Unable to find: [", a, "] or [", b, "]")
	}

	return PathResult{}
}

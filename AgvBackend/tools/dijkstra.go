package tools

import (
	"github.com/RyanCarrier/dijkstra/v2"
	"turtle/core/lgr"
)

func TestDjikstra() {
	graph := dijkstra.NewGraph()
	//Add the 3 verticies
	graph.AddEmptyVertex(0)
	graph.AddEmptyVertex(1)
	graph.AddEmptyVertex(2)
	//Add the arcs
	graph.AddArc(0, 1, 1)
	graph.AddArc(0, 2, 1)
	graph.AddArc(1, 2, 2)

	best, err := graph.Shortest(0, 2)
	if err != nil {
		lgr.Error(err.Error())
	}
	lgr.Info("Shortest distance is", best.Distance, "following path ", best.Path)

	best, err = graph.Longest(0, 2)
	if err != nil {
		lgr.Error(err.Error())
	}
	lgr.Info("Longest distance is", best.Distance, "following path ", best.Path)

}

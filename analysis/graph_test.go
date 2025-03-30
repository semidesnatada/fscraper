package analysis

import (
	"testing"
)


func TestGraphCreationAndBFS_1(t *testing.T) {

	// initialise graph
	g := NewGraph()

	// initialise ids for nodes
	n1 := NewID()
	n2 := NewID()
	n3 := NewID()
	n4 := NewID()
	n5 := NewID()

	// add nodes and adjacencies to graph
	g.AddNode(n1, AdjacencyMap{
			n2: 1,
			n3: 5,
	})
	g.AddNode(n2, AdjacencyMap{
			n1: 1,
			n4: 7,
	})
	g.AddNode(n3, AdjacencyMap{
			n1: 5,
			n4: 15,
			n5: 3,
	})
	g.AddNode(n4, AdjacencyMap{
			n3: 15,
			n2: 7,
	})
	g.AddNode(n5, AdjacencyMap{
			n3:3,
	})

	// find the shortest path
	expectedPath := IDList{n5, n3, n1}
	actualPath, pErr := g.BfsShortestPath(n1, n5)
	if pErr != nil {
		t.Errorf("bfs algorithm error %v", pErr)
	}
	if len(actualPath) != len(expectedPath) {
		t.Errorf("length of identified path is not the same as the actual path. actual: %v, expected: %v", actualPath, expectedPath)
	}
	for i := range actualPath {
		step := actualPath[i]
		expectedStep := expectedPath[i]
		if step != expectedStep {
			t.Errorf("didn't find the correct steps in the path. actual: %v, expected: %v", actualPath, expectedPath)
		}
	}
}

func TestGraphCreationAndBFS_2(t *testing.T) {

	// initialise graph
	g := NewGraph()

	// initialise ids for nodes
	n1 := NewID()
	n2 := NewID()
	n3 := NewID()
	n4 := NewID()
	n5 := NewID()
	n6 := NewID()
	n7 := NewID()
	n8 := NewID()

	// add nodes and adjacencies to graph
	g.AddNode(n1, AdjacencyMap{
		n2: 1,
	})
	g.AddNode(n2, AdjacencyMap{
		n1: 1,
		n3: 7,
	})
	g.AddNode(n3, AdjacencyMap{
		n2: 7,
		n4: 15,
	})
	g.AddNode(n4, AdjacencyMap{
		n3: 15,
		n5: 7,
		n6: 10,
	})
	g.AddNode(n5, AdjacencyMap{
		n4:7,
		n7:8,
	})
	g.AddNode(n6, AdjacencyMap{
		n4:10,
	})
	g.AddNode(n7, AdjacencyMap{
		n5:8,
		n8:2,
	})
	g.AddNode(n8, AdjacencyMap{
		n7:2,
	})

	// find the shortest path
	expectedPath := IDList{n8, n7, n5, n4, n3, n2, n1}
	actualPath, pErr := g.BfsShortestPath(n1, n8)
	if pErr != nil {
		t.Errorf("bfs algorithm error %v", pErr)
	}
	if len(actualPath) != len(expectedPath) {
		t.Errorf("length of identified path is not the same as the actual path. actual: %v, expected: %v", actualPath, expectedPath)
	}
	for i := range actualPath {
		step := actualPath[i]
		expectedStep := expectedPath[i]
		if step != expectedStep {
			t.Errorf("didn't find the correct steps in the path. actual: %v, expected: %v", actualPath, expectedPath)
		}
	}
}
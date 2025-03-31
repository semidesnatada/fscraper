package analysis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"
	"sort"
	"sync"

	"github.com/google/uuid"
	"github.com/semidesnatada/fscraper/config"
)

type PlayerAdjacencyGraph struct {
	Nodes map[PlayerID]AdjacencyMap 
	Mu *sync.Mutex
}
type IDAdjacencyPair struct {
	ID PlayerID
	Adjacencies AdjacencyMap
}
type PlayerID uuid.UUID
type AdjacencyMap map[PlayerID]int
type IDList []PlayerID
type Path map[PlayerID]PlayerID

func NewID() PlayerID {
	return PlayerID(uuid.New())
}

func NewGraph() PlayerAdjacencyGraph {
	return PlayerAdjacencyGraph{Nodes: make(map[PlayerID]AdjacencyMap), Mu: &sync.Mutex{}}
}

func NewGraphFromDB(s *config.State) (PlayerAdjacencyGraph, error) {
	// initialise graph
	g := NewGraph()

	// get the list of players we want to extract records from the database about
	x, idErr := s.DB.GetPlayerUUIDsOrderedByUrl(context.Background())
	if idErr != nil {
		return PlayerAdjacencyGraph{}, idErr
	}
	ids := []PlayerID{}
	for _, id := range x {
		ids = append(ids, PlayerID(id))
	}

	fmt.Println("got ids from db")
	fmt.Println("number of ids is:", len(ids))

	// initialise a channel to receive new records as and when they come from the db
	ch := make(chan IDAdjacencyPair, len(ids))

	fmt.Println("iterating over ids to get adjacencies from db")
	// call goroutines which return the database entry to the channel
	for i, id := range ids {
		// end := i //+ 1 == len(ids)
		g.GetAdjacenciesFromDB(s, ch, id)
		if i % 100 == 0 {
			fmt.Printf("getting player number %d from database\n", i)
		}
		// fmt.Printf("going tot get id: %v from db\n", id)
	}

	// as records are piped into the channel, add them to the graph
	var count int
	for record := range ch {
		g.AddNode(record.ID, record.Adjacencies)
		count ++
		if count == len(ids) {close(ch)}
		// fmt.Printf("added node with id: %v to graph\n", record.ID)
	}
	fmt.Println("finished adding nodes to graph")

	return g, nil
}

func (g PlayerAdjacencyGraph) GetAdjacenciesFromDB(s *config.State, ch chan IDAdjacencyPair, playerID PlayerID) error {
	queryResponse := IDAdjacencyPair{ID: playerID}

	sharedRows, dbErr := s.DB.GetAllPlayersAndSharedMinsByID(context.Background(), uuid.UUID(queryResponse.ID))
	if dbErr != nil {
		ch <- IDAdjacencyPair{}
		fmt.Println(dbErr.Error())
		return dbErr
	}

	adjacencyMap := make(AdjacencyMap)

	for _, row := range sharedRows {
		adjacencyMap[PlayerID(row.OtherPlayerID)] = int(row.SharedMinutes)
	}

	queryResponse.Adjacencies = adjacencyMap

	ch <- queryResponse
	// if last {close(ch)}
	return nil
}

func (g PlayerAdjacencyGraph) AddNode(id PlayerID, adjacencies AdjacencyMap) error {
	// adds a node to the graph, if there isn't already a node with that ID in the graph.
	g.Mu.Lock()
	defer g.Mu.Unlock()
	_, ok := g.Nodes[id]
	if ok {
		return errors.New("can't add this node to the graph again")
	}
	g.Nodes[id] = adjacencies
	return nil
}

func (g PlayerAdjacencyGraph) RemoveNode(id PlayerID) error {
	// note: this does not remove references to this node from other node adjacency lists, so use this with caution.
	g.Mu.Lock()
	defer g.Mu.Unlock()
	if _, ok := g.Nodes[id]; ok {
		delete(g.Nodes, id)
		return nil
	 } else {
		return errors.New("this node can't be deleted as it is not in the graph")
	 }
}

func (g PlayerAdjacencyGraph) GetSortedListOfNeighbours(id PlayerID) IDList {
	// returns a list of player ids which are sorted by the number of minutes
	// shared by the players
	// this should not be called while undertaking write/delete operations on the graph.
	neighbours := g.Nodes[id]

    sortedIDs := make(IDList, 0, len(neighbours))
    for key := range neighbours {
        sortedIDs = append(sortedIDs, key)
    }
    sort.Slice(sortedIDs, func(i, j int) bool { return neighbours[sortedIDs[i]] > neighbours[sortedIDs[j]] })

	return sortedIDs
}

func (s *IDList) Pop() PlayerID {
	// removes the first item of the ID list.
	// returns the item that is removed and the original list is modified in place
	rv := (*s)[0]
	if len(*s) > 1 {
		*s = (*s)[1:]
	} else {
		*s = (*s)[:0]
	}
	return rv
}

func (g PlayerAdjacencyGraph) BfsShortestPath(start, end PlayerID) (IDList, error) {

	visited := IDList{}
	toVisit := IDList{start}
	pathMap := make(Path)

	for l := len(toVisit); l > 0 ; {
		current := toVisit.Pop()
		visited = append(visited, current)

		// if we have found the endpoint, we need to simply reconstruct
		//  the path of how we got here and return it
		if current == end {
			shortestPath := IDList{}
			for c := current; c != start; {
				shortestPath = append(shortestPath, c)
				c = pathMap[c]
			}
			shortestPath = append(shortestPath, start)
			return shortestPath, nil
		}

		//if we haven't found the endpoint, then find all of the neighbours and add them to the list
		//of nodes to traverse
		neighbs := g.GetSortedListOfNeighbours(current)
		for _, neighb := range neighbs {
			if !slices.Contains(visited, neighb) && !slices.Contains(toVisit, neighb) {
				toVisit = append(toVisit, neighb)
				pathMap[neighb] = current
			}

		}
		l = len(toVisit)
	}

	return IDList{}, errors.New("couldn't find a path between the two players")
}

func (g PlayerAdjacencyGraph) GetShortestConnectionBetweenPlayerUrls(s *config.State, playerUrl1, playerUrl2 string) ([]string, []int, error) {

	p1ID, err1 := s.DB.GetPlayerIdFromUrl(context.Background(), playerUrl1)
	if err1 != nil {
		return []string{}, []int{}, err1
	}
	p2ID, err2 := s.DB.GetPlayerIdFromUrl(context.Background(), playerUrl2)
	if err2 != nil {
		return []string{}, []int{}, err2
	}

	ids, err3 := g.BfsShortestPath(PlayerID(p1ID), PlayerID(p2ID))
	if err3 != nil {
		return []string{}, []int{}, err3
	}

	var outputs []string
	var outputi []int
	var prevId PlayerID

	for i, id := range ids {
		url, err4 := s.DB.GetPlayerUrlFromId(context.Background(), uuid.UUID(id))
		if err4 != nil {
			return []string{}, []int{}, err4
		}
		outputs = append(outputs, url)
		if i > 0 {
			outputi = append(outputi, g.Nodes[prevId][id])
		} else {
			outputi = append(outputi, 0)
		}
		prevId = id
	}
	return outputs, outputi, nil
}

func (g PlayerAdjacencyGraph) GetPathsBelowGivenLength(startUrl string, searchDepth int, s *config.State) (map[int][][]string, error) {

	startu, err := s.DB.GetPlayerIdFromUrl(context.Background(), startUrl)
	if err != nil {
		return map[int][][]string{}, err
	}
	start := PlayerID(startu)

	interMap, err2 := g.BfsGetAllValidPathsBelowGivenSize(start, searchDepth, s)
	if err2 != nil {
		return map[int][][]string{}, err2
	}

	out := make(map[int][][]string)
	for level, IDList := range interMap {
		urlsList := [][]string{}
		for _, idPath := range IDList {
			urls := []string{}
			for _, id := range idPath {
				url, uErr := s.DB.GetPlayerUrlFromId(context.Background(), uuid.UUID(id))
				if uErr != nil {
					return map[int][][]string{}, uErr
				}
				urls = append(urls, url)
		}
		urlsList = append(urlsList, urls)
		}
		out[level] = urlsList
	}

	return out, nil
}

func (g PlayerAdjacencyGraph) BfsGetAllValidPathsBelowGivenSize(start PlayerID, searchDepth int, s *config.State) (map[int][]IDList, error) {

	paths := map[int][]IDList{}

	visited := IDList{}
	// toVisit := IDList{start}
	pathMap := make(Path)

	nodesAtPrevDepth := IDList{start}
	// toVisitAtThisDepth := IDList{}

	for i := 1; i <= searchDepth ; i++ {
		// at depth 1, add a path back to the start for each of the neighbours of the start node to the map at key 1
		// at depth 2, loop through the neighbours of the start node and add
		toVisitAtThisDepth := IDList{}
		pathsAtThisDepth := []IDList{}

		for _, node := range nodesAtPrevDepth {
			neighbs := g.GetSortedListOfNeighbours(node)
			for _, neigh := range neighbs {
				if !slices.Contains(visited, neigh) {
					toVisitAtThisDepth = append(toVisitAtThisDepth, neigh)
					visited = append(visited, neigh)
					pathMap[neigh] = node
				}
			}
		}

		for _, node := range toVisitAtThisDepth {
			path := IDList{}
			for c := node; c != start; {
				path = append(path, c)
				c = pathMap[c]
			}
			path = append(path, start)
			pathsAtThisDepth = append(pathsAtThisDepth, path)
		}
		nodesAtPrevDepth = toVisitAtThisDepth
		paths[i] = pathsAtThisDepth
	}

	return paths, nil
}

func (g PlayerAdjacencyGraph) Write(filename string) error {
	dat, err1 := json.Marshal(g)
	if err1 != nil {
		return err1
	}
	dir, _ := os.Getwd()
	err2 := os.WriteFile(dir+"/graph_dumps/"+filename+".json", dat, 0666)
	if err2 != nil {
		fmt.Println("failed to write erroneous record to json file")
		fmt.Println(err2.Error())
	}
	return nil
}

func ReadGraphFromFile(filename string) (PlayerAdjacencyGraph, error) {
	
	dir, _ := os.Getwd()
	// fmt.Println("getting graph stored here:")
	// fmt.Println(dir+"/graph_dumps/"+filename+".json")

	data, err := os.ReadFile(dir+"/graph_dumps/"+filename+".json")
	if err != nil {
		return PlayerAdjacencyGraph{}, err
	}

	var g PlayerAdjacencyGraph
	err2 := g.UnmarshalJSON(data)
	if err2 != nil {
		return PlayerAdjacencyGraph{}, err2
	}

	return g, nil

}

// func (a AdjacencyMap) MarshalJSON() ([]byte, error) {

// 	bytes := []byte{}

// 	return bytes, nil
// }

// func (p PlayerID) MarshalJSON() ([]byte, error) {

// 	bytes := []byte{}

// 	return bytes, nil
// }

// func (u *PlayerID) MarshalJSON() ([]byte, error) {
//     return []byte(fmt.Sprintf("\"%s\"", uuid.UUID(*u).String())), nil
// }
// func (u *PlayerID) UnmarshalJSON(b []byte) error {
//     id, err := uuid.Parse(string(b[:]))
//     if err != nil {
//             return err
//     }
//     *u = PlayerID(id)
//     return nil
// }

// functions for marshalling the graph to json
func (u *PlayerID) ProcessEncode() (string, error) {
	return fmt.Sprintf("\"%s\"", uuid.UUID(*u).String()), nil
}

func ReversePlayerIDProcessEncode(s string) (PlayerID, error) {
	id, err := uuid.Parse(s)
    if err != nil {
            return PlayerID{}, err
    }
    return PlayerID(id), nil
}

func (a AdjacencyMap) ProcessEncode() (map[string]int, error) {
	out := make(map[string]int)
	for key, val := range a {
		processedKey, pErr := key.ProcessEncode()
		if pErr != nil {
			return map[string]int{}, pErr
		}
		out[processedKey] = val
	}
	return out, nil
}

func ReverseAdjacencyMapProcessEncode(m map[string]int) (AdjacencyMap, error) {

	out := make(AdjacencyMap)
	for key, val := range m {
		reversedKey, err := ReversePlayerIDProcessEncode(key)
		if err != nil {
			return AdjacencyMap{}, err
		}
		out[reversedKey] = val
	}

	return out, nil
}

func (g PlayerAdjacencyGraph) MarshalJSON() ([]byte, error) {
	encoded := make(map[string]map[string]int)

	for nodeID, adjacencyMap := range g.Nodes {
		encodedNodeID, enErr := nodeID.ProcessEncode()
		if enErr != nil {
			return []byte{}, enErr
		}
		encodedAdMap, eaErr := adjacencyMap.ProcessEncode()
		if eaErr != nil {
			return []byte{}, eaErr
		}
		encoded[encodedNodeID] = encodedAdMap
	}

	out, mErr := json.Marshal(encoded)
	if mErr != nil {
		return []byte{}, mErr
	}

	return out, nil
}

func (g *PlayerAdjacencyGraph) UnmarshalJSON(b []byte) error {

	unmarshalled := make(map[string]map[string]int)

	err := json.Unmarshal(b, &unmarshalled)
	if err != nil {
		return err
	}

	out := make(map[PlayerID]AdjacencyMap)

	for stringedPlayerID, processedAdjacencyMap := range unmarshalled {

		playerID, rErr := ReversePlayerIDProcessEncode(stringedPlayerID)
		if rErr != nil {
			return rErr
		}

		adMap, adErr := ReverseAdjacencyMapProcessEncode(processedAdjacencyMap)
		if adErr != nil {
			return adErr
		}
		out[playerID] = adMap
	}

	*g = PlayerAdjacencyGraph{Nodes: out, Mu: &sync.Mutex{}}

	return nil
}
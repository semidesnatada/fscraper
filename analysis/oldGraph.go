package analysis

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/semidesnatada/fscraper/config"
)

type Graph struct {
	numVertices int
	weightMatrix [][]int
	nodeIdentifiers map[int]uuid.UUID
	reverseNodeIdentifiers map[uuid.UUID]int
}

func NewPlayerConnectionGraph(s *config.State) (*Graph, error) {

	fmt.Println("starting creation of graph")

	ids, idErr := s.DB.GetPlayerUUIDsOrderedByUrl(context.Background())
	if idErr != nil {
		return &Graph{}, idErr
	}
	numVertices := len(ids)
	fmt.Println("successfully queried database for UUIDs")
	fmt.Printf("The graph will have this many vertices: %d\n",numVertices)

	nodeIdentifiers := make(map[int]uuid.UUID, numVertices)
	reverseNodeIdentifiers := make(map[uuid.UUID]int, numVertices)
	for ind, id := range ids {
		nodeIdentifiers[ind] = id
		reverseNodeIdentifiers[id] = ind
	}
	// fmt.Println(reverseNodeIdentifiers)
	// fmt.Println("created node-id identifier maps")
	// checkingId, _ := uuid.Parse("57aad65e-2532-4b3b-a6be-9cd3f7ac1703")
	// checkingId2, _ := uuid.Parse("5b033e03-6f3d-4878-871b-d2f81f19715f")
	// fmt.Println("checky chekcy", reverseNodeIdentifiers[checkingId])
	// fmt.Println("checky chekcy2", reverseNodeIdentifiers[checkingId2])

	weightMatrix := make([][]int, numVertices)
	for i := range weightMatrix {
		weightMatrix[i] = make([]int, numVertices)
		// now query the database to get the values required
		// fmt.Println(nodeIdentifiers[i])
		sharedMins, minErr := s.DB.GetAllPlayersAndSharedMinsByID(context.Background(),nodeIdentifiers[i])
		if minErr != nil {
			return &Graph{}, minErr
		}
		// if i == 0 || i == 3249 || i == 2221 {
		// 	fmt.Println()
		// 	fmt.Println(i)
		// 	for _, item := range sharedMins {
		// 		if int(item.SharedMinutes) > 1000 {
		// 			fmt.Println(item.OtherPlayerID)
		// 		}
		// 	}
		// }

		for _, item := range sharedMins {
			otherID := item.OtherPlayerID
			secondInd, ok := reverseNodeIdentifiers[otherID]
			if !ok {
				continue
			}
			mins := item.SharedMinutes
			// fmt.Println("player id", nodeIdentifiers[i], "player index", i, "other player id", otherID, "other player index", secondInd, "shared mins", mins)
			weightMatrix[i][secondInd] = int(mins)
		}
		fmt.Printf("finished creating the %d th/st row of the matrix\n", i)

	}

	fmt.Println("completed initialising the graph")

	return &Graph{
		numVertices: numVertices,
		weightMatrix: weightMatrix,
		nodeIdentifiers: nodeIdentifiers,
		reverseNodeIdentifiers: reverseNodeIdentifiers,
	}, nil
}


func (g *Graph) CheckPlayerCombo(s *config.State, url1, url2 string) (int, error) {

	id1, err1 := s.DB.GetPlayerIdFromUrl(context.Background(), url1)
	if err1 != nil {
		return 0, err1
	}

	id2, err2 := s.DB.GetPlayerIdFromUrl(context.Background(), url2)
	if err2 != nil {
		return 0, err2
	}

	fmt.Printf("got ids successfully\n id1: %v \n id2: %v \n", id1, id2)

	ind1, ok1 := g.reverseNodeIdentifiers[id1]
	ind2, ok2 := g.reverseNodeIdentifiers[id2]
	if !ok1 || !ok2 {
		return 0, errors.New("couldn't find these players in the graph")
	}
	fmt.Printf("both identifiers successfuly got: no1: %d no2: %d \n", ind1, ind2)

	min1 := g.weightMatrix[ind1][ind2]
	min2 := g.weightMatrix[ind2][ind1]

	fmt.Printf("got mins no1: %d no2: %d\n", min1, min2)

	if min1 != min2 {
		return 0, errors.New("issue with graph - numbers don't add up")
	}

	return min1, nil

}
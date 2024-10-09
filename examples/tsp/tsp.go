package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	antcolony "vaktibabat/ant_colony"
)

type TravelingSalesman struct {
	graph   antcolony.Graph
	weights [][]float64
}

// Used when computing the pheromones for ACO: the pheromones are set to the repricorial of the length of a
// hamilitonian cycle found with a greedy nearest-neighbour search
func (tsp TravelingSalesman) greedySolution() float64 {
	tour := make([]antcolony.Edge, 0)
	initComponent := uint(rand.Intn(len(tsp.graph.Nodes)))
	currComponent := initComponent
	memory := make(map[uint]bool)
	tourCost := 0.0

	// Our tour should be as long as the number of vertices
	for len(tour) != len(tsp.graph.Nodes) {
		var bestEdge antcolony.Edge
		bestWeight := math.Inf(1)

		for _, edge := range tsp.graph.Edges[currComponent] {
			memory[currComponent] = true

			if !memory[edge.B] && edge.A != edge.B {
				if tsp.weights[edge.A][edge.B] < bestWeight {
					bestEdge = edge
					bestWeight = tsp.weights[edge.A][edge.B]
				}
				bestEdge = edge
			}
		}

		// Go through the edge and change our current location
		currComponent = bestEdge.B
		tourCost += tsp.weights[bestEdge.A][bestEdge.B]
		tour = append(tour, bestEdge)
		// If we only have one edge left, we mark the initial location (the start of the cycle)
		// as unvisited again
		if len(tour) == len(tsp.graph.Nodes)-1 {
			memory[uint(initComponent)] = false
		}
	}

	return tourCost
}

func (tsp *TravelingSalesman) ConstructGraph() antcolony.Graph {
	return tsp.graph
}

func (tsp *TravelingSalesman) InitPheromones(num_ants uint) [][]float64 {
	pheromones := make([][]float64, 0)

	for i := 0; i < len(tsp.graph.Nodes); i++ {
		pheromone := make([]float64, 0)

		for j := 0; j < len(tsp.graph.Nodes); j++ {
			pheromone = append(pheromone, float64(num_ants)/tsp.greedySolution())
		}

		pheromones = append(pheromones, pheromone)
	}

	return pheromones
}

func (tsp *TravelingSalesman) InitHeuristics() [][]float64 {
	heuristics := make([][]float64, 0)

	for i := 0; i < len(tsp.graph.Nodes); i++ {
		heuristic := make([]float64, 0)

		for j := 0; j < len(tsp.graph.Nodes); j++ {
			heuristic = append(heuristic, 1.0/(tsp.weights[i][j]+1e-8))
		}

		heuristics = append(heuristics, heuristic)
	}

	return heuristics
}

func newCompleteGraph(num_nodes uint) antcolony.Graph {
	nodes := make([]uint, 0)
	edges := make([][]antcolony.Edge, 0)

	for i := 0; i < int(num_nodes); i++ {
		nodes = append(nodes, uint(i))
		curr_edges := make([]antcolony.Edge, 0)

		for j := 0; j < int(num_nodes); j++ {
			curr_edges = append(curr_edges, antcolony.Edge{A: uint(i), B: uint(j)})
		}

		edges = append(edges, curr_edges)
	}

	return antcolony.Graph{Nodes: nodes, Edges: edges}
}

func randomWeights(num_nodes uint) [][]float64 {
	weights := make([][]float64, 0)

	for i := 0; i < int(num_nodes); i++ {
		curr_weights := make([]float64, 0)

		for j := 0; j < int(num_nodes); j++ {
			curr_weights = append(curr_weights, 0.0)
		}

		weights = append(weights, curr_weights)
	}

	for i := 0; i < int(num_nodes); i++ {
		for j := 0; j < i; j++ {
			w := rand.Float64()
			weights[i][j] = w
			weights[j][i] = w
		}
	}

	return weights
}

func weightsFromFile(path string) [][]float64 {
	weights := make([][]float64, 0)
	file, _ := os.Open(path)

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		curr_weights := make([]float64, 0)

		for _, weight := range strings.Split(line, " ") {
			weight_parsed, _ := strconv.ParseFloat(weight, 64)
			curr_weights = append(curr_weights, weight_parsed)
		}

		weights = append(weights, curr_weights)
	}

	return weights
}

func main() {
	graph := newCompleteGraph(20)
	weights := weightsFromFile("./dist_mat")

	tsp := TravelingSalesman{graph: graph, weights: weights}

	antColony := antcolony.NewAntColony(&tsp, 200)
	antColony.RunSimulation(100)

	cycle := antColony.GetSolution()

	for _, edge := range cycle {
		fmt.Printf("(%d, %d)\n", edge.A, edge.B)
	}
}

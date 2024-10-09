package antcolony

import (
	"math"
	"math/rand"
	"sort"
)

// Exp. decay rate for the pheromone
const rho = 0.5

// Pheromone weight
const alpha = 1.0

// Heuristic weight
const beta = 3.0

// Ant-Cycle Implementation

type ACOptimizable interface {
	// How to construct a graph from this problem?
	ConstructGraph() Graph
	// How should the pheromones be initialized? For example,
	// for TSP, a common heuristic is to initialize all pheromones as m / C^{nn}, where m
	// is the number of ants, and C^{nn} is
	// the length of a cycle constructed with a nearest neighbour (greedy) heuristic
	InitPheromones(num_ants uint) [][]float64
	// Similarily, how should the heuristics be initialized?
	InitHeuristics() [][]float64
}

type AntColony struct {
	// The construction graph G = (C, L) of the problem
	// C is the set of components (e.g. cities in TSP or items in KS)
	// and L is the set of connections (in TSP, for example, all pairs of cities are connected)
	constructionGraph Graph
	// Pheromones on connections - this is increased every time an ant steps on the edge
	Pheromones [][]float64
	// We can also have heuristic information on the arcs - for TSP, this is the repriocorial of the cost of the edge
	heuristics [][]float64
	// The ants
	ants     []Ant
	num_ants uint
}

// An individual ant
type Ant struct {
	// The index of the current component, i.e. the current vertex in the construction graph
	currComponent uint
	// Which components has this ant already visited?
	// Used to define constraints
	memory map[uint]bool
	// We also store the explicit edges to compute the pheromones
	tour []Edge
}

// Construct a new ant colony for an ACOptimizable problem with num_ants ants
func NewAntColony(problem ACOptimizable, num_ants uint) *AntColony {
	colony := new(AntColony)
	colony.constructionGraph = problem.ConstructGraph()
	colony.Pheromones = problem.InitPheromones(num_ants)
	colony.heuristics = problem.InitHeuristics()
	colony.num_ants = num_ants
	colony.ants = make([]Ant, 0)

	// Initialize all the ants
	for i := 0; i < int(num_ants); i++ {
		// Generate a random city
		rand_component := rand.Intn(len(colony.constructionGraph.Nodes))
		// Append the ant to the ant list
		ant_memory := make(map[uint]bool)
		//ant_memory[uint(rand_component)] = true
		colony.ants = append(colony.ants, Ant{uint(rand_component), ant_memory, make([]Edge, 0)})
	}

	return colony
}

func (colony *AntColony) RunSimulation(num_iters int) {
	for i := 0; i < num_iters; i++ {
		// Have each ant complete a cycle
		for i := 0; i < int(colony.num_ants); i++ {
			colony.ants[i].DoCycle(colony)
		}

		// Evaporate the pheromones to avoid converging on a suboptimal solution
		colony.EvaporatePheromones()
		// Update the pheromones from all the ants
		for _, ant := range colony.ants {
			ant.DepositPheromones(colony)
			// We want a clean slate for our ant in the next iteration
			ant.ResetSolution(colony)
		}
	}
}

func (colony *AntColony) GetSolution() []Edge {
	colony.ants[0].DoCycle(colony)

	return colony.ants[0].tour
}

func (colony *AntColony) EvaporatePheromones() {
	for i := 0; i < len(colony.constructionGraph.Nodes); i++ {
		for j := 0; j < len(colony.constructionGraph.Nodes); j++ {
			colony.Pheromones[i][j] *= (rho)
		}
	}
}

func (ant *Ant) DoCycle(colony *AntColony) {
	initLocation := ant.currComponent

	// Our tour should be as long as the number of vertices
	for len(ant.tour) != len(colony.constructionGraph.Nodes) {
		ant.memory[ant.currComponent] = true
		// What is the probability of going to each edge in our neighbourhood?
		// For simplicity, we also track the probabilities of nodes not in our neighbourhood (and set them to 0)
		weights := make(map[uint]float64)
		// We track the sum of the edge scores so that we can normalize by it
		// and convert it to a valid probability distribution
		denom := 0.0

		for _, edge := range colony.constructionGraph.Edges[ant.currComponent] {
			// ant.memory[ant.currComponent] = true

			if !ant.memory[edge.B] && edge.A != edge.B {
				// The score for this edge is affected by the current amount of pheromones on it
				// and its heuristic (e.g. in TSP the heuristic is inversely proportional to the weight of the edge)
				score := math.Pow(colony.Pheromones[edge.A][edge.B], alpha) * math.Pow(colony.heuristics[edge.A][edge.B], beta)
				weights[edge.B] = score
				denom += score
			} else {
				// If this edge either (1) goes from the current node to itself or (2) the node it goes to has
				// already been visited, set its probability to 0
				weights[edge.B] = 0
			}
		}

		// Normalize the scores to convert into a valid probability distribution
		for dest := range weights {
			weights[dest] /= denom
		}

		// Sample one of the edges according to the probability distribution
		dest := weightedSampling(weights)
		edge := Edge{A: ant.currComponent, B: uint(dest)}
		// Go through the edge and change our current location
		ant.currComponent = edge.B
		ant.tour = append(ant.tour, edge)
		// If we only have one edge left, we mark the initial location (the start of the cycle)
		// as unvisited again
		if len(ant.tour) == len(colony.constructionGraph.Nodes)-1 {
			ant.memory[initLocation] = false
		}
	}
}

func (ant *Ant) DepositPheromones(colony *AntColony) {
	tourCost := 0.0

	for _, edge := range ant.tour {
		tourCost += 1.0 / colony.heuristics[edge.A][edge.B]
	}

	for _, edge := range ant.tour {
		colony.Pheromones[edge.A][edge.B] += 1.0 / tourCost
	}
}

func (ant *Ant) ResetSolution(colony *AntColony) {
	ant.memory = make(map[uint]bool)
	ant.currComponent = uint(rand.Intn(len(colony.constructionGraph.Nodes)))
	ant.tour = make([]Edge, 0)
}

// Sample from a discrete distribution where the probability of sampling v_i is p_i: P(v_i) = p_i
func weightedSampling(weights map[uint]float64) int {
	// Generate a random number 0 <= x < 1
	x := rand.Float64()
	// Sort the map by probability
	type KeyValue struct {
		idx    uint
		weight float64
	}

	kvs := make([]KeyValue, 0, len(weights))

	for idx, weight := range weights {
		kvs = append(kvs, KeyValue{idx, weight})
	}

	sort.Slice(kvs, func(i, j int) bool { return kvs[i].weight > kvs[j].weight })
	// Track culminative probability
	culm := 0.0

	for _, kv := range kvs {
		if culm < x && x < culm+kv.weight {
			return int(kv.idx)
		}

		culm += kv.weight
	}

	return 0
}

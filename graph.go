package antcolony

// An edge (a, b) in an undirected graph G
type Edge struct {
	A uint
	B uint
}

// A graph G = (V, E)
type Graph struct {
	// The list of node indices V
	Nodes []uint
	// We store the edges in a slice: entry i in the slice is the list of all edges from vertex i
	Edges [][]Edge
}

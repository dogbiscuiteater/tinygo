package main

import (
	"strconv"
  "math"
	"syscall/js"
)

func main() {
}

type node struct {
	data int

	// For each node, which node it can most efficiently be reached from.
	// If a node can be reached from many nodes, cameFrom will eventually contain the
	// most efficient previous step.
	previous *node

	// For each node, the cost of getting from the start node to that node.
	gScore float64

	// For each node, the total cost of getting from the start node to the goal
	// by passing by that node. That value is partly known, partly heuristic.
	fScore float64

	// The heuristic cost
	hScore float64

	world *squareGridWorld

	index int

	costModifier float64
}

type squareGridWorld struct {
	nodes      []*node
	m          map[int]*node
	N          int
	neighbours map[int]map[int]int
}


func (g squareGridWorld) coords(i int) (int, int) {
	x := i % g.N
	y := i / g.N
	return x, y
}

type x struct {
  a int
}

//go:export add
func add(a, b int) int {
	return a + b
}

//go:export update
func update() {
	document := js.Global().Get("document")
	aStr := document.Call("getElementById", "a").Get("value").String()
	bStr := document.Call("getElementById", "b").Get("value").String()
	a, _ := strconv.Atoi(aStr)
	b, _ := strconv.Atoi(bStr)
	result := add(a, b)
  x1 := x{ a: 100 }
	g := NewGrid(100) 
  document.Call("getElementById", "result").Set("value", g.N + x1.a + result)
}

//go:export NewGrid
func NewGrid(n int) *squareGridWorld {
	nodes := make([]*node, n*n, n*n)
	g := &squareGridWorld{
		nodes: nodes, N: n,
		neighbours: make(map[int]map[int]int),
		m:          make(map[int]*node),
	}

	nodes[0] = &node{data: 0, gScore: 0, index: 0, world: g}
	g.m[0] = nodes[0]
	for i := 1; i < len(nodes); i++ {
		x1, y1 := g.coords(i)
		x2, y2 := g.coords(n)
		nodes[i] = &node{data: i, previous: nil, gScore: DistanceBetween(x1, y1, x2, y2), index: i, world: g}
		g.m[i] = nodes[i]
	}
	
  for i := range nodes {
		//g.apply(i, g.N, up, upright, right, downright, down, downleft, left, upleft)
    right(i, n)
	}
  g.updateHValues(0, n * n -1)
	return g
}

//go:export DistanceBetween
func DistanceBetween(x1, y1, x2, y2 int) float64 {
	dx, dy := x2-x1, y2-y1
	return math.Sqrt(float64(dx*dx +dy*dy))
}

//go:export neighbour 
func (g *squareGridWorld) neighbour(i, n int) {
		if g.neighbours[i] == nil {
			g.neighbours[i] = make(map[int]int)
		}
		g.neighbours[i][n] = 1
}

//go:export right
func right(i, n int) {
	resolve(i, i+1, n)
}

//go:export left
func left(i, n int) {
	resolve(i, i-1, n)
}

//go:export down
func down(i, n int) {
	resolve(i, i+n, n)
}

//go:export up
func up(i, n int) {
	resolve(i, i-n, n)
}

//go:export downright
func downright(i, n int) {
	resolve(i, i+n+1, n)
}

//go:export downleft
func downleft(i, n int) {
	resolve(i, i+n-1, n)
}

//go:export upright
func upright(i, n int) {
	resolve(i, i-n+1, n)
}

//go:export upleft
func upleft(i, n int) {
	resolve(i, i-n-1, n)
}

//go:export resolve
func resolve(i, i2, n int){
	if i2 >= (n*n) || i2 < 0 {return }
	if i%n == 0 && i2%n == n-1 {return 	}
  if i%n == n-1 && i2%n == 0 {return 	}
  //neighbour(i, n) 
}

//go:export updateHValues
func (g *squareGridWorld) updateHValues(i, j int) {
	for i < len(g.nodes) {
		x1, y1 := g.coords(i)
		x2, y2 := g.coords(j)
		g.nodes[i].hScore  = DistanceBetween(x1, y1, x2, y2)
		i++
	}
}

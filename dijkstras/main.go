package main

import (
	"container/heap"
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

const (
	mazeWidth  = 40
	mazeHeight = 20
)

type TerrainType int

const (
	Path TerrainType = iota
	Wall
	Water
	Sand
	Start
	End
	Treasure
)

type Position struct {
	row, col int
}

type Node struct {
	pos   Position
	cost  int
	prev  *Node
	index int
}

type PriorityQueue []*Node

func (pq PriorityQueue) Len() int           { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].cost < pq[j].cost }
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Node)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

type Maze struct {
	grid       [][]TerrainType
	start, end Position
	treasures  []Position
	difficulty int
}


func generateMaze(difficulty int) *Maze {
	maze := &Maze{
		grid:       make([][]TerrainType, mazeHeight),
		difficulty: difficulty,
	}

	for i := range maze.grid {
		maze.grid[i] = make([]TerrainType, mazeWidth)
		for j := range maze.grid[i] {
			r := rand.Float32()
			switch {
			case r < float32(difficulty)*0.1:
				maze.grid[i][j] = Wall
			case r < float32(difficulty)*0.15:
				maze.grid[i][j] = Water
			case r < float32(difficulty)*0.2:
				maze.grid[i][j] = Sand
			default:
				maze.grid[i][j] = Path
			}
		}
	}

	maze.start = Position{0, 0}
	maze.end = Position{mazeHeight - 1, mazeWidth - 1}
	maze.grid[maze.start.row][maze.start.col] = Start
	maze.grid[maze.end.row][maze.end.col] = End

	// Add treasures
	numTreasures := 3
	for i := 0; i < numTreasures; i++ {
		for {
			row, col := rand.Intn(mazeHeight), rand.Intn(mazeWidth)
			if maze.grid[row][col] == Path {
				maze.grid[row][col] = Treasure
				maze.treasures = append(maze.treasures, Position{row, col})
				break
			}
		}
	}

	return maze
}

func (m *Maze) getTerrainCost(t TerrainType) int {
	switch t {
	case Wall:
		return 1000
	case Water:
		return 5
	case Sand:
		return 3
	default:
		return 1
	}
}

func (m *Maze) getTerrainSymbol(t TerrainType) string {
	switch t {
	case Wall:
		return "██"
	case Water:
		return "≈≈"
	case Sand:
		return "░░"
	case Start:
		return "S "
	case End:
		return "E "
	default:
		return "  "
	}
}

func (m *Maze) getTerrainColor(t TerrainType) color.Attribute {
	switch t {
	case Wall:
		return color.FgWhite
	case Water:
		return color.FgBlue
	case Sand:
		return color.FgYellow
	case Start:
		return color.FgGreen
	case End:
		return color.FgRed
	default:
		return color.FgBlack
	}
}

func printMaze(maze *Maze, visited [][]bool, path []Position, explorer Position) {
	fmt.Print("\033[H\033[2J") // Clear screen
	for i, row := range maze.grid {
		for j, cell := range row {
			if !visited[i][j] && !(i == explorer.row && j == explorer.col) {
				fmt.Print("  ") // Fog of war
				continue
			}
			if explorer.row == i && explorer.col == j {
				color.Set(color.FgHiGreen)
				fmt.Print("@◯")
			} else if isInPath(path, Position{i, j}) {
				color.Set(color.FgMagenta)
				fmt.Print("●●")
			} else if visited[i][j] {
				color.Set(color.FgHiBlack)
				fmt.Print("··")
			} else {
				color.Set(maze.getTerrainColor(cell))
				fmt.Print(maze.getTerrainSymbol(cell))
			}
		}
		fmt.Println()
	}
	color.Unset()
}

func isInPath(path []Position, pos Position) bool {
	for _, p := range path {
		if p == pos {
			return true
		}
	}
	return false
}

func solveMaze(maze *Maze) {
	startNode := &Node{pos: maze.start, cost: 0}
	pq := make(PriorityQueue, 1)
	pq[0] = startNode
	heap.Init(&pq)

	visited := make([][]bool, mazeHeight)
	for i := range visited {
		visited[i] = make([]bool, mazeWidth)
	}

	cameFrom := make(map[Position]*Node)
	costSoFar := make(map[Position]int)
	cameFrom[maze.start] = nil
	costSoFar[maze.start] = 0

	stepCount := 0
	startTime := time.Now()
	treasuresFound := 0
	explorer := maze.start

	for pq.Len() > 0 {
		current := heap.Pop(&pq).(*Node)
		stepCount++

		if current.pos == maze.end && treasuresFound == len(maze.treasures) {
			path := reconstructPath(cameFrom, maze.end)
			printMaze(maze, visited, path, explorer)
			duration := time.Since(startTime)
			fmt.Printf("Path found in %d steps! Time: %v\n", stepCount, duration)
			return
		}

		visited[current.pos.row][current.pos.col] = true
		explorer = current.pos

		if maze.grid[current.pos.row][current.pos.col] == Treasure {
			treasuresFound++
			maze.grid[current.pos.row][current.pos.col] = Path // Remove collected treasure
		}

		if stepCount%5 == 0 {
			printMaze(maze, visited, reconstructPath(cameFrom, current.pos), explorer)
			time.Sleep(50 * time.Millisecond)
		}

		directions := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
		for _, dir := range directions {
			newRow, newCol := current.pos.row+dir[0], current.pos.col+dir[1]
			if newRow >= 0 && newRow < mazeHeight && newCol >= 0 && newCol < mazeWidth {
				newPos := Position{newRow, newCol}
				newCost := costSoFar[current.pos] + maze.getTerrainCost(maze.grid[newRow][newCol])

				if cost, exists := costSoFar[newPos]; !exists || newCost < cost {
					costSoFar[newPos] = newCost
					newNode := &Node{pos: newPos, cost: newCost, prev: current}
					heap.Push(&pq, newNode)
					cameFrom[newPos] = current
				}
			}
		}
	}

	fmt.Println("No path found!")
}

func reconstructPath(cameFrom map[Position]*Node, end Position) []Position {
	path := []Position{end}
	current := end
	for cameFrom[current] != nil {
		current = cameFrom[current].pos
		path = append([]Position{current}, path...)
	}
	return path
}

func showMenu() int {
	fmt.Println("Welcome to the Advanced Maze Solver!")
	fmt.Println("1. Start Easy Game")
	fmt.Println("2. Start Medium Game")
	fmt.Println("3. Start Hard Game")
	fmt.Println("4. Quit")
	fmt.Print("Enter your choice: ")

	var choice int
	fmt.Scanf("%d", &choice)
	return choice
}

func main() {
	rand.Seed(time.Now().UnixNano())

	for {
		choice := showMenu()
		var difficulty int

		switch choice {
		case 1:
			difficulty = 1
		case 2:
			difficulty = 2
		case 3:
			difficulty = 3
		case 4:
			fmt.Println("Thanks for playing!")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
			continue
		}

		maze := generateMaze(difficulty)
		printMaze(maze, make([][]bool, mazeHeight), nil, maze.start)
		fmt.Println("Solving maze... Press Enter to start.")
		fmt.Scanln() // Wait for user input
		solveMaze(maze)

		fmt.Println("Press Enter to continue...")
		fmt.Scanln() 
	}
}
/*
Objectives
This project is meant to make you code a digital version of an ant farm.

Create a program lem-in that will read from a file (describing the ants and the colony) given in the arguments.

Upon successfully finding the quickest path, lem-in will display the content of the file passed as argument and each move the ants make from room to room.

How does it work?

You make an ant farm with tunnels and rooms.
You place the ants on one side and look at how they find the exit.
You need to find the quickest way to get n ants across a colony (composed of rooms and tunnels).

At the beginning of the game, all the ants are in the room ##start. The goal is to bring them to the room ##end with as few moves as possible.
The shortest path is not necessarily the simplest.
Some colonies will have many rooms and many links, but no path between ##start and ##end.
Some will have rooms that link to themselves, sending your path-search spinning in circles. Some will have too many/too few ants, no ##start or ##end, duplicated rooms, links to unknown rooms, rooms with invalid coordinates and a variety of other invalid or poorly-formatted input. In those cases the program will return an error message ERROR: invalid data format. If you wish, you can elaborate a more specific error message (example: ERROR: invalid data format, invalid number of Ants or ERROR: invalid data format, no start room found).
You must display your results on the standard output in the following format :

number_of_ants
the_rooms
the_links

Lx-y Lz-w Lr-o ...
x, z, r represents the ants numbers (going from 1 to number_of_ants) and y, w, o represents the rooms names.

A room is defined by "name coord_x coord_y", and will usually look like "Room 1 2", "nameoftheroom 1 6", "4 6 7".

The links are defined by "name1-name2" and will usually look like "1-2", "2-5".

Here is an example of this in practice :

##start
1 23 3
2 16 7
#comment
3 16 3
4 16 5
5 9 3
6 1 5
7 4 8
##end
0 9 5
0-4
0-6
1-3
4-3
5-2
3-5
#another comment
4-2
2-1
7-6
7-2
7-4
6-5
Which corresponds to the following representation :

	       _________________
	      /                 \
	 ____[5]----[3]--[1]     |
	/            |    /      |

[6]---[0]----[4]  /       |

	\   ________/|  /        |
	 \ /        [2]/________/
	 [7]_________/

Instructions
You need to create tunnels and rooms.
A room will never start with the letter L or with # and must have no spaces.
You join the rooms together with as many tunnels as you need.
A tunnel joins only two rooms together never more than that.
A room can be linked to multiple rooms.
Two rooms can't have more than one tunnel connecting them.
Each room can only contain one ant at a time (except at ##start and ##end which can contain as many ants as necessary).
Each tunnel can only be used once per turn.
To be the first to arrive, ants will need to take the shortest path or paths. They will also need to avoid traffic jams as well as walking all over their fellow ants.
You will only display the ants that moved at each turn, and you can move each ant only once and through a tunnel (the room at the receiving end must be empty).
The rooms names will not necessarily be numbers, and in order.
Any unknown command will be ignored.
The program must handle errors carefully. In no way can it quit in an unexpected manner.
The coordinates of the rooms will always be int.
Your project must be written in Go.
The code must respect the good practices.
It is recommended to have test files for unit testing.
Allowed packages
Only the standard Go packages are allowed.
Usage
Example 1 :

$ go run . test0.txt
3
##start
1 23 3
2 16 7
3 16 3
4 16 5
5 9 3
6 1 5
7 4 8
##end
0 9 5
0-4
0-6
1-3
4-3
5-2
3-5
4-2
2-1
7-6
7-2
7-4
6-5

L1-3 L2-2
L1-4 L2-5 L3-3
L1-0 L2-6 L3-4
L2-0 L3-0
$
Example 2 :

$ go run . test1.txt
3
##start
0 1 0
##end
1 5 0
2 9 0
3 13 0
0-2
2-3
3-1

L1-2
L1-3 L2-2
L1-1 L2-3 L3-2
L2-1 L3-3
L3-1
$
Example 3 :

$ go run . test1.txt
3
2 5 0
##start
0 1 2
##end
1 9 2
3 5 4
0-2
0-3
2-1
3-1
2-3

L1-2 L2-3
L1-1 L2-1 L3-2
L3-1
$
Bonus
As a bonus you can create an ant farm visualizer that shows the ants moving trough the colony.

Here is an usage example : ./lem-in ant-farm.txt | ./visualizer

The coordinates of the room will be useful only here.

This project will help you learn about :

Algorithmic
Ways to receive data
Ways to output data
Manipulation of strings
Manipulation of structures
*/
package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type Room struct {
	name string
}

type Link struct {
	from, to string
}

var (
	numAnts   int
	rooms     map[string]Room
	links     []Link
	startRoom string
	endRoom   string
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . <input_file>")
		return
	}
	inputFile := os.Args[1]

	if err := ParseInput(inputFile); err != nil {
		fmt.Println("Ошибка чтения файла:", err)
		return
	}

	if startRoom == "" {
		fmt.Println("ERROR: invalid data format, no start room found")
		return
	}
	if endRoom == "" {
		fmt.Println("ERROR: invalid data format, no end room found")
		return
	}
	fmt.Println("startRoom", startRoom)
	fmt.Println("endRoom", endRoom)

	fmt.Println("Количество муравьев:", numAnts)
	fmt.Println("Комнаты:")
	for name := range rooms {
		fmt.Printf("  %s\n", name)
	}
	fmt.Println("Туннели:")
	for _, link := range links {
		fmt.Printf("  %s-%s\n", link.from, link.to)
	}

	// Вызываем функцию поиска кратчайших путей
	distances, _ := Dijkstra(startRoom, endRoom, rooms, links)

	// Перемещаем муравьев
	MoveAnts(numAnts, startRoom, endRoom, distances, links)
}

func ParseInput(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	rooms = make(map[string]Room)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if err := parseLine(line); err != nil {
			return err
		}
	}

	return scanner.Err()
}

var isStart, isEnd bool

func parseLine(line string) error {
	if strings.HasPrefix(line, "#") {
		if line == "##start" {
			isStart = true
		} else if line == "##end" {
			isEnd = true
		}
		return nil
	}

	parts := strings.Fields(line)
	if len(parts) == 1 && strings.Contains(parts[0], "-") {
		roomNames := strings.Split(parts[0], "-")
		if len(roomNames) != 2 {
			return fmt.Errorf("invalid link format: %s", parts[0])
		}
		if _, ok := rooms[roomNames[0]]; !ok {
			return fmt.Errorf("room %s doesn't exist", roomNames[0])
		}
		if _, ok := rooms[roomNames[1]]; !ok {
			return fmt.Errorf("room %s doesn't exist", roomNames[1])
		}
		links = append(links, Link{roomNames[0], roomNames[1]})
		return nil
	}

	switch {
	case len(parts) == 1:
		num, err := strconv.Atoi(parts[0])
		if err != nil {
			return fmt.Errorf("invalid number of ants: %s", parts[0])
		}
		numAnts = num

	case len(parts) == 3:
		name := parts[0]
		rooms[name] = Room{name}
		if isStart {
			startRoom = name
			isStart = false
		} else if isEnd {
			endRoom = name
			isEnd = false
		}

	default:
		return fmt.Errorf("invalid line format: %s", line)
	}

	return nil
}

// Dijkstra's algorithm implementation with early exit
func Dijkstra(start string, end string, rooms map[string]Room, links []Link) (map[string]int, bool) {
	// Initialize distances with infinity for all rooms except start
	distances := make(map[string]int)
	for name := range rooms {
		distances[name] = math.MaxInt32
	}
	distances[start] = 0

	// Initialize priority queue (min heap) for vertices
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	// Push start vertex to priority queue
	heap.Push(&pq, &Item{value: start, priority: 0})

	// Process vertices until priority queue is empty
	for pq.Len() > 0 {
		// Extract min distance vertex from priority queue
		current := heap.Pop(&pq).(*Item).value

		// If the current room is the end room, we have found the shortest path
		if current == end {
			return distances, true
		}

		// Iterate over all adjacent vertices
		for _, link := range links {
			if link.from == current {
				neighbor := link.to
				// Check if the neighbor exists in the distances map
				if _, ok := distances[neighbor]; !ok {
					continue
				}
				// Calculate new distance
				newDist := distances[current] + 1
				// If the new distance is shorter, update distances and priority queue
				if newDist < distances[neighbor] {
					distances[neighbor] = newDist
					// Update priority queue
					heap.Push(&pq, &Item{value: neighbor, priority: newDist})
				}
			}
		}
	}

	// If we reach here, there is no path from start to end
	return nil, false
}

type Item struct {
	value    string // The value of the item; arbitrary.
	priority int    // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// Ant represents an ant
type Ant struct {
	id    int
	path  []string
	index int
}

func MoveAnts(numAnts int, startRoom string, endRoom string, distances map[string]int, links []Link) {
	// Initialize ants
	turn := 0

	ants := make([]Ant, numAnts)
	for i := range ants {
		ants[i] = Ant{id: i + 1, path: []string{startRoom}, index: 0}
	}

	// Initialize room occupancy
	occupancy := make(map[string]int)
	occupancy[startRoom] = numAnts

	// Move ants until all reach the end
	antsReachedEnd := 0 // Track the number of ants that reached the end
	for antsReachedEnd < numAnts {
		moved := false
		for i, ant := range ants {
			if ant.path[ant.index] == endRoom {
				antsReachedEnd++
				continue
			}

			// Find next room in the path
			nextRoom := ""
			minDist := math.MaxInt32
			for _, link := range links {
				if link.from == ant.path[ant.index] && distances[link.to] < minDist && occupancy[link.to] == 0 {
					nextRoom = link.to
					minDist = distances[link.to]
				}
			}

			// Move ant to the next room
			if nextRoom != "" {
				occupancy[ant.path[ant.index]]--
				ant.path = append(ant.path, nextRoom)
				ant.index++
				occupancy[nextRoom]++
				ants[i] = ant
				moved = true
				// Print move
				turn++
				fmt.Printf("L%d-%s ", ant.id, nextRoom)
			}
		}

		if !moved {
			break // Exit loop if no ants moved
		}

		fmt.Println()
		if turn > 100000 { // Protection against infinite loop
			fmt.Println("Breaking to prevent a potential infinite loop.")
			break
		}
	}

	// Print total number of moves
	fmt.Println("Total number of moves:", len(ants[0].path)-1)
}

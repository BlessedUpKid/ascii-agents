package agents

import (
	"math"
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
)

type Agents struct {
	agents map[int]Agent
	x      int
	y      int
}

func Coord(x, y, width, height int) int {
	coord := x*width + y
	return coord
}

func (a *Agents) Add(x, y int, ag Agent) {
	a.agents[Coord(x, y, a.x, a.y)] = ag
}

func (a *Agents) Set(x, y int, ag Agent) {
	a.Remove(ag)
	ag.SetPos(x, y)
	a.Add(x, y, ag)
}

func (a *Agents) Remove(ag Agent) {
	x, y := ag.Pos()
	delete(a.agents, Coord(x, y, a.x, a.y))
}

func (a *Agents) FindAgent(x, y int) Agent {
	agent, ok := a.agents[Coord(x, y, a.x, a.y)]
	if !ok {
		return nil
	}

	return agent
}

func (a *Agents) Ns(x, y, d int) []Agent {
	var allCoords [][]int

	for i := 1; i <= d; i++ {
		allCoords = append(allCoords, Ncoords(x, y, i)...)
	}

	var neighbors []Agent
	for _, coord := range allCoords {
		i := coord[0]
		j := coord[1]

		u := wrapCoord(i, a.x)
		w := wrapCoord(j, a.y)

		nextAgent := a.FindAgent(u, w)
		if nextAgent != nil {
			neighbors = append(neighbors, nextAgent)
		}
	}

	return neighbors
}

func (a *Agents) List() []Agent {
	var agents []Agent
	for _, ag := range a.agents {
		agents = append(agents, ag)
	}
	return agents
}

func NewAgentsManager(x, y int) *Agents {
	return &Agents{
		agents: make(map[int]Agent, 0),
		x:      x,
		y:      y,
	}
}

type Agent interface {
	Update()
	Val() int
	Color() tcell.Color
	Pos() (int, int)
	SetPos(x, y int)
}

type agentA struct {
	agent
}

func NewAgentA(mgmt *Agents) *agentA {
	r := rand.New(rand.NewSource(time.Now().UnixMicro()))
	return &agentA{
		agent: agent{
			val:        1,
			randomizer: r,
			mgmt:       mgmt,
			dirX:       0,
			dirY:       0,
			p:          0,
			g:          32,
			x:          0,
			y:          0,
		},
	}
}

func (a *agentA) Update() {
	a.p += 0.1
	n := 501 + int(500*math.Sin(a.p))

	// sa.dirX = int(4*math.Sin(sa.p)) * (sa.randomizer.Intn(2) - 1)
	// sa.dirY = int(4*math.Sin(sa.p)) * (sa.randomizer.Intn(2) - 1)

	a.val = n
	ns := a.mgmt.Ns(a.x, a.y, 1)
	checked := len(ns)
	if checked > 2 {
		a.dirX = 0
		a.dirY = 0
	}

	lookX, lookY := a.agent.Lookhead()

	there := a.mgmt.FindAgent(lookX, lookY)
	if there == nil {
		a.mgmt.Set(lookX, lookY, a)
	}

	if n > 949 && n <= 1000 {
		a.Spwan()
	}

	if checked > 3 {
		if n < 100 && a.p > 4*math.Pi {
			a.mgmt.Remove(a)
		}
		return
	}

	if n < 100 && a.p > math.Pi {
		a.mgmt.Remove(a)
	}
}

func (a *agentA) Spwan() {
	spawn := NewAgentA(a.mgmt)
	coin := a.randomizer.Intn(4)
	var x, y int
	switch coin {
	case 0:
		x = -1
		y = -1
	case 1:
		x = 1
		y = 1
	case 2:
		x = 1
		y = -1
	case 3:
		x = -1
		y = 1
	default:
		x = 0
		y = 0
	}
	ox := x
	oy := y
	spawn.x = wrapCoord(a.x+ox, a.mgmt.x)
	spawn.y = wrapCoord(a.y+oy, a.mgmt.y)
	spawn.dirX = ox
	spawn.dirY = oy
	spawn.g = clampHigh(a.g+16, 256)

	occupied := a.mgmt.FindAgent(spawn.x, spawn.y)
	if occupied == nil {
		a.mgmt.Add(spawn.x, spawn.y, spawn)
	}
}

func (a *agentA) Val() int {
	return a.val
}

func (a *agentA) Pos() (int, int) {
	return a.x, a.y
}

func (a *agentA) Color() tcell.Color {
	g32 := int32(a.g)
	return tcell.NewRGBColor(g32, 255, g32)
}

func (a *agentA) SetPos(x, y int) {
	a.x = x
	a.y = y
}

type agentB struct {
	agent
}

func NewAgentB(mgmt *Agents) *agentB {
	r := rand.New(rand.NewSource(time.Now().UnixMicro()))
	return &agentB{
		agent{
			val:        1,
			randomizer: r,
			mgmt:       mgmt,
			dirX:       0,
			dirY:       0,
			g:          32,
			x:          0,
			y:          0,
		},
	}
}

func (a *agentB) Update() {
	a.p += 0.1
	n := 501 + int(500*math.Sin(a.p))

	a.val = n
	lookX, lookY := a.agent.Lookhead()

	there := a.mgmt.FindAgent(lookX, lookY)
	if there == nil {
		a.mgmt.Set(lookX, lookY, a)
	}

	if n > 900 && n < 1000 {
		a.Spwan(1)
		a.Spwan(2)
		a.Spwan(3)
		a.Spwan(0)
	}

	if n < 100 && a.p > 3*math.Pi {
		a.mgmt.Remove(a)
	}
}

func (a *agentB) Spwan(z int) {
	spawn := NewAgentB(a.mgmt)
	// coin := a.randomizer.Intn(4)
	coin := z
	var x, y int
	switch coin {
	case 0:
		x = -1
		y = 0
	case 1:
		x = 1
		y = 0
	case 2:
		x = 0
		y = 1
	case 3:
		x = 0
		y = -1
	default:
		x = 0
		y = 0
	}

	spawn.x = wrapCoord(a.x+x, a.mgmt.x)
	spawn.y = wrapCoord(a.y+y, a.mgmt.y)
	spawn.dirX = 0
	spawn.dirY = 0
	spawn.g = clampHigh(a.g+16, 256)

	occupied := a.mgmt.FindAgent(spawn.x, spawn.y)
	if occupied == nil {
		a.mgmt.Add(spawn.x, spawn.y, spawn)
	}
}

func (a *agentB) Val() int {
	return a.val
}

func (a *agentB) Pos() (int, int) {
	return a.x, a.y
}

func (a *agentB) Color() tcell.Color {
	g32 := int32(a.g)
	return tcell.NewRGBColor(g32, g32, 255)
}

func (a *agentB) SetPos(x, y int) {
	a.x = x
	a.y = y
}

type agentC struct {
	agent
}

func NewAgentC(mgmt *Agents) *agentC {
	r := rand.New(rand.NewSource(time.Now().UnixMicro()))
	return &agentC{
		agent{
			val:        0,
			randomizer: r,
			mgmt:       mgmt,
			dirX:       0,
			dirY:       0,
			g:          32,
			x:          0,
			y:          0,
		},
	}
}

// "Explodes" killing everything around and itself
func (a *agentC) Update() {
	a.p += 0.1
	n := 501 + int(500*math.Sin(a.p))

	a.val = n

	if n > 950 {
		ns := a.mgmt.Ns(a.x, a.y, 3)
		for _, n := range ns {
			a.mgmt.Remove(n)
		}

		a.mgmt.Remove(a)
	}
}

func (a *agentC) Val() int {
	return a.val
}

func (a *agentC) Pos() (int, int) {
	return a.x, a.y
}

func (a *agentC) Color() tcell.Color {
	g32 := int32(a.g)
	return tcell.NewRGBColor(255, g32, g32)
}

func (a *agentC) SetPos(x, y int) {
	a.x = x
	a.y = y
}

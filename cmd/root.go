package commands

import (
	"fmt"
	xlog "log"
	"os"
	"time"
	"viz/pkg/agents"
	"viz/pkg/render"
	"viz/pkg/types"
	"viz/pkg/ui"

	tcell "github.com/gdamore/tcell/v2"
	"github.com/spf13/cobra"
)

func init() {
}

func LogMessage(str string) {
	fmt.Printf("%s\n", str)
}

var rootCmd = &cobra.Command{
	Use:   "viz",
	Short: "viz is for some visual interest on stream.",
	Run: func(cmd *cobra.Command, args []string) {
		s, err := tcell.NewScreen()
		if err != nil {
			xlog.Fatalf("%+v", err)
		}
		if err := s.Init(); err != nil {
			xlog.Fatalf("%+v", err)
		}

		defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
		s.SetStyle(defStyle)

		quit := func() {
			// You have to catch panics in a defer, clean up, and
			// re-raise them - otherwise your application can
			// die without leaving any diagnostic trace.
			maybePanic := recover()
			s.Fini()
			if maybePanic != nil {
				panic(maybePanic)
			}
		}
		defer quit()

		var w, l int
		w, l = s.Size()
		mat := types.NewIntMatrix(w, l)
		mat2 := types.NewIntMatrix(10, 10)

		// Init randomizer
		// randomizer := rand.New(rand.NewSource(time.Now().UnixMicro()))
		startTime := time.Now()
		diff := time.Duration(0)

		evCh := make(chan tcell.Event)
		quitCh := make(chan struct{})

		amgmt := agents.NewAgentsManager(w, l)
		a2 := agents.NewAgentB(amgmt)
		a2.SetPos(w/2, l/2)
		amgmt.Add(w/2, l/2, a2)

		submgmt := agents.NewAgentsManager(10, 10)

		var updateEnabled bool = true
		var brush int = 4
		var tick time.Duration = 100
		s.EnableMouse()
		go s.ChannelEvents(evCh, quitCh)

	loop:
		for {
			// Process events aynsc
			select {
			case event := <-evCh:
				if e, ok := event.(*tcell.EventKey); ok {
					if e.Key() == tcell.KeyCtrlC || e.Key() == tcell.KeyESC {
						close(quitCh)
						break loop
					}

					if e.Key() == tcell.KeyCtrlA {
						updateEnabled = !updateEnabled
					}

					if e.Key() == tcell.KeyCtrlB {
						brush = (brush + 1) % 5
						for _, agent := range submgmt.List() {
							submgmt.Remove(agent)
						}

						a := agents.GetAgentFromBrush(submgmt, brush)
						a.SetPos(1, 1)
						submgmt.Add(1, 1, a)
					}

					if e.Key() == tcell.KeyCtrlD {
						for _, agent := range amgmt.List() {
							amgmt.Remove(agent)
						}
					}

					if e.Key() == tcell.KeyLeft {
						tick = tick * 2
					}

					if e.Key() == tcell.KeyRight {
						tick = (tick / 2) + 1
					}
				}

				if e, ok := event.(*tcell.EventMouse); ok {
					posx, posy := e.Position()
					btns := e.Buttons()

					switch btns {
					case tcell.Button1:

						agent := agents.GetAgentFromBrush(amgmt, brush)
						agent.SetPos(posx, posy)
						amgmt.Add(posx, posy, agent)

					}
				}

				if e, ok := event.(*tcell.EventResize); ok {
					w, l = e.Size()
				}

			default:
			}

			diff = time.Since(startTime)
			if diff > time.Millisecond*tick {
				diff = 0
				startTime = time.Now()
				s.Clear()
				mat.Clear()
				mat2.Clear()

				for _, agent := range amgmt.List() {
					if updateEnabled {
						agent.Update()
					}
					x, y := agent.Pos()
					mat.Set(x, y, agent.Val())
					mat.SetColor(x, y, agent.Color())
				}

				for _, agent := range submgmt.List() {
					if updateEnabled {
						agent.Update()
					}
					x, y := agent.Pos()
					mat2.Set(x, y, agent.Val())
					mat2.SetColor(x, y, agent.Color())
				}

				ui.DrawRect(mat2, 0, 0, mat2.Cols()-1, mat2.Rows()-1)

				w, l = s.Size()
				mat.Composite(mat2, w-10, l-10)

				render.RenderMat(s, mat)
				// render.RenderMat(s, mat2)
			}
			s.Show()
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

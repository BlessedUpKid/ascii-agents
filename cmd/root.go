package commands

import (
	"fmt"
	xlog "log"
	"os"
	"time"
	"viz/pkg/agents"
	"viz/pkg/render"
	"viz/pkg/types"

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

		// Init randomizer
		// randomizer := rand.New(rand.NewSource(time.Now().UnixMicro()))
		startTime := time.Now()
		diff := time.Duration(0)

		evCh := make(chan tcell.Event)
		quitCh := make(chan struct{})

		amgmt := agents.NewAgentsManager(w, l)
		a1 := agents.NewAgentA(amgmt)
		amgmt.Add(0, 0, a1)

		a2 := agents.NewAgentB(amgmt)

		a2.SetPos(w/2, l/2)
		amgmt.Add(w/2, l/2, a2)

		var updateEnabled bool = true
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
				}

				if e, ok := event.(*tcell.EventMouse); ok {
					posx, posy := e.Position()
					btns := e.Buttons()

					switch btns {
					case tcell.Button1:
						a := agents.NewAgentC(amgmt)
						a.SetPos(posx, posy)
						amgmt.Add(posx, posy, a)
					}
				}

				if e, ok := event.(*tcell.EventResize); ok {
					w, l = e.Size()
				}

			default:
			}

			diff = time.Since(startTime)
			if diff > time.Millisecond*100 {
				diff = 0
				startTime = time.Now()
				s.Clear()
				mat.Clear()

				for _, agent := range amgmt.List() {
					if updateEnabled {
						agent.Update()
					}
					x, y := agent.Pos()
					mat.Set(x, y, agent.Val())
					mat.SetColor(x, y, agent.Color())
				}
				render.RenderMat(s, mat)
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

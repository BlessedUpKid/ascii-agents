package commands

import (
	"fmt"
	"viz/pkg/agents"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(testCmd)
}

// func newAgentAtPos(m *agents.Agents, x, y int) agents.Agent {
// 	a := agents.NewAgentA(m)
// 	a.SetPos(x, y)
// 	m.Add(x, y, a)
// 	return a
// }

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Runs tests.",
	Run: func(cmd *cobra.Command, args []string) {
		coords := agents.Ncoords(0, 0, 2)

		for _, coord := range coords {
			fmt.Printf("%d %d\n", coord[0], coord[1])
		}

		fmt.Printf("len: %d", len(coords))
	},
}

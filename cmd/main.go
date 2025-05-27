package cmd

import (
	"fmt"
	"github.com/EduardMikhrin/territorial-fishermans/cmd/service"
	"github.com/spf13/cobra"
	"os"
)

func Execute() {
	root := &cobra.Command{
		Use:   "relayer-svc",
		Short: "Relayer Service is developed to broadcast userOPs to network",
	}

	root.AddCommand(service.Cmd)

	if err := root.Execute(); err != nil {
		fmt.Printf("Error occured: %v\n", err)
		os.Exit(1)
	}
}

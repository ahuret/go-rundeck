package cmds

import (
	"github.com/ahuret/go-rundeck/pkg/cli"
	"github.com/spf13/cobra"
)

func listSystemACLPoliciesFunc(cmd *cobra.Command, args []string) error {
	policies, err := cli.Client.ListSystemACLPolicies()
	if err != nil {
		return err
	}
	cli.OutputFormatter.SetHeaders([]string{
		"Name",
		"Path",
		"Type",
		"HRef",
		"Parent",
		"Parent Type",
	})
	parent := "/"
	if policies.Path != "" {
		parent = policies.Path
	}
	for _, p := range policies.Resources {
		if err := cli.OutputFormatter.AddRow([]string{
			p.Name,
			p.Path,
			p.Type,
			p.Href,
			parent,
			policies.Type,
		}); err != nil {
			return err
		}
	}
	cli.OutputFormatter.Draw()
	return nil
}

func listSystemACLPoliciesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "lists system acl policies",
		RunE:  listSystemACLPoliciesFunc,
	}
	rootCmd := cli.New(cmd)
	return rootCmd
}

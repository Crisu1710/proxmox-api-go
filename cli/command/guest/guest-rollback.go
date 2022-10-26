package guest

import (
	"fmt"

	"github.com/Telmate/proxmox-api-go/cli"
	"github.com/Telmate/proxmox-api-go/proxmox"
	"github.com/spf13/cobra"
)

var guest_rollbackCmd = &cobra.Command{
	Use:   "rollback GUESTID SNAPSHOT",
	Short: "Shuts the speciefid guest down",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		vmr := proxmox.NewVmRef(cli.ValidateIntIDset(args, "GuestID"))
		snapName := cli.RequiredIDset(args, 1, "SnapshotName")
		c := cli.NewClient()
		_, err = c.RollbackSnapshot(vmr, snapName)
		if err == nil {
			fmt.Fprintf(GuestCmd.OutOrStdout(), "Guest with id (%d) has been rolled back to snapshot (%s)\n", vmr.VmId(), snapName)
		}
		return
	},
}

func init() {
	GuestCmd.AddCommand(guest_rollbackCmd)
}

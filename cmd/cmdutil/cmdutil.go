// Package cmdutil holds small helpers shared by the cmd subpackages.
package cmdutil

import (
	"time"

	"github.com/UnipayFI/okx-cli/common"
	"github.com/spf13/cobra"
)

// ParseTime reads a time-valued string flag and parses it with the shared CLI
// rules (unix-ms or "YYYY-MM-DD HH:MM:SS"). An unset flag yields a zero time.
func ParseTime(cmd *cobra.Command, flag string) (time.Time, error) {
	raw, _ := cmd.Flags().GetString(flag)
	parsed, _, err := common.ParseTimeFlag("--"+flag, raw)
	if err != nil {
		return time.Time{}, err
	}
	return parsed, nil
}

package ignitecmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/chaincmd"
	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/ignite/pkg/xos"
	"github.com/ignite/cli/ignite/services/network/networkchain"
)

func newNetworkChainShowGenesis() *cobra.Command {
	c := &cobra.Command{
		Use:   "genesis [launch-id]",
		Short: "Show the chain genesis file",
		Args:  cobra.ExactArgs(1),
		RunE:  networkChainShowGenesisHandler,
	}

	flagSetClearCache(c)
	c.Flags().String(flagOut, "./genesis.json", "path to output Genesis file")

	return c
}

func networkChainShowGenesisHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New(cliui.StartSpinner())
	defer session.End()

	out, _ := cmd.Flags().GetString(flagOut)

	cacheStorage, err := newCache(cmd)
	if err != nil {
		return err
	}

	nb, launchID, err := networkChainLaunch(cmd, args, session)
	if err != nil {
		return err
	}
	n, err := nb.Network()
	if err != nil {
		return err
	}

	chainLaunch, err := n.ChainLaunch(cmd.Context(), launchID)
	if err != nil {
		return err
	}

	networkOptions := []networkchain.Option{
		networkchain.WithKeyringBackend(chaincmd.KeyringBackendTest),
	}

	c, err := nb.Chain(networkchain.SourceLaunch(chainLaunch), networkOptions...)
	if err != nil {
		return err
	}

	// generate the genesis in a temp dir
	tmpHome, err := os.MkdirTemp("", "*-spn")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpHome)

	c.SetHome(tmpHome)

	if err := prepareFromGenesisInformation(
		cmd,
		cacheStorage,
		launchID,
		n,
		c,
		chainLaunch,
	); err != nil {
		return err
	}

	// get the new genesis path
	genesisPath, err := c.GenesisPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(out), 0o744); err != nil {
		return err
	}

	if err := xos.Rename(genesisPath, out); err != nil {
		return err
	}

	return session.Printf("%s Genesis generated: %s\n", icons.Bullet, out)
}

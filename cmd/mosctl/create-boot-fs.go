package main

import (
	"github.com/apex/log"
	"github.com/project-machine/mos/pkg/mosconfig"
	"github.com/urfave/cli"
)

var createBootFsCmd = cli.Command{
	Name: "create-boot-fs",
	Usage: "Create a boot filesystem",
	Action: doCreateBootfs,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name: "config-dir, c",
			Usage: "Directory where mos config is found",
			Value: "/config",
		},
		cli.StringFlag{
			Name: "atomfs-store, a",
			Usage: "Directory under which atomfs store is kept",
			Value: "/atomfs-store",
		},
		cli.StringFlag{
			Name: "scratch-dir, s",
			Usage: "Directory under which storage should keep overlays and tempdirs",
			Value: "/scratch-writes",
		},
		cli.StringFlag{
			Name: "dest",
			Usage: "Directory over which to mount the rfs",
			Value: "/sysroot",
		},
	},
}

// Setup a rootfs to which dracut should pivot.
// Note, setup of luks keys, SUDI keys, and extension of PCR7
// must already have been done.
func doCreateBootfs(ctx *cli.Context) error {
	opts := mosconfig.DefaultMosOptions()
	opts.ConfigDir = ctx.String("config-dir")
	opts.StorageCache = ctx.String("atomfs-store")
	opts.ScratchWrites = ctx.String("scratch-dir")

	m, err := mosconfig.OpenMos(opts)
	if err != nil {
		return err
	}

	t, err := m.Current("hostfs")
	if err != nil {
		return err
	}

	dest := ctx.String("dest")
	//err = m.Storage.MountWriteable(t, dest)
	//if err != nil {
		//return err
	//}
	log.Infof("Would mount %s:%s to %s", t.Fullname, t.Version, dest)

	log.Infof("Rootfs has been setup under %s", dest)
	return nil
}

package mosconfig

import (
	"fmt"
	"os"
)

type MosOptions struct {
	// Storage type - atomfs for now
	StorageType StorageType

	// Where install manifest is found
	ConfigDir     string

	// Directory where atomfs/puzzlefs cache is rooted
	StorageCache  string

	// Directory where atomfs keeps its working storage
	// e.g. upperdirs and temporary mounts
	ScratchWrites  string

	// whether this will be a session which writes to some mos state
	LayersReadOnly bool

	// whether this will be a session which writes to the system manifest
	ManifestReadOnly bool

	// During initial install, we can't read the provisioned host certs
	NoHostCerts bool
}

func DefaultMosOptions() MosOptions {
	return MosOptions{
		StorageType:      AtomfsStorageType,
		ConfigDir:        "/config",
		StorageCache:     "/atomfs-store",
		ScratchWrites:    "/scratch-writes",
		LayersReadOnly:   true,
		ManifestReadOnly: true,
		NoHostCerts:      false,
	}
}

type Mos struct {
	//storage   Storage
	//bootmgr   Bootmgr

	opts        MosOptions
	lockfile    *os.File
}

func NewMos(configDir, storeDir string) (*Mos, error) {
	opts := MosOptions{
		StorageType: AtomfsStorageType,
		ConfigDir: configDir,
		StorageCache: storeDir,
		LayersReadOnly: false,
		ManifestReadOnly: false,
		NoHostCerts: true,
	}

	mos := &Mos{
		opts: opts,
		lockfile: nil,
	}
	err := mos.acquireLock()
	if err != nil {
		return mos, err
	}
	return mos, nil
}

func OpenMos(opts MosOptions) (*Mos, error) {
	mos := &Mos{opts: opts}

	err := mos.acquireLock()
	if err != nil {
		return nil, err
	}
	return mos, nil
}

func (mos *Mos) Close() {
	if mos.lockfile != nil {
		mos.lockfile.Close()
		mos.lockfile = nil
	}
}

// Give the current information for the target named @name, for instance
// 'hostfs' or 'zot.  Returns a *Target containing the full target
// information from the manifest
func (mos *Mos) Current(name string) (*Target, error) {
	manifest, err := mos.CurrentManifest()
	if err != nil {
		return nil, err
	}

	for _, t := range manifest.Targets {
		if t.Name == name {
			return &t, nil
		}
	}

	return nil, fmt.Errorf("Target %s not found", name)
}

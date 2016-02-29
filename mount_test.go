package fuse_test

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"golang.org/x/net/context"

	"github.com/jacobsa/fuse"
	"github.com/jacobsa/fuse/fuseutil"
)

////////////////////////////////////////////////////////////////////////
// minimalFS
////////////////////////////////////////////////////////////////////////

// A minimal fuseutil.FileSystem that can successfully mount but do nothing
// else.
type minimalFS struct {
	fuseutil.NotImplementedFileSystem
}

////////////////////////////////////////////////////////////////////////
// Tests
////////////////////////////////////////////////////////////////////////

func TestSuccessfulMount(t *testing.T) {
	ctx := context.Background()

	// Set up a temporary directory.
	dir, err := ioutil.TempDir("", "mount_test")
	if err != nil {
		t.Fatal("ioutil.TempDir: %v", err)
	}

	defer os.RemoveAll(dir)

	// Mount.
	fs := &minimalFS{}
	mfs, err := fuse.Mount(
		dir,
		fuseutil.NewFileSystemServer(fs),
		&fuse.MountConfig{})

	if err != nil {
		t.Fatalf("fuse.Mount: %v", err)
	}

	defer func() {
		if err := mfs.Join(ctx); err != nil {
			t.Errorf("Joining: %v", err)
		}
	}()

	defer fuse.Unmount(mfs.Dir())
}

func TestNonEmptyMountPoint(t *testing.T) {
	// Set up a temporary directory.
	dir, err := ioutil.TempDir("", "mount_test")
	if err != nil {
		t.Fatal("ioutil.TempDir: %v", err)
	}

	defer os.RemoveAll(dir)

	// Add a file within it.
	err = ioutil.WriteFile(path.Join(dir, "foo"), []byte{}, 0600)
	if err != nil {
		t.Fatalf("ioutil.WriteFile: %v", err)
	}

	// Attempt to mount.
	fs := &minimalFS{}
	_, err = fuse.Mount(
		dir,
		fuseutil.NewFileSystemServer(fs),
		&fuse.MountConfig{})

	const want = "not empty"
	if err == nil || !strings.Contains(err.Error(), want) {
		t.Errorf("Unexpected error: %v", err)
	}
}

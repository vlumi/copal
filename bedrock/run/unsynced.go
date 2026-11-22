package run

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// unsyncedSource is a library that answers a program's requests to wait
// for the disk at once: fsync and fdatasync, and on macOS the F_FULLFSYNC
// the indexed file handler asks for on every close of a file it wrote.
const unsyncedSource = `#include <fcntl.h>
#include <stdarg.h>
#include <unistd.h>

#ifdef __APPLE__
static int unsynced_fsync(int fd) { (void)fd; return 0; }
static int unsynced_fcntl(int fd, int cmd, ...) {
    va_list ap;
    va_start(ap, cmd);
    void *arg = va_arg(ap, void *);
    va_end(ap);
    return cmd == F_FULLFSYNC ? 0 : fcntl(fd, cmd, arg);
}
#define INTERPOSE(with, of) \
    __attribute__((used)) static const struct { const void *with_, *of_; } interpose_##of \
        __attribute__((section("__DATA,__interpose"))) = { (const void *)&with, (const void *)&of };
INTERPOSE(unsynced_fsync, fsync)
INTERPOSE(unsynced_fcntl, fcntl)
#else
int fsync(int fd) { (void)fd; return 0; }
int fdatasync(int fd) { (void)fd; return 0; }
#endif
`

// Unsynced builds, in a directory, the library a Step runs a program with
// so that the program does not wait for the disk: for runs whose files are
// thrown away, thousands of them, where the waiting is most of the time
// and nothing is lost by a file a crash would have left half written. It
// is built with the C compiler GnuCOBOL compiles with, and returns the
// library's path.
func Unsynced(dir string) (string, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	source, library := filepath.Join(dir, "unsynced.c"), filepath.Join(dir, "unsynced.so")
	if err := os.WriteFile(source, []byte(unsyncedSource), 0o644); err != nil {
		return "", err
	}
	kind := "-shared"
	if runtime.GOOS == "darwin" {
		kind = "-dynamiclib"
	}
	if output, err := exec.Command("cc", kind, "-fPIC", "-o", library, source).CombinedOutput(); err != nil {
		return "", fmt.Errorf("cc unsynced.c: %w\n%s", err, output)
	}
	return library, nil
}

// preload is the environment variable the dynamic loader takes a library
// to load before a program's own from.
func preload() string {
	if runtime.GOOS == "darwin" {
		return "DYLD_INSERT_LIBRARIES"
	}
	return "LD_PRELOAD"
}

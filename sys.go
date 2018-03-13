package gosys

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"strings"
)

func TempDir() string {
	return strings.TrimRight(os.TempDir(), `/\`)
}

func TempFilename(ext string) string {
	if ext != "" {
		ext = "." + ext
	}
	return fmt.Sprintf("%s/tmp%016x%s", TempDir(), rand.Uint64(), ext)
}

func FileExists(path string) bool {
	st, err := os.Stat(path)
	return err == nil && st != nil
}

// FileExt returns file extension  without dot.
// For ex: FileExt("//path/path/file.ZIP") -> "zip"
//         FileExt("//path/path/file")     -> ""
func FileExt(filename string) string {
	return strings.ToLower(strings.TrimPrefix(path.Ext(filename), "."))
}

func FileSize(path string) int64 {
	st, err := os.Stat(path)
	for err == nil && (st.Mode()&os.ModeSymlink) != 0 { // is symlink
		if path, err = os.Readlink(path); err == nil {
			st, err = os.Stat(path)
		} else {
			return 0
		}
	}
	if err == nil && st != nil {
		return st.Size()
	}
	return 0
}

func IsDir(filename string) bool {
	f, _ := os.Stat(filename)
	return f != nil && f.IsDir()
}

func IsSymLink(filename string) bool {
	st, err := os.Stat(filename)
	return err == nil && (st.Mode()&os.ModeSymlink) != 0
}

func UserHomeDir() (dir string) {
	for _, name := range [...]string{
		"HOME",               // /Users/{username} *-nix
		"HOMEPATH",           // \Users\{username}
		"LOCALAPPDATA",       // C:\Users\{username}\AppData\Local
		"APPDATA",            // C:\Users\{username}\AppData\Roaming
		"CSIDL_APPDATA",      // C:\Users\{username}\AppData\Roaming
		"ProgramData",        // C:\ProgramData
		"CommonProgramFiles", // C:\Program Files\Common Files
		"CD",                 // The current directory
	} {
		if dir = os.Getenv(name); dir != "" && FileExists(dir) { // dir exists
			dir = strings.TrimRight(dir, "/")
			dir = strings.TrimRight(dir, "\\")
			break
		}
	}
	return
}

func DirSize(dir string) (n int64) {
	FetchDir(dir, func(info os.FileInfo) error {
		if info.IsDir() {
			n += DirSize(dir + "/" + info.Name())
		} else {
			n += info.Size()
		}
		return nil
	})
	return
}

func FetchDir(dir string, fn func(os.FileInfo) error) error {
	f, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer f.Close()

	names, err := f.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, filename := range names {
		if info, err := os.Lstat(dir + "/" + filename); os.IsNotExist(err) {
			continue
		} else if err != nil {
			return err
		} else if err = fn(info); err != nil {
			return err
		}
	}
	return nil
}

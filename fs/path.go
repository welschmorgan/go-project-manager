package fs

import (
	io_fs "io/fs"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
)

type Path string
type PathEnv map[string]func() string

var env PathEnv = make(PathEnv)

func PutPathEnv(k string, v func() string) func() string {
	previousVal, _ := env[k]
	env[k] = v
	return previousVal
}

func SetPathEnv(e PathEnv) {
	env = e
}

func GetPathEnv() PathEnv {
	return env
}

func (p Path) Raw() string {
	return string(p)
}

func (p Path) String() string {
	return p.Expand()
}

func (p Path) Join(with ...string) Path {
	parts := append([]string{}, p.Raw())
	parts = append(parts, with...)
	return Path(filepath.Join(parts...))
}

func (p Path) IsAbs() bool {
	return filepath.IsAbs(p.Expand())
}

func (p Path) Mkdir() error {
	return os.MkdirAll(p.Expand(), 0755)
}

func (p Path) Stat() (os.FileInfo, error) {
	return os.Stat(p.Expand())
}

func (p Path) IsDir() bool {
	if fi, err := p.Stat(); err != nil || os.IsNotExist(err) {
		return false
	} else {
		return fi.IsDir()
	}
}

func (p Path) Abs() (string, error) {
	return filepath.Abs(p.Expand())
}

func (p Path) Base() string {
	return filepath.Base(p.Raw())
}

func (p Path) Dir() Path {
	return Path(filepath.Dir(p.Raw()))
}

func (p Path) Exists() bool {
	_, err := os.Stat(p.Raw())
	return err == nil || os.IsExist(err)
}

func (p Path) Replace(s, repl string, n int) Path {
	return Path(strings.Replace(p.Raw(), s, repl, n))
}

func (p Path) ReplaceAll(s, repl string) Path {
	return Path(strings.ReplaceAll(p.Raw(), s, repl))
}

func (p Path) TrimSpace() Path {
	return Path(strings.TrimSpace(p.Raw()))
}

func (p Path) ReadFile() (data []byte, err error) {
	return os.ReadFile(p.Expand())
}

func (p Path) WriteFile(data []byte) error {
	return os.WriteFile(p.Expand(), data, 0755)
}

func (p Path) Chdir() error {
	if err := os.Chdir(p.Expand()); err != nil {
		return err
	}
	return nil
}
func (p Path) ReadDir() (entries []io_fs.DirEntry, err error) {
	return os.ReadDir(p.Expand())
}

func (p Path) Expand() string {
	s := string(p)
	re := regexp.MustCompile(`\$\{([\w]+)\}`)
	matches := re.FindAllStringSubmatch(s, -1)
	for _, match := range matches {
		getter := env[match[1]]
		if getter == nil {
			getter = func() string { return "" }
		}
		s = strings.ReplaceAll(s, match[0], getter())
	}
	return s
}

func ExpandPath(p string) string {
	return Path(p).Expand()
}

func init() {
	env["tmp_dir"] = func() string { return os.TempDir() }
	env["tmp"] = env["tmp_dir"]
	if usr, err := user.Current(); err != nil {
		panic(err)
	} else {
		env["home"] = func() string { return usr.HomeDir }
	}
}

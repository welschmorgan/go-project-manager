package fs

import (
	"os"
	"os/user"
	"strings"
)

func SanitizePath(p string) (string, error) {
	home := ""
	temp := os.TempDir()
	if usr, err := user.Current(); err != nil {
		return "", err
	} else {
		home = usr.HomeDir
	}
	p = strings.Replace(p, "$HOME", home, -1)
	p = strings.Replace(p, "$TMP_DIR", temp, -1)
	p = strings.Replace(p, "$TMP", temp, -1)
	return p, nil
}

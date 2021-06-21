package fs

import (
	"os"
	"os/user"
	"strings"
)

func SanitizePath(p string, repl map[string]string) string {
	home := ""
	temp := os.TempDir()
	if usr, err := user.Current(); err != nil {
		panic(err)
	} else {
		home = usr.HomeDir
	}
	p = strings.ReplaceAll(p, "$HOME", home)
	p = strings.ReplaceAll(p, "$TMP_DIR", temp)
	p = strings.ReplaceAll(p, "$TMP", temp)
	for k, v := range repl {
		k = "$" + strings.TrimPrefix(k, "$")
		p = strings.ReplaceAll(p, k, v)
	}
	return p
}

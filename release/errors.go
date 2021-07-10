package release

import "errors"

var ErrUserAbort error = errors.New("user-aborted release")

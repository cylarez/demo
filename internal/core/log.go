package core

import (
    "log"
    "os"
)

var Log = log.New(os.Stdout, "MY APP >> ", log.LstdFlags|log.Lshortfile)

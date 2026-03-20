package logger

import (
    "log"
    "os"
)

// Made by Claude

var Log *log.Logger

func Init() error {
    home, err := os.UserHomeDir()
    if err != nil {
        return err
    }
    file, err := os.OpenFile(home + "/.config/clit/clit.log", os.O_CREATE | os.O_WRONLY | os.O_APPEND, 0644)
    if err != nil {
        return err
    }
    Log = log.New(file, "", log.Ltime|log.Ldate|log.Lshortfile)
    return nil
}

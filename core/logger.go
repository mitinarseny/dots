package core

import (
    "bufio"
    "io"
    "log"
    "os"

    "github.com/mitinarseny/gol"
)

var (
    logger = gol.New(log.New(os.Stdout, "", 0))
)

func loggerWriter() io.Writer {
    pipeReader, pipeWriter := io.Pipe()

    scanner := bufio.NewScanner(pipeReader)
    scanner.Split(bufio.ScanLines)

    go func() {
        for scanner.Scan() {
            logger.Println(string(scanner.Bytes()))
        }
    }()
    return pipeWriter
}

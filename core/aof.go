package core

import (
	"fmt"
	"os"
	"strings"

	"github.com/gvp-alekhya/VelociStore/config"
)

func GenerateDumpAOF() {
	// Create a file
	fp, err := os.OpenFile(config.AOF_FILE_NAME, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModeAppend)

	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer fp.Close() // Close the file when done
	for key, value := range RedisStore {
		dumpKey(fp, key, value)
	}

	fmt.Println("Dumped AOF file")
}

func dumpKey(fp *os.File, key string, obj *Obj) {
	command := fmt.Sprintf("SET %s %s", key, obj.Value)
	tokens := strings.Split(command, " ")
	fp.Write(Encode(tokens, false))
}

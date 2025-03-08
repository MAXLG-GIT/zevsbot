package messages

import (
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"os"
	"sync"
)

//type Once struct {
//	// contains filtered or unexported fields
//	messages map[string]string
//}

var (
	once     sync.Once
	messages map[string]string
)

func Init(filePath string) {
	once.Do(func() {
		file, err := os.ReadFile(filePath)
		if err != nil {
			log.Error().Msg(err.Error())
			return
		}

		err = yaml.Unmarshal(file, &messages)
		if err != nil {
			return
		}
	})
}

func GetMessage(msgName string) string {

	return messages[msgName]
}

package utilities

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
)

func DownloadFile(url string) (*os.File, error) {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Error().Msg(err.Error())
		}
	}(resp.Body)

	// Create the file
	out, err := os.CreateTemp(os.Getenv("ZEVS_TMP_DIR"), "*"+path.Base(url))
	if err != nil {
		return nil, err
	}
	//defer out.Close()
	// Write the body to file

	_, err = io.Copy(out, resp.Body)
	_, err = out.Seek(0, 0)

	return out, err
}

func GetSubstring(str string, offset, length int) (string, error) {

	if len(str) <= offset || offset+length > len(str) {
		return "", fmt.Errorf("invalid offset/length")
	}
	return string([]rune(str)[offset : offset+length]), nil
}

type FlexibleString string

func (fs *FlexibleString) UnmarshalJSON(data []byte) error {

	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*fs = FlexibleString(s)
		return nil
	}

	var i int
	if err := json.Unmarshal(data, &i); err == nil {
		*fs = FlexibleString(strconv.Itoa(i))
		return nil
	}

	return fmt.Errorf("failed to parse %s", string(data))
}

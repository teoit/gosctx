package core

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/teoit/gosctx/common"

	"github.com/teoit/gosctx/configs"
)

func CreateIdDate() int64 {
	mu := sync.Mutex{}
	mu.Lock()

	currentTime := time.Now()
	location, err := time.LoadLocation(configs.Location)
	if err != nil {
		return currentTime.UnixNano()
	}

	timeLoc := currentTime.In(location)
	time := timeLoc.Format(common.FORMAT_ID_DATE)
	time = strings.ReplaceAll(time, ".", "")

	id, err := strconv.ParseInt(time, 10, 64)

	if err != nil {
		return timeLoc.UnixNano()
	}

	defer mu.Unlock()
	return id
}

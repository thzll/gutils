package rbacModel

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

var oldTime time.Time
var instanceIndex int
var newInstanceIdLock sync.Mutex

func GetNewOrderId() int64 {
	newInstanceIdLock.Lock()
	defer newInstanceIdLock.Unlock()
	nowt := time.Now()
	if oldTime.Unix() == nowt.Unix() {
		instanceIndex++
	} else {
		oldTime = nowt
		instanceIndex = 0
	}
	instanceName := fmt.Sprintf("%02d%02d%02d%02d%02d%02d%03d",
		nowt.Year(), nowt.Month(), nowt.Day(), nowt.Hour(), nowt.Minute(), nowt.Second(), instanceIndex)
	id, _ := strconv.ParseInt(instanceName, 10, 64)
	return id
}

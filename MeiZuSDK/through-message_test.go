package meizupush

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestJson(t *testing.T) {
	pushTimeInfo := T_PushTimeInfo{
		OffLine:   1,
		ValidTime: 3,
	}

	advanceInfo := T_AdvanceInfo{
		FixSpeed:     0,
		FixSpeedRate: 0,
	}

	message := ThroughMessage{
		Title:        "title",
		Content:      "content",
		PushTimeInfo: pushTimeInfo,
		AdvanceInfo:  advanceInfo,
	}
	j, _ := json.MarshalIndent(message, "", "")

	fmt.Println(string(j))

}

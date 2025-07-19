package log

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestJsonBytes_MarshalLogObject(t *testing.T) {
	ctx := BackgroundCtxWithRandomId()
	tmp, _ := json.Marshal(map[string]interface{}{
		"a": []map[string]interface{}{{"12": []int{1, 2, 3, 4}}},
		"b": nil,
		"c": "test",
	})
	cases := [][]byte{
		nil,
		[]byte(""),
		[]byte("test"),
		[]byte("go tool test2json -t /private/var/folders"),
		[]byte(`{"1": 12}`),
		tmp,
	}

	for i, str := range cases {
		DebugEx(ctx, fmt.Sprintf("test %d", i), Json("data", str))
	}

	DebugEx(ctx, "id", "data", RetrieveSessionId(ctx))
}

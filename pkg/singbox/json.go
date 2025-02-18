package singbox

import (
	"bytes"
	
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common/json"
)

func MarshalJSON(o *option.Options) ([]byte, error) {
	b := bytes.NewBuffer(nil)

	e := json.NewEncoder(b)
	e.SetIndent("", "  ")

	if err := e.Encode(o); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

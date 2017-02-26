package lua

import (
	"testing"
	"time"
)

func TestEvalConfig(t *testing.T) {
	config, err := EvalConfig(`
        user = setmetatable({}, {
            __tostring = "user"
        })
        config = {
            proxy = "caddy",
            addr = ":8000",
            env = user,
        }
    `, 10*time.Second)
	if err != nil {
		t.Fatalf("error while parsing config: %v", err)
	}
	if !(config["proxy"].(string) == "caddy") {
		t.Fail()
	}
	if !(config["addr"].(string) == ":8000") {
		t.Fail()
	}
	t.Logf("%+v", config["env"])
	if !(config["env"].(string) == "user") {
		t.Fail()
	}
	t.Logf("config: %+v", config)
}

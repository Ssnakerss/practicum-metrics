package app

import "testing"

func TestJsonCfg_Load(t *testing.T) {
	sCfg := serverJsonConfig{}
	err := loadJsonConfig("./server_config.json", &sCfg)
	if err != nil {
		t.Error(err)
	}
	t.Log(sCfg)

	aCfg := agentJsonConfig{}
	err = loadJsonConfig("./agent_config.json", &aCfg)
	if err != nil {
		t.Error(err)
	}
	t.Log(aCfg)
}

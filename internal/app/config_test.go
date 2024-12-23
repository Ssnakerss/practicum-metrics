package app

import "testing"

func TestJsonCfg_Load(t *testing.T) {
	sCfg := serverJSONConfig{}
	err := loadJSONConfig("./server_config.json", &sCfg)
	if err != nil {
		t.Error(err)
	}
	t.Log(sCfg)

	aCfg := agentJSONConfig{}
	err = loadJSONConfig("./agent_config.json", &aCfg)
	if err != nil {
		t.Error(err)
	}
	t.Log(aCfg)
}

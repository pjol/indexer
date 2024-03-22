package indexer

type Listener struct {
	Owner    string `json:"listener_owner"`
	Contract string `json:"contract"`
	Address  string `json:"address"`
	Service  string `json:"service"`
	Secret   string `json:"secret"`
	Value    int    `json:"value"`
}

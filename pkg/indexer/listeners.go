package indexer

type Listener struct {
	Owner   string `json:"listener_owner"`
	Address string `json:"address"`
	Service string `json:"service"`
	Secret  string `json:"endpoint"`
	Value   int    `json:"value"`
}

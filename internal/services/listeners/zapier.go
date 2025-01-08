package listeners

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/citizenwallet/indexer/pkg/indexer"
)

func SendZapierTransaction(ls *indexer.Listener, tx *indexer.Transfer) error {
	url := ls.LocationId
	amount := float64(int(tx.Value.Int64())*ls.Value) / 1000000

	body := []byte(fmt.Sprintf(`{
    "amount": %.2f,
    "txhash": "%s",
    "from": "%s"
  }`, amount, tx.TxHash, tx.From))

	r, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error creating http request: %s", err)
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", ls.Secret))

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err = io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	fmt.Printf("sending zap for address %s\n", tx.To)

	if res.StatusCode != http.StatusOK {
		fmt.Printf("error: http request rejected with code %d\n", res.StatusCode)
		fmt.Println(string(body))
		return fmt.Errorf("error: http request rejected with code %d", res.StatusCode)
	}

	return nil
}

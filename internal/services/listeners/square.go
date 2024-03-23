package listeners

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/citizenwallet/indexer/pkg/indexer"
	"github.com/google/uuid"
)

type Post struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	UserId int    `json:"userId"`
}

func SendSquareTransaction(ls *indexer.Listener, tx *indexer.Transfer) error {

	url := "https://connect.squareup.com/v2/payments"
	note := fmt.Sprintf("Payment in %s from %s", ls.TokenName, tx.From)
	amount := (int(tx.Value.Int64()) * ls.Value) / 1000000
	location := ls.LocationId

	body := []byte(fmt.Sprintf(`{
    "amount_money": {
      "amount": %d,
      "currency": "USD"
    },
    "idempotency_key": "%s",
    "source_id": "EXTERNAL",
    "location_id": "%s",
    "external_details": {
      "source": "%s",
      "type": "CRYPTO",
      "source_id": "%s"
    },
    "note": "%s",
    "accept_partial_authorization": false,
    "customer_id": ""
  }`, amount, uuid.New(), location, ls.TokenName, tx.From, note))

	r, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	r.Header.Add("Square-Version", "2024-03-20")
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", ls.Secret))
	r.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return err
	}

	body, err = io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	// bodyString := string(body)
	// fmt.Println(bodyString)

	post := &Post{}

	err = json.NewDecoder(res.Body).Decode(post)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusAccepted {
		fmt.Println(post)
	}
	return nil
}

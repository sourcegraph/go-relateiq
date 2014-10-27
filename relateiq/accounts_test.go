package relateiq

import (
	"bytes"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestAccount_json(t *testing.T) {
	jsonStr := []byte(`{"id":"abcd1234","name":"ABC Corp","modifiedDate":1414001420560}`)
	obj := Account{ID: "abcd1234", Name: "ABC Corp", ModifiedDate: Time{time.Unix(1414001420, 560000000).In(time.UTC)}}

	var gotObj Account
	if err := json.Unmarshal(jsonStr, &gotObj); err != nil {
		t.Error(err)
	}
	if gotObj != obj {
		t.Errorf("got obj %+v, want %+v", gotObj, obj)
	}

	gotJSONStr, err := json.Marshal(obj)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(gotJSONStr, jsonStr) {
		t.Errorf("got JSON %q, want %q", gotJSONStr, jsonStr)
	}
}

func TestAccountsService_List(t *testing.T) {
	setup()
	defer teardown()

	want := []*Account{
		{ID: "id", Name: "name", ModifiedDate: Time{time.Unix(123, 0)}},
	}

	var called bool
	mux.HandleFunc("/accounts", func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"_ids":   "a,b",
			"_start": "1",
			"_limit": "1",
		})
		writeJSON(w, struct {
			Objects interface{} `json:"objects"`
		}{want})
	})

	opt := AccountsListOptions{
		IDs:         []string{"a", "b"},
		ListOptions: ListOptions{Start: 1, Limit: 1},
	}
	accounts, _, err := client.Accounts.List(opt)
	if err != nil {
		t.Fatal(err)
	}

	if !called {
		t.Fatal("!called")
	}

	for _, a := range want {
		normalizeTime(&a.ModifiedDate)
	}

	if !reflect.DeepEqual(accounts, want) {
		t.Errorf("Accounts.List returned %+v, want %+v", accounts, want)
	}
}

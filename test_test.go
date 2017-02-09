package main

import (
	"testing"
	"net/http"

	"net/http/httptest"

	"os"
	"net/url"
	"bytes"

	"io/ioutil"
)


func TestHorse(t *testing.T) {
	resp, err := http.Get("http://localhost:8000/horse/e5")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Ожидается код %d. Получен код %d", http.StatusOK, resp.StatusCode)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if len(body)/3 != 8 {
		t.Errorf("Ожидается %d ходов. Получено %d ходов - %s",8, len(body)/3, body)
	}
}
 func TestTask1(t *testing.T) {

	 buffer := new(bytes.Buffer)
	 params := url.Values{}
	 params.Set("id", "1")
	 params.Set("text", "test")
	 buffer.WriteString(params.Encode())
	req, err := http.NewRequest("POST","http://localhost:8000/md5",buffer)
	if err != nil {
		t.Fatal(err)
	}
	 req.Header.Set("content-type", "application/x-www-form-urlencoded")

	 client := &http.Client{}
	 resp, err := client.Do(req)
	 defer resp.Body.Close()
	 body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Ожидается код %d. Получен код %d", http.StatusOK, resp.StatusCode)
	}



	 if string(body) != "974413f148af9aed469d5911aea2d02d" {
		 t.Errorf("Неверный результат! Ожидается %s Получен %s", "974413f148af9aed469d5911aea2d02d",body)
	 }

}

func TestMain(m *testing.M) {
	server := httptest.NewServer(getEngine())
	defer server.Close()

	os.Exit(m.Run())
}
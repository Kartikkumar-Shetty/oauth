package main

import (
	"bytes"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func main() {
	http.HandleFunc("/tenants", func(w http.ResponseWriter, r *http.Request) {
		req, err := http.NewRequest(http.MethodPost, "https://management.azure.com/tenants?api-version=2016-06-01", bytes.NewReader([]byte("")))
		if err != nil {
			fmt.Println(err)
			return
		}

		token := ``

		req.Header.Add("Authorization", "Bearer"+" "+token)

		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(body))
		return

	})
	http.HandleFunc("/path", func(w http.ResponseWriter, r *http.Request) {
		code, _ := fmt.Fprintf(w, "%s", r.URL.Query().Get("code"))
		fmt.Println(code)
		return

	})
	http.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://login.microsoftonline.com/common/oauth2/authorize?client_id=xxxxxxxxxx-xxxxxxx-x-xxx&response_type=code&redirect_uri=http%3A%2F%2Flocalhost%3A8081%2Fpath&response_mode=query&resource=https%3A%2F%2Fmanagement.azure.com%2F&state=12345", http.StatusSeeOther)
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		data := url.Values{
			"client_id":     {""},
			"resource":      {"https://management.azure.com/"},
			"code":          {""},
			"redirect_uri":  {"http://localhost:8081/path"},
			"grant_type":    {"authorization_code"},
			"client_secret": {""},
			"prompt":        {"none"},
		}

		client := http.Client{}
		resp, err := client.PostForm("https://login.microsoftonline.com/common/oauth2/token", data)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		w.Write(body)
		fmt.Println(string(body))
		return
	})
	http.HandleFunc("/starttenant", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://login.microsoftonline.com/<tenantid>/oauth2/authorize?client_id=xxxxxxxxxxxxx-xxx-x-x-x&response_type=code&redirect_uri=http%3A%2F%2Flocalhost%3A8081%2Fpath&response_mode=query&state=234234", http.StatusSeeOther)
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	http.HandleFunc("/authtenant", func(w http.ResponseWriter, r *http.Request) {
		data := url.Values{
			"client_id":     {""},
			"resource":      {"https://management.azure.com/"},
			"code":          {""},
			"redirect_uri":  {"http://localhost:8081/path"},
			"grant_type":    {"authorization_code"},
			"client_secret": {""},
		}

		client := http.Client{}
		resp, err := client.PostForm("https://login.microsoftonline.com/<tenantid>/oauth2/token", data)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		w.Write(body)
		fmt.Println(string(body))
		return

	})

	log.Fatal(http.ListenAndServe(":8081", nil))

}

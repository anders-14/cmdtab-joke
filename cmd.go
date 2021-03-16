package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/rwxrob/cmdtab"
	"github.com/rwxrob/conf-go"
)

const (
	url string = "https://icanhazdadjoke.com/"
)

func init() {
	x := cmdtab.New("joke")
	x.Summary = "Fetches you a joke"
	x.Usage = "[]"

	x.Method = func(args []string) error {
		config, err := conf.New()
		if err != nil {
			return err
		}
		err = config.Load()
		if err != nil {
			return err
		}

		if len(args) > 0 {

		} else {
			joke, err := fetchJoke()
			if err != nil {
				return err
			}
			fmt.Println(joke.Joke)
		}
		return nil
	}
}

type Response struct {
	ID   string `json:"id"`
	Joke string `json:"joke"`
}

func fetchJoke() (Response, error) {
	var joke Response

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Response{}, err
	}
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: time.Second * 5}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(body, &joke)
	if err != nil {
		log.Fatal(err)
	}

	return joke, nil
}

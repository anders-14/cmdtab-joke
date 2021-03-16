package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/rwxrob/cmdtab"
	"github.com/rwxrob/conf-go"
)

const (
	url string = "https://icanhazdadjoke.com/"
)

func init() {
	x := cmdtab.New("joke", "save", "list")
	x.Summary = `Fetches you a joke`
	x.Usage = `[save|list]`
	x.Description = ``

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
			switch args[0] {
			case "save":
				saveJoke(config.Get("joke.last"), config)
			case "list":
				jokes := getJokes(config)
				for k, v := range jokes {
					fmt.Printf("%v) %v\n", k+1, v)
				}
			case "delete":
        i, err := strconv.ParseInt(args[1], 10, 0)
        if err != nil {
          return err
        }
				deleteJoke(i, config)
			default:
				return x.UsageError()
			}
			return nil
		}
		res, err := fetchJoke()
		if err != nil {
			return err
		}
		fmt.Println(res.Joke)
		config.SetSave("joke.last", res.Joke)
		return nil
	}
}

func getJokes(config *conf.Config) []string {
	var jokes = []string{}
	json.Unmarshal([]byte(config.Get("joke.saved")), &jokes)
	return jokes
}

func saveJoke(joke string, config *conf.Config) {
	jokes := getJokes(config)
	jokes = append(jokes, joke)
	encodedJokes, _ := json.Marshal(jokes)
	config.SetSave("joke.saved", string(encodedJokes))
}

func deleteJoke(i int64, config *conf.Config) {
  i = i-1
	jokes := getJokes(config)
  if int(i) > len(jokes)-1 || i < 0 {
    return
  }
	copy(jokes[i:], jokes[i+1:])
	jokes[len(jokes)-1] = ""
	jokes = jokes[:len(jokes)-1]
	encodedJokes, _ := json.Marshal(jokes)
	config.SetSave("joke.saved", string(encodedJokes))
}

type response struct {
	Joke string `json:"joke"`
}

func fetchJoke() (response, error) {
	var joke response

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return response{}, err
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

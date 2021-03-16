package cmd

import (
	"encoding/json"
	"fmt"
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
	x := cmdtab.New("joke", "save", "list", "delete")
	x.Summary = `Fetches you a joke`
	x.Usage = `[save|list|delete]`
	x.Description = `
    Get a funny (or not) joke by calling the command *joke*.
    The joke is stored in *joke.last*.

    When *save* is passed the joke stored in *joke.last* is saved
    to *joke.saved*.

    When *list* is passed all the jokes stored in *jokes.saved* is
    shown with a numbered index.

    When *delete is passed and an argument corresponding to the
    index of the joke you want to delete, shown using *list*, the
    joke in question is removed from *joke.saved*.`

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
				err = deleteJoke(i, config)
				if err != nil {
					return err
				}
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

func deleteJoke(i int64, config *conf.Config) error {
	i = i - 1
	jokes := getJokes(config)
	if int(i) > len(jokes)-1 || i < 0 {
		return fmt.Errorf("Invalid argument")
	}
	copy(jokes[i:], jokes[i+1:])
	jokes[len(jokes)-1] = ""
	jokes = jokes[:len(jokes)-1]
	encodedJokes, _ := json.Marshal(jokes)
	config.SetSave("joke.saved", string(encodedJokes))
	return nil
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

	err = json.NewDecoder(res.Body).Decode(&joke)
	if err != nil {
		log.Fatal(err)
	}

	return joke, nil
}

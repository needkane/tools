package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/urfave/cli"
)

type QueryJson struct {
	Query string `json:"query"`
}
type StartForkResult struct {
	Data Data `json:"data"`
}
type Data struct {
	Repository Repository `json:"repository"`
}
type Repository struct {
	Stargazers Stargazers `json:"stargazers"`
	Forks      Forks      `json"forks"`
}
type Stargazers struct {
	TotalCount int    `json:"totalCount"`
	Edges      []Edge `json"edges"`
}
type Forks struct {
	TotalCount int    `json:"totalCount"`
	Edges      []Edge `json"edges"`
}
type Edge struct {
	Node Node `json"node"`
}
type Node struct {
	Url string `json:"url"`
}

const queryContent = `query { repository(owner:"dappledger", name:"AnnChain") { stargazers (last : 100){ totalCount edges {  node { url } } }  forks(last : 100) { totalCount edges { node { url  } } } } }`

func main() {
	app := cli.NewApp()
	app.Name = "StarFork"
	app.Version = "0.1"
	app.UsageText = "./StarFork -u [Your github username]"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "username, u",
			Usage: "Your github username",
		},
	}
	app.Action = queryByGithubAPI
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
func queryByGithubAPI(ctx *cli.Context) error {
	if ctx.String("username") == "" {
		return errors.New("Please specify a username by '-u [Your github username]'")
	}
	qj := QueryJson{queryContent}
	bytez, err := json.Marshal(qj)
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewBuffer(bytez))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "bearer [Your github account token]")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bytez, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	sfr := StartForkResult{}
	err = json.Unmarshal(bytez, &sfr)
	if err != nil {
		log.Fatal(err)
	}
	username := ctx.String("username")
	err = fmt.Errorf("%s, please star github.com/dappledger/AnnChain", username)
	for _, v := range sfr.Data.Repository.Stargazers.Edges {
		if v.Node.Url[19:] == username {
			err = nil
			break
		}
	}
	if err != nil {
		return err
	}
	err = fmt.Errorf("%s, please fork github.com/dappledger/AnnChain", username)
	for _, v := range sfr.Data.Repository.Forks.Edges {
		l := len(v.Node.Url)
		if v.Node.Url[19:l-9] == username {
			err = nil
			break
		}
	}
	return err
}

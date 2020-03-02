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
	"strings"
	"sync"

	"github.com/dappledger/AnnChain/eth/crypto"
	"github.com/urfave/cli"
)

type QueryJsonV3 struct {
}
type StarInfoV3 struct {
	StarredAt string `json:"starred_at"`
	User      User   `json:"user"`
}
type User struct {
	Login string `json:"login"`
}
type QueryJson struct {
	Query string `json:"query"`
}
type StarForkResult struct {
	Data Data `json:"data"`
	//failed response
	Message string `json:"message"`
}
type FailedResponse struct {
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

var (
	bearToken    = "Your github account token"
	pageCount    = 30
	maxLastCount = 100
	queryContent = fmt.Sprintf(`query { repository(owner:"dappledger", name:"AnnChain") { stargazers (last : %d){ totalCount edges {  node { url } } }  forks(last : %d) { totalCount edges { node { url  } } } } }`, maxLastCount, maxLastCount)
	apiV3StarUrl = "https://api.github.com/repos/dappledger/annchain/stargazers"
	apiV3ForkUrl = "https://api.github.com/repos/dappledger/annchain/fork"
)

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

type Result struct {
	Contract string `json"contract"`
	Privkey  string `json:"privkey"`
	Address  string `json:"address"`
}

func queryByAPIV3(url, username string, surplus int) error {
	pages := surplus / pageCount
	if surplus%pageCount != 0 {
		pages += 1
	}
	var wg sync.WaitGroup
	wg.Add(pages)
	getInfo := make(chan bool, 1)
	getInfo <- false
	fGetStarInfo := func(pageNo int) {
		req, err := http.NewRequest("GET", fmt.Sprintf("%s?page=%d&access_token=%s", url, pageNo, bearToken), nil)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Accept", "application/vnd.github.v3.star+json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("error is (%v),please check your network to make sure you can connect to github api", err)
		}
		defer resp.Body.Close()
		bytez, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		var siV3s = []StarInfoV3{}
		err = json.Unmarshal(bytez, &siV3s)
		if err != nil {
			log.Fatal(err)
		}
		for _, v := range siV3s {
			if v.User.Login == username {
				<-getInfo
				getInfo <- true
				break
			}
		}
		wg.Done()
	}
	for i := 1; i <= pages; i++ {
		j := i
		go func() {
			fGetStarInfo(j)
		}()
	}
	wg.Wait()
	if <-getInfo {
		return nil
	}
	return fmt.Errorf("%s, please star and fork github.com/dappledger/AnnChain", username)
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
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", bearToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("error is (%v),please check your network to make sure you can connect to github api", err)
	}
	defer resp.Body.Close()
	bytez, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var sfr = new(StarForkResult)
	err = json.Unmarshal(bytez, &sfr)
	if err != nil {
		log.Fatal(err)
	}
	if sfr.Message != "" {
		return fmt.Errorf("query failed: %s", sfr.Message)
	}

	username := ctx.String("username")
	err = fmt.Errorf("%s, please star and fork github.com/dappledger/AnnChain", username)
	for _, v := range sfr.Data.Repository.Stargazers.Edges {
		if v.Node.Url[19:] == username {
			err = nil
			break
		}
	}
	// query history forks too slow
	ignoreQueryForks := false
	if err != nil {
		if sfr.Data.Repository.Stargazers.TotalCount > maxLastCount {
			err = queryByAPIV3(apiV3StarUrl, username, sfr.Data.Repository.Stargazers.TotalCount-maxLastCount)
			if err != nil {
				return err
			} else {
				ignoreQueryForks = true
			}
		} else {
			return err
		}
	}
	if !ignoreQueryForks {
		err = fmt.Errorf("%s, please fork github.com/dappledger/AnnChain", username)
		for _, v := range sfr.Data.Repository.Forks.Edges {
			l := len(v.Node.Url)
			if v.Node.Url[19:l-9] == username {
				err = nil
				break
			}
		}
		if err != nil {
			return err
		}
	}
	result, err := createAccount()
	if err != nil {
		return err
	}
	result.Contract = strings.ToLower("0x04D7A824b3301e67Ef34024E9dc79445E54D5aF7")
	bytez, err = json.Marshal(result)
	if err != nil {
		return err
	}
	fmt.Println(string(bytez))
	return err
}

func createAccount() (result Result, err error) {
	var (
		privkeyBytes []byte
		addrBytes    []byte
	)

	privkey, errG := crypto.GenerateKey()
	if errG != nil {
		err = errG
		return
	}

	privkeyBytes = crypto.FromECDSA(privkey)

	address := crypto.PubkeyToAddress(privkey.PublicKey)
	addrBytes = address.Bytes()

	result.Privkey = fmt.Sprintf("%x", privkeyBytes)
	result.Address = fmt.Sprintf("%x", addrBytes)
	return
}

package main

import (
	"fmt"
	"strings"

	jms "github.com/wlin-jin/jumpserver-go-sdk"
)

func main() {

	jmsURL := "http://localhost"
	//authURL := "http://localhost/api/users/v1/auth/"
	//username := "admin"
	//password := "admin"

	tp := jms.TokenAuthTransport{
		//Username:  strings.TrimSpace(username),
		//Password:  strings.TrimSpace(password),
		//AuthURL:   strings.TrimSpace(authURL),
		Token: "909f6f7bf3e147c3809d81a0e27addb3",
	}

	client, err := jms.NewClient(tp.Client(), strings.TrimSpace(jmsURL))
	if err != nil {
		fmt.Printf("\nerror: %v\n", err)
		return
	}

	// search
	matchUserList, err := client.Users.Search("weilin")
	if err != nil {
		fmt.Printf("\nerror: %v\n", err)
		return
	}
	for _, user := range matchUserList {
		fmt.Println(user)
	}
	//
	// list
	userList, _, err := client.Users.GetList()
	if err != nil {
		fmt.Printf("\nerror: %v\n", err)
		return
	}
	for _, user := range userList {
		fmt.Println(user.Email)
	}

	groups, err := client.Users.Groups("test")
	if err != nil {
		fmt.Printf("\nerror: %v\n", err)
		return
	}
	for _, group := range groups {
		fmt.Println(group)
	}

	// test for asserts
	assets, _, err := client.Assets.GetList()
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, item := range assets {
		fmt.Println(item.Protocols)
		fmt.Println(item.Platform)
		fmt.Println(item.Domain)
		fmt.Println(item.AdminUser)
		fmt.Println(item.Nodes)
		fmt.Println(item.Number)
		fmt.Println(item.Comment)
	}

	// test delete asset
	err = client.Assets.Delete("172.24.35.31")
	if err != nil {
		fmt.Println(err.Error())
	}

	// test create host
	asset, err := client.Assets.Create(&jms.AssetBody{
		AdminUser: "5a6c9561-2665-4ffe-ab5f-88c42078d85c",
		Comment:   "just for test",
		Hostname:  "testnode8",
		IP:        "10.99.99.97",
		IsActive:  true,
		Labels:    []string{},
		Nodes:     []string{"92fc7e3a-ec08-47d8-b850-0de6c3d9f8be"},
		Platform:  "Linux",
		Port:      22,
		Protocol:  "ssh",
		Protocols: []string{"ssh/22"},
	})
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(asset)
	}

	// test refresh
	resp, err := client.Perms.RefreshCache()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(resp)
	}

}

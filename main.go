package main

import (
	"os"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: qingting-dl <channel-id>")
		return
	}
	channelId := os.Args[1]
	fmt.Println("Fetching audios of channel " + channelId)

	infoUrl := GetInfoUrl(channelId, "7714479")
	response, err := http.Get(infoUrl)
	if err != nil {
		fmt.Println("Error in fetching JSON")
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))
}

func GetInfoUrl(channelId string, programId string) string {
	return "http://i.qingting.fm/wapi/channels/" + channelId + "/programs/" + programId
}

func GetDownloadUrl(filePath string) string {
	return "http://od.qingting.fm/" + filePath
}

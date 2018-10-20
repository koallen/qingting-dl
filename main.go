package main

import (
	"os"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"encoding/json"
)

type ChannelApiResponse struct {
	Data []AudioInfo
	Code int
	Total int
}

type AudioInfo struct {
	FilePath string `json:"file_path"`
	Name string
	ResId int `json:"res_id"`
	UpdateTime string `json:"update_time"`
	Duration int
	Playcount string
	Id int
	Desc string
	ChannelId string `json:"channel_id"`
	Type string
	ImgUrl string `json:"img_url"`
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: qingting-dl <channel-id>")
		return
	}
	channelId := os.Args[1]
	fmt.Println("Fetching audios of channel " + channelId)

	// request API and parse it
	infoUrl := GetChannelInfoUrl(channelId)
	response, err := http.Get(infoUrl)
	if err != nil {
		fmt.Println("Error in fetching JSON")
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	var parsedJson ChannelApiResponse
	json.Unmarshal(body, &parsedJson)

	for _, audioInfo := range parsedJson.Data {
		fmt.Println("Downloading " + audioInfo.Name)
		DownloadFile(audioInfo.Name + ".m4a", GetDownloadUrl(audioInfo.FilePath))
	}
}

func GetChannelInfoUrl(channelId string) string {
	return "http://i.qingting.fm/wapi/channels/" + channelId + "/programs/page/1/pagesize/250"
}

func GetProgramInfoUrl(channelId string, programId string) string {
	return "http://i.qingting.fm/wapi/channels/" + channelId + "/programs/" + programId
}

func GetDownloadUrl(filePath string) string {
	return "http://od.qingting.fm/" + filePath
}

func DownloadFile(filename string, url string) error {
	// Create the file
	out, err := os.Create(filename)
	if err != nil  {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil  {
		return err
	}

	return nil
}

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type ChannelInfoApi struct {
	Data ChannelInfo
	Code int
}

type ChannelInfo struct {
	ProgramCount int `json:"program_count"`
	Name         string
}

type ChannelAudioInfoApi struct {
	Data  []AudioInfo
	Code  int
	Total int
}

type AudioInfo struct {
	FilePath   string `json:"file_path"`
	Name       string
	ResId      int    `json:"res_id"`
	UpdateTime string `json:"update_time"`
	Duration   int
	Playcount  string
	Id         int
	Desc       string
	ChannelId  string `json:"channel_id"`
	Type       string
	ImgUrl     string `json:"img_url"`
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: qingting-dl <channel-id>")
		return
	}
	channelId := os.Args[1]

	// get info of the channel
	channelInfoUrl := GetChannelInfoUrl(channelId)
	channelInfoResponse, err := http.Get(channelInfoUrl)
	if err != nil {
		fmt.Println("Error in fetching JSON")
	}
	defer channelInfoResponse.Body.Close()
	channelInfoBody, err := ioutil.ReadAll(channelInfoResponse.Body)
	var parsedChannelJson ChannelInfoApi
	json.Unmarshal(channelInfoBody, &parsedChannelJson)
	fmt.Println("节目 \"" + parsedChannelJson.Data.Name + "\" 共有 " + strconv.Itoa(parsedChannelJson.Data.ProgramCount) + " 段音频")

	// request API and parse it
	audioInfoUrl := GetChannelAudioInfoUrl(channelId)
	response, err := http.Get(audioInfoUrl)
	if err != nil {
		fmt.Println("Error in fetching JSON")
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	var parsedJson ChannelAudioInfoApi
	json.Unmarshal(body, &parsedJson)

	// download all audios
	for index, audioInfo := range parsedJson.Data {
		fmt.Printf("[%2d/%2d] %s\n", index, parsedJson.Total, audioInfo.Name)
		DownloadFile(audioInfo.Name+".m4a", GetDownloadUrl(audioInfo.FilePath))
	}
}

func GetChannelInfoUrl(channelId string) string {
	return "http://i.qingting.fm/wapi/channels/" + channelId
}

func GetChannelAudioInfoUrl(channelId string) string {
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
	if err != nil {
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
	if err != nil {
		return err
	}

	return nil
}

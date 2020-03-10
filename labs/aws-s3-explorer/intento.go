package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"flag"
)

var name string

func getXML(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("Status error: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("Read body: %v", err)
	}

	return data, nil
}
type ListBucketResult struct {
	XMLName     xml.Name `xml:"ListBucketResult"`
	Text        string   `xml:",chardata"`
	Xmlns       string   `xml:"xmlns,attr"`
	Name        string   `xml:"Name"`
	Prefix      string   `xml:"Prefix"`
	Marker      string   `xml:"Marker"`
	MaxKeys     string   `xml:"MaxKeys"`
	IsTruncated string   `xml:"IsTruncated"`
	Contents    []struct {
		Text         string `xml:",chardata"`
		Key          string `xml:"Key"`
		LastModified string `xml:"LastModified"`
		ETag         string `xml:"ETag"`
		Size         string `xml:"Size"`
		StorageClass string `xml:"StorageClass"`
		Owner        struct {
			Text        string `xml:",chardata"`
			ID          string `xml:"ID"`
			DisplayName string `xml:"DisplayName"`
		} `xml:"Owner"`
	} `xml:"Contents"`
} 
func main() {
	name := flag.String("bucket", name, "bucket's name")
	flag.Parse()
	bucket := *name
	var result ListBucketResult
	if xmlBytes, err := getXML("https://"+bucket+".s3.amazonaws.com/"); err != nil {
		log.Printf("Failed to get XML: %v", err)
	} else {
		xml.Unmarshal(xmlBytes, &result)
	}

	extensions := make(map[string]int)
	objs, dirs := 0, 0
	for _, s := range result.Contents {
		key:=s.Key
		if(key[len(key)-1:]=="/"){
			dirs++
		}else{
			var ext string
			for i,c:=range key{
				if(c=='.'){
					ext=key[i+1:]
				}
			}
			extensions[ext]++
		}
		objs++
	}
	fmt.Println("AWS S3 Explorer")
	fmt.Println("Bucket Name            : ", *name)
	fmt.Println("Number of objects      : ", objs)
	fmt.Println("Number of directories  : ", dirs)
	fmt.Printf("Extensions             : ")
	for s,e := range extensions{
		fmt.Print(s,"(", e, "), ")	
	}
}

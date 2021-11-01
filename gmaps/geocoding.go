package gmaps

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/darienkentanu/suit/constants"
)

func Geocoding(address string) (lat, lng string) {
	newAddress := strings.Replace(address, " ", "+", -1)
	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?address=%s&key=%s",
		newAddress, constants.API_KEY,
	)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return "", ""
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", ""
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", ""
	}

	location := strings.Index(string(body), "location")
	lat1 := string(body[location:])
	temp := strings.SplitN(lat1[37:], ",", 2)
	lat2 := temp[0]
	lat = strings.TrimSpace(lat2)

	lng1 := temp[1]
	temp2 := strings.SplitN(lng1[24:], "},", 2)
	lng2 := temp2[0]
	lng = strings.TrimSpace(lng2)

	return lat, lng
}

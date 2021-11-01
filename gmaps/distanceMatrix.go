package gmaps

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/darienkentanu/suit/constants"
)

func Distancematrix(origin, destination string) (float64, int) {

	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/distancematrix/json?origins=%s&destinations=%s&key=%s",
		origin, destination, constants.API_KEY,
	)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return -1, -1
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return -1, -1
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return -1, -1
	}

	temp1 := strings.Index(string(body), "text")
	temp2 := string(body)[temp1:]
	temp3 := temp2[9:]
	temp4 := strings.SplitN(temp3, " ", 2)
	temp5 := temp4[0]
	temp6 := strings.TrimSpace(temp5)
	km, err := strconv.ParseFloat(temp6, 64)
	if err != nil {
		fmt.Println(err)
	}

	tmp1 := strings.Index(string(body), "duration")
	tmp2 := string(body)[tmp1:]
	tmp3 := tmp2[42:]
	tmp4 := strings.SplitN(tmp3, " ", 2)
	tmp5 := tmp4[0]
	tmp6 := strings.TrimSpace(tmp5)
	minutes, err := strconv.Atoi(tmp6)
	if err != nil {
		fmt.Println(err)
	}
	return km, minutes
}

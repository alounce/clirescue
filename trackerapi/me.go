package trackerapi

import (
	"path"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/user"

	"github.com/alounce/clirescue/cmdutil"
)

const (
	url string = "https://www.pivotaltracker.com/services/v5/me"
	settingsFileName string = ".tracker"
)

// MeResponse represents service response
type MeResponse struct {
	APIToken string `json:"api_token"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Initials string `json:"initials"`
	Timezone struct {
		Kind      string `json:"kind"`
		Offset    string `json:"offset"`
		OlsonName string `json:"olson_name"`
	} `json:"time_zone"`
}

// Me obtains information about my account
func Me() (*MeResponse, error) {
	
	response, err := makeMeRequest(url)
	if err != nil {
		return nil, err
	}
	
	user, err := parseMeResponse(response)
	if err != nil {
		return nil, err
	}

	writeCachedToken(user.APIToken)

	return user, nil
}


// Private members --------------------------------------------------------

func makeMeRequest(meURL string) (string, error) {
	
	client := &http.Client{}
	
	req, err := http.NewRequest("GET", meURL, nil)
	if err != nil {
		return "", fmt.Errorf("Could not create HTTP request: %v", err)
	}

	token := readCachedToken()
	if token == "" {
		fmt.Println("Token was not found in the cache, so requesting user name and password...")
		//user, password := retrieveUserCredentials()
		username, err := cmdutil.Ask("Username", false)
		if err != nil {
			return "", fmt.Errorf("Unable to collect user name: %v", err)
		}

		password, err := cmdutil.Ask("Password", true)
		if err != nil {
			return "", fmt.Errorf("Unable to collect password: %v", err)
		}
		req.SetBasicAuth(username, password)
	} else {
		fmt.Println("Token was found in the cache, so let's use it...")
		req.Header.Set("X-TrackerToken", token)
	}
	
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}


	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err)
		return "", err
	}
	body := string(data)
	fmt.Printf("\n****\nAPI response: \n%s\n", body)
	return body, nil
}

func parseMeResponse(body string) (*MeResponse, error) {
	response := new(MeResponse)
	err := json.Unmarshal([]byte(body), &response)
	return response, err
}

func getCachedTokenFileName() string {
	usr, err := user.Current()
	if err != nil {
		return ""
	}
	fileName := path.Join(usr.HomeDir, settingsFileName)
	return fileName
}

func readCachedToken() string {
	fileName := getCachedTokenFileName()
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil { 
		return "" 
	}
	return string(bytes)
}

func writeCachedToken(token string) {
	fileName := getCachedTokenFileName()
	ioutil.WriteFile(fileName, []byte(token), 0644)
}

package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"math"
	"net/http"
	"net/url"
	"time"
)

type ItemReqResponse struct {
	Data []struct {
		Name        string `json:"name"`
		Code        string `json:"code"`
		Level       int    `json:"level"`
		Type        string `json:"type"`
		Subtype     string `json:"subtype"`
		Description string `json:"description"`
		Effects     []struct {
			Name  string `json:"name"`
			Value int    `json:"value"`
		} `json:"effects"`
		Craft struct {
			Skill string `json:"skill"`
			Level int    `json:"level"`
			Items []struct {
				Code     string `json:"code"`
				Quantity int    `json:"quantity"`
			} `json:"items"`
			Quantity int `json:"quantity"`
		} `json:"craft"`
		Tradeable bool `json:"tradeable"`
	} `json:"data"`
	Total int `json:"total"`
	Page  int `json:"page"`
	Size  int `json:"size"`
	Pages int `json:"pages"`
}

type CharacterResponse struct {
	Data []CharacterState
}

type ActionResponse struct {
	Data struct {
		Cooldown struct {
			TotalSeconds     int       `json:"total_seconds"`
			RemainingSeconds int       `json:"remaining_seconds"`
			StartedAt        time.Time `json:"started_at"`
			Expiration       time.Time `json:"expiration"`
			Reason           string    `json:"reason"`
		} `json:"cooldown"`
		Details struct {
			Xp    int `json:"xp"`
			Items []struct {
				Code     string `json:"code"`
				Quantity int    `json:"quantity"`
			} `json:"items"`
		} `json:"details"`
		Character CharacterState `json:"character"`
	} `json:"data"`
}

type MonsterResponse struct {
	Data []struct {
		Name        string `json:"name"`
		Code        string `json:"code"`
		Level       int    `json:"level"`
		Hp          int    `json:"hp"`
		AttackFire  int    `json:"attack_fire"`
		AttackEarth int    `json:"attack_earth"`
		AttackWater int    `json:"attack_water"`
		AttackAir   int    `json:"attack_air"`
		ResFire     int    `json:"res_fire"`
		ResEarth    int    `json:"res_earth"`
		ResWater    int    `json:"res_water"`
		ResAir      int    `json:"res_air"`
		MinGold     int    `json:"min_gold"`
		MaxGold     int    `json:"max_gold"`
		Drops       []struct {
			Code        string `json:"code"`
			Rate        int    `json:"rate"`
			MinQuantity int    `json:"min_quantity"`
			MaxQuantity int    `json:"max_quantity"`
		} `json:"drops"`
	} `json:"data"`
	Total int `json:"total"`
	Page  int `json:"page"`
	Size  int `json:"size"`
	Pages int `json:"pages"`
}

type MapResponse struct {
	Data []struct {
		Name    string `json:"name"`
		Skin    string `json:"skin"`
		X       int    `json:"x"`
		Y       int    `json:"y"`
		Content struct {
			Type string `json:"type"`
			Code string `json:"code"`
		} `json:"content"`
	} `json:"data"`
	Total int `json:"total"`
	Page  int `json:"page"`
	Size  int `json:"size"`
	Pages int `json:"pages"`
}

type MoveRequest struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func doRequest(req *http.Request, response any) error {
	slog.Debug("sending request", "request", req)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read and display the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return ResponseCodeError{ArtifactsResponseCode(resp.StatusCode)}
	}

	err = json.Unmarshal(body, &response)

	if err != nil {
		slog.Error("Error parsing response: ", "error", err)
		return err
	}
	slog.Debug("received response", "response", response)
	return nil
}

func (state *CharacterState) getBankLocation() (*MoveRequest, error) {
	retval, err := state.getClosestLocation("bank", "bank")
	if err != nil {
		return nil, err
	}
	return retval, nil
}

func (state *CharacterState) getMonsterLocation(monsterName string) (*MoveRequest, error) {
	retval, err := state.getClosestLocation(monsterName, "monster")
	if err != nil {
		return nil, err
	}
	return retval, nil
}

func (state *CharacterState) getResourceLocation(resourceName string) (*MoveRequest, error) {
	retval, err := state.getClosestLocation(resourceName, "resource")
	if err != nil {
		return retval, err
	}
	return retval, nil
}

// return a MoveRequest representing the closest map location with given content
func (state *CharacterState) getClosestLocation(contentCode string, contentType string) (*MoveRequest, error) {
	myX := state.X
	myY := state.Y
	minDst := 9001.0

	retval := new(MoveRequest)

	resp, err := makeMapQuery(contentCode, contentType)
	if err != nil {
		return nil, err
	}

	for _, spot := range resp.Data {
		xDiff := math.Abs(float64(spot.X - myX))
		yDiff := math.Abs(float64(spot.Y - myY))
		dst := math.Sqrt(xDiff*xDiff + yDiff*yDiff)
		if dst < minDst {
			minDst = dst
			retval.X = spot.X
			retval.Y = spot.Y
		}
	}
	return retval, err
}

func makeMapQuery(contentCode string, contentType string) (*MapResponse, error) {
	response := new(MapResponse)

	// Define the endpoint and token
	apiURL := "https://api.artifactsmmo.com/maps"

	u, _ := url.Parse(apiURL)

	q := u.Query()
	q.Add("content_code", contentCode)
	q.Add("content_type", contentType)

	u.RawQuery = q.Encode()

	// Create the HTTP request
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+ApiToken)

	err = doRequest(req, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

type ItemQueryResponse struct {
	Data struct {
		Name        string `json:"name"`
		Code        string `json:"code"`
		Level       int    `json:"level"`
		Type        string `json:"type"`
		Subtype     string `json:"subtype"`
		Description string `json:"description"`
		Effects     []struct {
			Name  string `json:"name"`
			Value int    `json:"value"`
		} `json:"effects"`
		Craft struct {
			Skill string `json:"skill"`
			Level int    `json:"level"`
			Items []struct {
				Code     string `json:"code"`
				Quantity int    `json:"quantity"`
			} `json:"items"`
			Quantity int `json:"quantity"`
		} `json:"craft"`
		Tradeable bool `json:"tradeable"`
	} `json:"data"`
}

func (state *CharacterState) getItemInfo(itemCode string) (*ItemQueryResponse, error) {
	response := new(ItemQueryResponse)

	// Define the endpoint and token
	apiURL := "https://api.artifactsmmo.com/items/" + itemCode

	// Create the HTTP request
	req, err := http.NewRequest("GET", apiURL, bytes.NewBuffer([]byte{}))
	if err != nil {
		state.Logger.Error("Error creating request", "error", err)
		return nil, err
	}

	// Set headers
	req.Header.Set("Accept", "application/json")

	err = doRequest(req, response)

	if err != nil {
		return nil, err
	}

	return response, nil
}

// perform given action and block until cooldown is up
func (state *CharacterState) performActionAndWait(actionName string, actionData []byte) (*ActionResponse, error) {
	response := new(ActionResponse)

	// Define the endpoint and token
	apiURL := "https://api.artifactsmmo.com/my/" + state.Name + "/action/" + actionName

	// Create the HTTP request
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(actionData))
	if err != nil {
		state.Logger.Error("Error creating request", "error", err)
		return nil, err
	}

	// Set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+ApiToken)

	err = doRequest(req, response)
	if err != nil {
		return nil, err
	}

	state.updateState(response)

	cooldown := response.Data.Cooldown.RemainingSeconds
	state.Logger.Debug("Waiting finish action", "cooldown", cooldown, "action", actionName)
	time.Sleep(time.Duration(cooldown) * time.Second)

	return response, err
}

// query game for initial status of all characters
func getGameStatus() ([]CharacterState, error) {
	response := new(CharacterResponse)

	// Define the endpoint and token
	apiURL := "https://api.artifactsmmo.com/my/characters"

	// Create the HTTP request
	req, err := http.NewRequest("GET", apiURL, bytes.NewBuffer([]byte{}))
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+ApiToken)

	err = doRequest(req, response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// query game for initial status of all characters
func getMonsterDB() (*MonsterResponse, error) {
	response := new(MonsterResponse)

	// Define the endpoint and token
	apiURL := "https://api.artifactsmmo.com/monsters"

	// Create the HTTP request
	req, err := http.NewRequest("GET", apiURL, bytes.NewBuffer([]byte{}))
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+ApiToken)

	err = doRequest(req, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

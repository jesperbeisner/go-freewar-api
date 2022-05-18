package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Player struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Xp     int    `json:"xp"`
	Race   string `json:"race"`
	ClanId int    `json:"clan_id"`
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/worlds/{id}", WorldsHandler)

	log.Println("Listening on http://localhost:8050")
	log.Fatalln(http.ListenAndServe(":8050", r))
}

func WorldsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	vars := mux.Vars(r)
	worldId := vars["id"]

	worlds := map[string]string{
		"1":     "welt1",
		"2":     "welt2",
		"3":     "welt3",
		"4":     "welt4",
		"5":     "welt5",
		"6":     "welt6",
		"7":     "welt7",
		"8":     "welt8",
		"9":     "welt9",
		"10":    "welt10",
		"11":    "welt11",
		"12":    "welt12",
		"13":    "welt13",
		"14":    "welt14",
		"rp":    "rpsrv",
		"af":    "afsrv",
		"chaos": "chaos",
	}

	if _, err := worlds[worldId]; err == false {
		NotFoundResponse(w, "World not found")
		return
	}

	var players []Player

	resp, _ := http.Get("https://" + worlds[worldId] + ".freewar.de/freewar/dump_players.php")
	body, _ := ioutil.ReadAll(resp.Body)

	stringArray := strings.Split(strings.TrimSuffix(string(body), "\n"), "\n")

	for _, value := range stringArray {
		playerStats := strings.Split(value, "\t")

		id, err := strconv.Atoi(strings.TrimSpace(playerStats[0]))
		if err != nil {
			log.Fatalln(err)
		}

		xp, err := strconv.Atoi(strings.TrimSpace(playerStats[2]))
		if err != nil {
			log.Fatalln(err)
		}

		clanId, err := strconv.Atoi(strings.TrimSpace(playerStats[4]))
		if err != nil {
			log.Fatalln(err)
		}

		player := Player{
			Id:     id,
			Name:   strings.TrimSpace(playerStats[1]),
			Xp:     xp,
			Race:   strings.TrimSpace(playerStats[3]),
			ClanId: clanId,
		}

		players = append(players, player)
	}

	playersJson, err := json.Marshal(players)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}

	w.Write(playersJson)
}

func NotFoundResponse(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusNotFound)

	resp := make(map[string]string)
	resp["message"] = message

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}

	w.Write(jsonResp)
}

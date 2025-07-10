package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type PokemonResponse struct {
	Name   string   `json:"name"`
	ID     int      `json:"id"`
	Height int      `json:"height"`
	Weight int      `json:"weight"`
	Types  []string `json:"types"`
}

type PokemonAPIResponse struct {
	Name   string        `json:"name"`
	ID     int           `json:"id"`
	Height int           `json:"height"`
	Weight int           `json:"weight"`
	Types  []PokemonType `json:"types"`
}

type PokemonType struct {
	Type struct {
		Name string `json:"name"`
	} `json:"type"`
}

func getPokemonHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	resp, err := http.Get(fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", name))
	if err != nil {
		http.Error(w, "Erro ao buscar o Pokémon", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Pokémon não encontrado", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Erro ao ler resposta da PokeAPI", http.StatusInternalServerError)
		return
	}

	var pokeData PokemonAPIResponse
	if err := json.Unmarshal(body, &pokeData); err != nil {
		http.Error(w, "Erro ao decodificar a resposta: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var types []string
	for _, t := range pokeData.Types {
		types = append(types, t.Type.Name)
	}

	result := PokemonResponse{
		Name:   pokeData.Name,
		ID:     pokeData.ID,
		Height: pokeData.Height,
		Weight: pokeData.Weight,
		Types:  types,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/pokemon/{name}", getPokemonHandler).Methods("GET")

	fmt.Println("Servidor inicializado em http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

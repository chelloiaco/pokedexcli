package internal

var nextUrl string = "https://pokeapi.co/api/v2/location-area?limit=20"
var previousUrl string = ""
var cacheInterval float32 = 5.0
var cache Cache = Cache{}
var catchRate int = 380
var pokedex map[string]Pokemon = map[string]Pokemon{}
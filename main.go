package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	// "strconv"

	"github.com/joho/godotenv"
	"github.com/shomali11/slacker"
)

type weatherData struct{
	Name string `json:"name"`
	Main struct{
		Kelvin string `json:"temp"`
	}`json:"main"`
}


func main(){
	// Load the environment variable
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// initialize the botapp
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	definition := &slacker.CommandDefinition{Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
		// Getting the User words from request param
		city := request.Param("city")

		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		apiConfig := os.Getenv("OpenWeatherMapApiKey")

		resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?" + "q=" + city + "&appid=" + apiConfig)
		if err != nil{
			log.Fatal(err)
		}
		defer resp.Body.Close()
	
		var data weatherData
		
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			log.Fatal(err)
		}
		
		// return response data to User
		response.Reply(data.Name)
	}}

	bot.Command("{city}", definition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
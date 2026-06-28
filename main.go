package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"geowar-bot/bot"
	"geowar-bot/config"
	"github.com/bwmarrin/discordgo"
)

func main() {
	dg, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		log.Fatal("Failed to create Discord session:", err)
	}

	dg.Identify.Intents = config.Intents

	os.MkdirAll("data", 0755)

	botCore := bot.New(dg)

	dg.AddHandler(botCore.MessageCreate)
	dg.AddHandler(botCore.MessageDelete)
	dg.AddHandler(botCore.MessageUpdate)
	dg.AddHandler(botCore.GuildMemberAdd)
	dg.AddHandler(botCore.GuildMemberRemove)
	dg.AddHandler(botCore.VoiceStateUpdate)
	dg.AddHandler(botCore.InteractionCreate)
	dg.AddHandler(botCore.MessageReactionAdd)

	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		fmt.Printf("Logged in as %s (ID: %s)\n", r.User.Username, r.User.ID)
		fmt.Printf("Prefix: %s\n", config.Prefix)
		fmt.Printf("Serving %d guilds\n", len(r.Guilds))
		s.UpdateCustomStatus(fmt.Sprintf("%shelp | GeoWar SMP", config.Prefix))

		_, err := s.ApplicationCommandCreate(s.State.User.ID, "", &discordgo.ApplicationCommand{
			Name:        "serverstatus",
			Description: "Shows the GeoWar Minecraft server status (updates every minute)",
		})
		if err != nil {
			fmt.Printf("Failed to register slash command: %v\n", err)
		} else {
			fmt.Println("Registered /serverstatus slash command")
		}
	})

	if err := dg.Open(); err != nil {
		log.Fatal("Failed to open connection:", err)
	}
	defer dg.Close()

	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "GeoWar Bot is running")
		})
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		fmt.Printf("Health server listening on :%s\n", port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Printf("Health server error: %v", err)
		}
	}()

	fmt.Println("Bot is running. Press Ctrl+C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc
	fmt.Println("\nShutting down...")
}

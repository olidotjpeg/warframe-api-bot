package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

// Bot parameters
var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	BotToken       = flag.String("token", "", "Bot access token")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)

var s *discordgo.Session

func init() { flag.Parse() }

func init() {
	var err error
	s, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name: "arbitration",
			// All commands and options must have a description
			// Commands/options without description will fail the registration
			// of the command.
			Description: "Get the active Arbitration",
		},
		{
			Name:        "void-trader",
			Description: "Get the active Void Trader",
		},
		{
			Name:        "sortie",
			Description: "Get the current Sortie State",
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"arbitration": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			arbitrationContent := getArbitration()
			arbitrationString := fmt.Sprintf(`
Hello Operator, Ordis have found some Arbitration information for you

Arbitration Mission Type: %s
Star Node %s
Enemy Type %s
`, arbitrationContent.Type, arbitrationContent.Node, arbitrationContent.Enemy)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: arbitrationString,
				},
			})
		},
		"sortie": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			sortieState := getSortieState()

			sortieString := "Hello Operator, Ordis have compiled a list of the current Sortie tasks"

			for index, variant := range sortieState.Variants {
				sortieString = fmt.Sprintf("%s \n Mission #%d is on %s(%s) - the type is %s", sortieString, index+1, variant.Node, variant.Planet, variant.MissionType)
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: sortieString,
				},
			})
		},
		"void-trader": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			voidTraderState := getVoidTraderState()
			voidItems := getVoidItemsDetailed()

			var voidString string

			if !voidTraderState.Active {
				voidString = fmt.Sprintf("Hello Operator, the Void Trader isn't active, please come back in approximately %s", voidTraderState.StartString)
			}

			if voidTraderState.Active {
				voidString = `Hello Operator, the Void trader is selling the following items
				`

				for _, item := range voidItems {
					voidString = fmt.Sprintf("%s \n %s : %s - %s", voidString, item.Name, item.LevelStats[0], item.Category)
				}
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: voidString,
				},
			})
		},
	}
)

func init() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func main() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, *GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	if *RemoveCommands {
		log.Println("Removing commands...")
		// // We need to fetch the commands, since deleting requires the command ID.
		// // We are doing this from the returned commands on line 375, because using
		// // this will delete all the commands, which might not be desirable, so we
		// // are deleting only the commands that we added.
		// registeredCommands, err := s.ApplicationCommands(s.State.User.ID, *GuildID)
		// if err != nil {
		// 	log.Fatalf("Could not fetch registered commands: %v", err)
		// }

		for _, v := range registeredCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, *GuildID, v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

	log.Println("Gracefully shutting down.")
}

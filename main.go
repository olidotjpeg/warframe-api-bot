package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

type Arbitration struct {
	Id          string `json:"id"`
	Activation  string `json:"activation"`
	Expiry      string `json:"expiry"`
	StartString string `json:"startString"`
	Active      bool   `json:"active"`
	Node        string `json:"node"`
	Enemy       string `json:"enemy"`
	EnemyKey    string `json:"enemyKey"`
	Type        string `json:"type"`
	TypeKey     string `json:"typeKey"`
	Archwing    bool   `json:"archwing"`
	Sharkwing   bool   `json:"sharkwing"`
}

type DarvoDeals struct {
	Sold          int    `json:"sold"`
	Item          string `json:"item"`
	Total         int    `json:"total"`
	Eta           string `json:"eta"`
	OriginalPrice int    `json:"originalPrice"`
	SalePrice     int    `json:"salePrice"`
	Discount      int    `json:"discount"`
	Expiry        string `json:"expiry"`
	Id            string `json:"id"`
}

type SortieState struct {
	Id          string
	Activation  string
	Expiry      string
	StartString string
	Active      bool
	RewardPool  string
	Variants    []Variant
	Boss        string
	Faction     string
	FactionKey  string
	Expired     bool
	Eta         string
}

type Variant struct {
	Node                string
	Boss                string
	MissionType         string
	Planet              string
	Modifier            string
	ModifierDescription string
}

type VoidItem struct {
	Item    string
	Ducat   int
	Credits int
}

type VoidTrader struct {
	Id          string
	Activation  string
	StartString string
	Expiry      string
	Active      bool
	Character   string
	Location    string
	Inventory   []VoidItem
	PsId        string
	EndString   string
}

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

			var voidString string

			if !voidTraderState.Active {
				voidString = fmt.Sprintf("Hello Operator, the Void Trader isn't active, please come back in approximately %s", voidTraderState.StartString)
			}

			if voidTraderState.Active {
				voidString = `Hello Operator, the Void trader is selling the following items
				`

				for _, item := range voidTraderState.Inventory {
					voidString = fmt.Sprintf("%s \n %s", voidString, item.Item)
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

func getSortieState() SortieState {
	response, err := http.Get("https://api.warframestat.us/pc/sortie?language=en")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObject SortieState
	json.Unmarshal(responseData, &responseObject)

	return responseObject
}

func getVoidTraderState() VoidTrader {
	response, err := http.Get("https://api.warframestat.us/pc/voidTrader/")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObject VoidTrader
	json.Unmarshal(responseData, &responseObject)

	return responseObject
}

func getArbitration() Arbitration {
	response, err := http.Get("https://api.warframestat.us/pc/arbitration/")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObject Arbitration
	json.Unmarshal(responseData, &responseObject)

	return responseObject
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

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/user"
	"text/tabwriter"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/ehazlett/dialogue/client"
	"github.com/howeyc/gopass"
)

type (
	Configuration struct {
		URL      string `json:"url"`
		Username string `json:"username"`
		Token    string `json:"token"`
	}
)

var (
	URL      string
	USERNAME string
	TOKEN    string
	log      = logrus.New()
	filename string
)

func init() {
	usr, _ := user.Current()
	filename = fmt.Sprintf("%s/.dialogue.cfg", usr.HomeDir)
}

// getConfig returns Configuration from dialogue client config file
func getConfig() (*Configuration, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, err
	}
	file, _ := os.Open(filename)
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	if err := decoder.Decode(&configuration); err != nil {
		log.Warn("Error parsing config file: %s", err)
		return nil, err
	}
	return &configuration, nil
}

// saveConfig saves Configuration to dialogue client config file
func saveConfig(cfg *Configuration) error {
	f, err := os.Create(filename)
	if err != nil {
		log.Errorf("Error creating configuration file: %s", err)
		return err
	}
	b, err := json.Marshal(cfg)
	if err != nil {
		log.Errorf("Error saving configuration: %s", err)
		return nil
	}
	if _, err := io.WriteString(f, string(b)); err != nil {
		return err
	}
	return nil
}

func getTableWriter() *tabwriter.Writer {
	w := tabwriter.NewWriter(os.Stdout, 12, 8, 0, '\t', 0)
	return w
}

func cliLogin(c *cli.Context) {
	var u string
	var user string
	fmt.Printf("URL: ")
	fmt.Scanf("%s", &u)
	fmt.Printf("Username: ")
	fmt.Scanf("%s", &user)
	fmt.Printf("Password: ")
	pass := gopass.GetPasswd()
	token, err := client.Authenticate(u, user, string(pass))
	if err != nil {
		log.Fatalf("Error logging in: %s", err)
	}
	if token == "" {
		log.Fatal("An error occurred while logging in")
	}
	// save config
	cfg := &Configuration{
		URL:      u,
		Username: user,
		Token:    token,
	}
	saveConfig(cfg)
	log.Info("Login successful")
}

func cliDeleteTopic(c *cli.Context) {
	id := c.String("id")
	if id == "" {
		log.Fatal("You must specify a topic ID")
	}
	client, err := client.NewDialogueClient(URL, USERNAME, TOKEN)
	if err != nil {
		log.Fatal(err)
	}
	if err := client.DeleteTopic(id); err != nil {
		log.Fatal(err)
	}
}

func cliCreateTopic(c *cli.Context) {
	title := c.String("title")
	if title == "" {
		log.Fatal("You must specify a title")
	}
	client, err := client.NewDialogueClient(URL, USERNAME, TOKEN)
	if err != nil {
		log.Fatal(err)
	}
	if err := client.CreateTopic(title); err != nil {
		log.Fatal(err)
	}
}

func cliListTopics(c *cli.Context) {
	client, err := client.NewDialogueClient(URL, USERNAME, TOKEN)
	if err != nil {
		log.Fatal(err)
	}
	topics, err := client.GetTopics()
	if err != nil {
		log.Fatal(err)
	}
	if len(topics) == 0 {
		return
	}
	w := getTableWriter()
	fmt.Fprint(w, "Title\tID\t\n")
	for _, t := range topics {
		fmt.Fprintf(w, "%s\t%v\t\n", t.Title, t.Id)
	}
	w.Flush()
}

func cliCreatePost(c *cli.Context) {
	content := c.String("content")
	topicId := c.String("topicId")
	if content == "" || topicId == "" {
		log.Fatal("You must specify a topic id and content")
	}
	client, err := client.NewDialogueClient(URL, USERNAME, TOKEN)
	if err != nil {
		log.Fatal(err)
	}
	if err := client.CreatePost(topicId, content); err != nil {
		log.Fatal(err)
	}
}

func cliListPosts(c *cli.Context) {
	topicId := c.String("topicId")
	showIds := c.Bool("ids")
	if topicId == "" {
		log.Fatal("You must specify a topic id")
	}
	client, err := client.NewDialogueClient(URL, USERNAME, TOKEN)
	if err != nil {
		log.Fatal(err)
	}
	posts, err := client.GetPosts(topicId)
	if err != nil {
		log.Fatal(err)
	}
	if len(posts) == 0 {
		return
	}
	w := getTableWriter()
	for _, p := range posts {
		fmt.Fprintf(w, "%v\t -%s", p.Content, p.Author)
		if showIds {
			fmt.Fprintf(w, "\t%s", p.Id)
		}
		fmt.Fprint(w, "\n")
	}
	w.Flush()
}

func cliDeletePost(c *cli.Context) {
	id := c.String("id")
	if id == "" {
		log.Fatal("You must specify a post ID")
	}
	client, err := client.NewDialogueClient(URL, USERNAME, TOKEN)
	if err != nil {
		log.Fatal(err)
	}
	if err := client.DeletePost(id); err != nil {
		log.Fatal(err)
	}
}

func main() {
	// load config
	config, _ := getConfig()
	if config != nil {
		URL = config.URL
		USERNAME = config.Username
		TOKEN = config.Token
	}
	app := cli.NewApp()
	app.Name = "dialogue"
	app.Version = "0.0.1"
	// commands
	app.Commands = []cli.Command{
		{
			Name:      "login",
			ShortName: "l",
			Usage:     "Login to Dialogue",
			Action:    cliLogin,
		},
		{
			Name:      "topics",
			ShortName: "t",
			Usage:     "Topic Commands",
			Subcommands: []cli.Command{
				{
					Name:      "create",
					ShortName: "c",
					Usage:     "create a new topic",
					Action:    cliCreateTopic,
					Flags: []cli.Flag{
						cli.StringFlag{"title, t", "", "Topic title"},
					},
				},
				{
					Name:      "delete",
					ShortName: "d",
					Usage:     "delete a topic",
					Action:    cliDeleteTopic,
					Flags: []cli.Flag{
						cli.StringFlag{"id, i", "", "Topic ID"},
					},
				},
				{
					Name:      "list",
					ShortName: "l",
					Usage:     "list topics",
					Action:    cliListTopics,
				},
			},
		},
		{
			Name:      "posts",
			ShortName: "p",
			Usage:     "Post Commands",
			Subcommands: []cli.Command{
				{
					Name:      "create",
					ShortName: "c",
					Usage:     "create a new post",
					Action:    cliCreatePost,
					Flags: []cli.Flag{
						cli.StringFlag{"topicId, i", "", "Topic ID"},
						cli.StringFlag{"content, c", "", "Post content"},
					},
				},
				{
					Name:      "delete",
					ShortName: "d",
					Usage:     "delete a post",
					Action:    cliDeletePost,
					Flags: []cli.Flag{
						cli.StringFlag{"id, i", "", "Post ID"},
					},
				},
				{
					Name:      "list",
					ShortName: "l",
					Usage:     "list posts",
					Action:    cliListPosts,
					Flags: []cli.Flag{
						cli.StringFlag{"topicId, i", "", "Topic ID"},
						cli.BoolFlag{"ids", "Show post ids"},
					},
				},
			},
		},
	}
	// run
	app.Run(os.Args)
}

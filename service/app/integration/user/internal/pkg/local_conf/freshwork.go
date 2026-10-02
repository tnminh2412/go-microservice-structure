package local_conf

type ChatClient struct {
	Token  string `bson:"token"`
	Host   string `bson:"host"`
	SdkUrl string `bson:"sdkUrl"`
}

type FreshWorkApp struct {
	ClientID     string `bson:"clientID"`
	LoginPageUri string `bson:"loginPageUri"`
	Name         string `bson:"name"`
	Description  string `bson:"description"`
}

type FreshWork struct {
	PlayerApp *App   `bson:"playerApp"`
	AgentID   string `bson:"agentID"`
	Enable    bool   `bson:"enable"`
}

type Chat struct {
	FreshWork *FreshWork `bson:"freshWork"`
	Jwt       *Jwt       `bson:"jwt"`
}

type App struct {
	Ticket      *FreshWorkApp `bson:"ticket"`
	ChatSupport *FreshWorkApp `bson:"chatSupport"`
	ChatClient  *ChatClient   `bson:"chatClient"`
}

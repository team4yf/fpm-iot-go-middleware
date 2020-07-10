//Package light the light controller factory
//it containes the keys of the apps
//use the brand + appid, we cant get the rest client/controller
//use the sync.pool to gen the client
package light

import (
	"errors"
	"sync"

	"github.com/team4yf/fpm-iot-go-middleware/external/device/light/lt10"
	"github.com/team4yf/fpm-iot-go-middleware/external/rest"
	"github.com/team4yf/fpm-iot-go-middleware/internal/service"
)

var (
	clientService service.ClientService
	clientConfigs map[string]*rest.Options
	restClients   map[string]rest.Client
	lock          sync.Mutex
)

//Init load the configs of the clients
func Init() error {
	clientService = service.NewSimpleClientService()
	clients, err := clientService.ListByCondition("type = ? and status = 1", "light")
	if err != nil {
		return err
	}
	clientConfigs = make(map[string]*rest.Options)
	restClients = make(map[string]rest.Client)
	for _, client := range clients {
		clientConfigs[client.Brand+"/"+client.AppKey] = &rest.Options{
			AppID:       client.AppKey,
			AppSecret:   client.SecretKey,
			Username:    client.Username,
			TokenExpire: client.Expired,
			Enviroment:  client.Environment,
			BaseURL:     client.APIBaseURL,
		}
	}

	return nil
}

//NewAPIClient create a new api client with the brand & appid
func NewAPIClient(brand, appid string) (rest.Client, error) {
	lock.Lock()
	defer lock.Unlock()
	key := brand + "/" + appid
	var client rest.Client
	client, ok := restClients[key]
	if ok {
		return client, nil
	}

	switch brand {
	case "lt10":
		client = lt10.NewClient(clientConfigs[key], clientService.GetCache())

	default:
		return nil, errors.New("not support brand: " + brand)
	}
	if err := client.Init(); err != nil {
		return nil, err
	}
	restClients[key] = client
	return client, nil
}

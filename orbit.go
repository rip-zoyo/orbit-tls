package orbit

import "github.com/rip-zoyo/orbit-tls/client"

type Client = client.Client
type Response = client.Response

var Chrome120 = client.Chrome120
var Chrome131 = client.Chrome131
var Chrome138 = client.Chrome138

var Firefox121 = client.Firefox121
var Firefox131 = client.Firefox131

var Safari17 = client.Safari17
var Safari18 = client.Safari18
var SafariiOS = client.SafariiOS

var Edge120 = client.Edge120

var MullvadBrowser = client.MullvadBrowser
var ChromeAndroid = client.ChromeAndroid
var Opera115 = client.Opera115
var Brave131 = client.Brave131
var Brave138 = client.Brave138

func New(profileName string) (*Client, error) {
	return client.New(profileName)
} 
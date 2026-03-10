package conf

import "net/http"

type Conf struct{
	Port      string
	Signature string
	HttpClient    *http.Client
}


var Config = Conf{}
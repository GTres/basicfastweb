package main


import(
	"fmt"
	"log"

	// pakemon
	"encoding/json"
	// "io/ioutil"
	// "net/http"
	"os/exec"
	

	//fasthttp

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)


type Server struct {
	Host				string		`json:"host,omitempty"`
	Server_Changed		bool		`json:"servers_changed,omitempty"`
	Ssl_Grade			string		`json:"ssl_grade,omitempty"`
	Previous_Ssl_Grade	string		`json:"previous_ssl_grade,omitempty"`
	Logo				string		`json:"logo,omitempty"`
	Title				string		`json:"title,omitempty"`
	Is_Down				bool		`json:"is_down,omitempty"`
	Endpoints			[]Servers	`json:"endpoints"`
}

type Servers struct {
	Address				string		`json:"ipAddress"`
	Ssl_grade			string		`json:"grade"`
	Country				string		`json:"country"`
	Owner				string		`json:"owner"`
}



func whois(ctx *fasthttp.RequestCtx) {

	app := "whois"
	arg := fmt.Sprintf("%v", ctx.UserValue("name"))

	cmd := exec.Command(app, arg)
	stdout, err :=  cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	
	fmt.Fprintf(ctx, "out:\n%s", stdout)
}

func doRequest(ctx *fasthttp.RequestCtx) {
	ssllab  := "https://api.ssllabs.com/api/v3/analyze?host="
	url := fmt.Sprintf("%v%v", ssllab, ctx.UserValue("name"))

	// output := fmt.Sprintf("%s%s","https://api.ssllabs.com/api/v3/analyze?host=", url)

	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)

	resp := fasthttp.AcquireResponse()

	err := fasthttp.Do(req, resp)

	if err != nil {
		fmt.Print(ctx, err.Error())
		// os.Exit(1)
	}
	
	// query := ctx.QueryArgs().Peek(serverName)

	bodyBytes := resp.Body()

	var s Server
	json.Unmarshal(bodyBytes, &s)

	// for i := 0; i < len(responseObject); i++ {
		// fmt.Fprintln(ctx, responseObject.host)
	// }
	// fmt.Println("\nResponse")
	fmt.Fprintln(ctx, s)
	fmt.Fprintln(ctx, "Host:", s.Host)
	fmt.Fprintln(ctx, "Server changed:", s.Server_Changed)
	fmt.Fprintln(ctx, "ssl grade:", s.Ssl_Grade)
	fmt.Fprintln(ctx, "Previous ssl grade", s.Previous_Ssl_Grade)
	fmt.Fprintln(ctx, "Logo:", s.Logo)
	fmt.Fprintln(ctx, "Title:", s.Title)
	fmt.Fprintln(ctx, "Is down:", s.Is_Down)

	fmt.Fprintln(ctx, s.Endpoints)
	fmt.Fprintln(ctx, "Endpoints:", len(s.Endpoints))
	for i := 0; i < len(s.Endpoints); i++ {
		fmt.Fprintln(ctx, "\nEndpoint: %d",i+1)
		fmt.Fprintln(ctx, "Address:", 	s.Endpoints[i].Address)
		fmt.Fprintln(ctx, "Ssl grade:", s.Endpoints[i].Ssl_grade)
		fmt.Fprintln(ctx, "Country:", 	s.Endpoints[i].Country)
		fmt.Fprintln(ctx, "Owner:", 	s.Endpoints[i].Owner)
	}

	// fmt.Println("Host: %s",responseObject[0])


	// var s interface{}
	// // var s Server
	// b := resp.Body()
	// json.Unmarshal(b, &s)
	


	// m := s.(map[string]interface{})
	// // m["servers"].(map[string]interface{})[]
	// fmt.Fprintln(ctx, m)
	// fmt.Fprintln(ctx, "Host:", m["host"])
	// fmt.Fprintln(ctx, "Protocol:", m["protocol"])
	// fmt.Fprintln(ctx, "Servers_changed:", m["servers_changed"])
	// fmt.Fprintln(ctx, "Ssl_grade:", m["ssl"])
	// fmt.Fprintln(ctx, "Previous_ssl_grade:", m["previous_ssl_grade"])
	// fmt.Fprintln(ctx, "Logo:", m["logo"])
	// fmt.Fprintln(ctx, "Title:", m["title"])
	// fmt.Fprintln(ctx, "Down:", m["is_down"])
	// fmt.Fprintln(ctx, "ssl_grade:", m["endpoints"].(map[string]interface{})["grade"])


	// for k,v := range  {
	// 	fmt.Println()
	// 	switch vv := v.(type){
	// 	case string:
	// 		fmt.Println(k, "is string", vv)
	// 	case int:
	// 		fmt.Println(k, "is int", vv)
	// 	case []interface{}:
	// 		fmt.Println(k, "is an array:")
	// 		for i, u := range vv {
	// 			fmt.Println(i, u)
	// 		}
	// 	default:
	// 		fmt.Println(k, "is of a type I don't know how to handble")
	// 	}
	// }
}

func Index(ctx *fasthttp.RequestCtx) {
	fmt.Fprintln(ctx, "Welcom!")
}

func Hello(ctx *fasthttp.RequestCtx) {
	// fmt.Println(str)
	// fmt.Fprintf(ctx, str)
	// fmt.Fprintf(ctx, "hello, %s!\n", ps.ByName("name"))
	// fmt.Fprintf(ctx, "hello, %s!\n", ctx.UserValue("name"))
	// u := fmt.Sprintf(ctx.UserValue("serverName"))
	// doRequest(ctx, "https://pokeapi.co/api/v2/pokedex/kanto/")
	// doRequest(ctx, "https://api.ssllabs.com/api/v3/analyze?host=google.com")
	// doRequest(ctx, ps.ByName("serverName"))
	doRequest(ctx)
	whois(ctx)
}

func main() {
	
	

	router := fasthttprouter.New()
	router.GET("/", Index)
	router.GET("/server/:name", Hello)

	log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))
}
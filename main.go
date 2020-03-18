package main


import(
	"fmt"
	"log"

	// pakemon
	"encoding/json"
	// "io/ioutil"
	// "net/http"
	"os/exec"
	"regexp"
	

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
	EndPoints			[]EndPoint	`json:"endpoints"`
}

type EndPoint struct {
	Address				string		`json:"ipAddress"`
	Ssl_Grade			string		`json:"grade"`
	Country				string		`json:"country"`
	Owner				string		`json:"owner"`
}


func Whois(address string) (owner string, country string){
	whoisstring := ""
	
	expressionhost, err := regexp.Compile(`(https?://)?(www\.)?(.*)`)
	if err != nil {
		fmt.Println(err.Error())
	}

	argumentcut := expressionhost.FindAllStringSubmatch(string(address), -1)
	if len(argumentcut) > 0 {
		// host := s.Address // Usar la direccion de los endpoints
		host := argumentcut[0][len(argumentcut[0])-1]
		app := "whois"
		wohiscommand := exec.Command(app, host)

		whoisoutput, err :=  wohiscommand.Output()
		if err != nil {
			fmt.Println(err.Error())
		}
		whoisstring = string(whoisoutput)
		if whoisstring != "" {
			owner = FindInformation(whoisstring, `(?P<t>[oO]rg-?[nN]ame:)(?P<s>\s*)(?P<o>.*)`, "o")
			country = FindInformation(whoisstring, `(?P<t>[cC]ountry:)(?P<s>\s+)(?P<c>[A-Z]{2})`, "c")
		}
	}
	return owner, country
}

// Bucar en un string la expresion regular indicada.
func FindInformation(str, expression, group string) (information string) {
		compileexpression, err := regexp.Compile(expression)
		if err != nil {
			fmt.Println(err.Error())
		}
		
		match := compileexpression.FindStringSubmatch(string(str))
		result := make(map[string]string)
		subname := compileexpression.SubexpNames()
		if len(match) > 0 {
			for i, name := range subname {
				if i != 0 && name != "" {
					result[name] = match[i]
				}
			}
			if len(subname)>0 {
				information = result[group]
			}
		}
		return information
}


// Buscar el logo de pagina de host indicado en el request.
func FindLogo(ctx *fasthttp.RequestCtx) (logo string) {
	addresshost := ctx.QueryArgs().Peek("address")
	host := fmt.Sprintf("%v", string(addresshost))

	if len(host) > 0 {		
		curlcommand := exec.Command("curl", host)
		html, err :=  curlcommand.Output()
		if err != nil {
			fmt.Println(err.Error())
		}

		compileexpression, err := regexp.Compile(`https?\S*(logo)?\S*\.png`)
		if err != nil {
			fmt.Println(err.Error())
		}
		
		logo = compileexpression.FindString(string(html))
		fmt.Println("Logo:",logo)
	}
	return logo
}

// Procesar el request.
func doRequest(ctx *fasthttp.RequestCtx) {
	ssllab  := "https://api.ssllabs.com/api/v3/analyze?host="
	a := ctx.QueryArgs().Peek("address")
	h := fmt.Sprintf("%v", string(a))
	url := fmt.Sprintf("%v%v", ssllab, h)


	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	resp := fasthttp.AcquireResponse()

	err := fasthttp.Do(req, resp)
	if err != nil {
		fmt.Print(ctx, err.Error())
	}
	
	bodyBytes := resp.Body()

	var s Server
	json.Unmarshal(bodyBytes, &s)

	// fmt.Println("\nResponse")
	s.Logo = FindLogo(ctx)
	fmt.Fprintln(ctx, "logo listo")
	fmt.Fprintln(ctx, s)
	fmt.Fprintln(ctx, "Host:", s.Host)
	fmt.Fprintln(ctx, "Server changed:", s.Server_Changed)
	fmt.Fprintln(ctx, "ssl grade:", s.Ssl_Grade)
	fmt.Fprintln(ctx, "Previous ssl grade", s.Previous_Ssl_Grade)
	fmt.Fprintln(ctx, "Logo:", s.Logo)
	fmt.Fprintln(ctx, "Title:", s.Title)
	fmt.Fprintln(ctx, "Is down:", s.Is_Down)
	

	sj := make(map[string]interface{})
	sj["Host"] = s.Host
	

	fmt.Fprintln(ctx, s.EndPoints)
	fmt.Fprintln(ctx, "EndPoints:", len(s.EndPoints))
	// type Points []EndPoint
	var points []EndPoint
	for i := 0; i < len(s.EndPoints); i++ {
		// s.EndPoints[i].FindOwner(ctx)
		// s.EndPoints[i].FindCountry(ctx)
		s.EndPoints[i].Owner, s.EndPoints[i].Country = Whois(s.EndPoints[i].Address)
		fmt.Fprintln(ctx, "\nEndpoint: %d",i+1)
		fmt.Fprintln(ctx, "Address:", 	s.EndPoints[i].Address)
		fmt.Fprintln(ctx, "Ssl grade:", s.EndPoints[i].Ssl_Grade)
		fmt.Fprintln(ctx, "Country:", 	s.EndPoints[i].Country)
		fmt.Fprintln(ctx, "Owner:", 	s.EndPoints[i].Owner)
		// response.EndPoint[i]{
		// 	Address:s.Endpoints[i].Address,
		// 	Ssl_Grade:s.Endpoints[i].Ssl_Grade,
		// 	Country:s.Endpoints[i].Country,
		// 	Owner:s.Endpoints[i].Owner
		// }
		// sj["Endpoints"][i] = map[string]interface{}[]{
		// 	"Address": s.EndPoints[i].Address,
		// 	"Country": s.EndPoints[i].Country,
		// }

		points[i] = EndPoint {
				Address: s.EndPoints[i].Address,
				Ssl_Grade: s.EndPoints[i].Ssl_Grade,
				Country: s.EndPoints[i].Country,
				Owner: s.EndPoints[i].Owner}
	}

	response := Server{
		Host:s.Host, 
		Server_Changed:s.Server_Changed,
		Ssl_Grade:s.Ssl_Grade,
		Previous_Ssl_Grade:s.Previous_Ssl_Grade,
		Logo:s.Logo,
		Title:s.Title,
		Is_Down:s.Is_Down,
		EndPoints:points}

	datatoresponse, err := json.MarshalIndent(response, "", "    ")
	if err != nil {
		fmt.Print(ctx, err.Error())
	}
	// d := fmt.Sprintln(datatoresponse)
	fmt.Println("%s\n", string(datatoresponse))

	

	type Employee struct {
		Name string
		Age int
		Salary int
	}
	emp_obj := Employee{Name:"Rachel", Age:24, Salary :344444}
	emp, _ := json.Marshal(emp_obj)
	fmt.Println(string(emp))
}

func Index(ctx *fasthttp.RequestCtx) {
	fmt.Fprintln(ctx, "Welcom!")
}

func Hello(ctx *fasthttp.RequestCtx) {
	doRequest(ctx)
}

func (s Server) L(ctx *fasthttp.RequestCtx) {
	h := ctx.QueryArgs().Peek("address")
	// fmt.Fprintln(ctx, string(ad))


	app := "curl"
	arg := fmt.Sprintf("%v", string(h))
	cmd := exec.Command(app, arg)
	exec, _ :=  cmd.Output()
	// fmt.Fprintln(ctx, string(logo))

	regLogo, err := regexp.Compile(`http[s?]\S*[logo?]\S*\.png`)
	if err != nil {
		fmt.Println(err.Error())
	}
	//strLogo, err := string(regCountry.FindAllStringSubmatch(string(logo), -1)[0][0])
	findLogo  := regLogo.FindAllStringSubmatch(string(exec), -1)
	fmt.Fprintln(ctx, len(findLogo))
	if len(findLogo) > 0 {
		s.Logo = string(findLogo[0][0])
	}
	if err!= nil {
		fmt.Println(err.Error())
	}
	// fmt.Fprintln(ctx, string(strLogo))
	
}

func main() {
	
	

	router := fasthttprouter.New()
	router.GET("/", Index)
	router.GET("/server/:name", Hello)
	router.GET("/par/est", Hello)

	log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))
}
package main

import(
	"fmt"
	"encoding/json"
	"os/exec"
	"regexp"

	"github.com/valyala/fasthttp"
)


func Whois(address string) (owner string, country string){
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

		whoisstring := string(whoisoutput)
		if whoisstring != "" {
			owner = FindInformation(whoisstring, `(?P<t>[oO]rg-?[nN]ame:)(?P<s>\s*)(?P<o>.*)`, "o")
			country = FindInformation(whoisstring, `(?P<t>[cC]ountry:)(?P<s>\s+)(?P<c>[A-Z]{2})`, "c")
		}
	}

	return owner, country
}


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


func LookLogoAndTitleInHTML(ctx *fasthttp.RequestCtx) (logo string, title string) {
	addresshost := ctx.QueryArgs().Peek("address")
	host := fmt.Sprintf("%v", string(addresshost))

	if len(host) > 0 {		
		curlcommand := exec.Command("curl", host)
		html, err :=  curlcommand.Output()
		if err != nil {
			fmt.Println(err.Error())
		}

		logoexpression, err := regexp.Compile(`https?\S*(logo)?\S*\.png`)
		if err != nil {
			fmt.Println(err.Error())
		}

		// titleexpression, err := regexp.Compile(`(?P<t1>\<title\>)(?P<t2>.*)(?P<t3>\</title\>)`)
		titleexpression, err := regexp.Compile(`(?P<t1><title>)(?P<t2>.*)(?P<t3></title>)`)
		if err != nil {
			fmt.Println(err.Error())
		}
		
		logo = logoexpression.FindString(string(html))


		match := titleexpression.FindStringSubmatch(string(html))
		result := make(map[string]string)
		subname := titleexpression.SubexpNames()

		if len(match) > 0 {
			for i, name := range subname {
				if i != 0 && name != "" {
					result[name] = match[i]
				}
			}
			if len(subname)>0 {
				title = result["t2"]
			}
		}
	}
	return logo, title
}




func doRequest(ctx *fasthttp.RequestCtx) (Server) {
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

	s.Logo, s.Title = LookLogoAndTitleInHTML(ctx)
	for i := 0; i < len(s.EndPoints); i++ {
		s.EndPoints[i].Owner, s.EndPoints[i].Country = Whois(s.EndPoints[i].Address)
	}

	return s
}
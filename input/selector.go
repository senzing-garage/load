package input

import (
	"context"
	"fmt"
	"net/url"

	"github.com/senzing/load/input/rabbitmq"
	"github.com/senzing/load/input/sqs"
)

// load is 6201:  https://github.com/Senzing/knowledge-base/blob/main/lists/senzing-product-ids.md
const productIdentifier = 6201

var idMessages = map[int]string{
	21: "Scheme: %s",
	22: "User full: %s",
	23: "User name: %s",
	24: "Password: %s",
	25: "Host full: %s",
	26: "Host: %s",
	27: "Port: %s",
	28: "Path: %s",
	29: "Fragment: %s",
	30: "Query string: %s",
	31: "Parsed query string: %s",
}

// ----------------------------------------------------------------------------
func parseURL(urlString string) *url.URL {
	fmt.Println("Parse url:", urlString)
	u, err := url.Parse(urlString)
	if err != nil {
		panic(err)
	}

	// msglog.Log(21, u.Scheme, messagelogger.LevelInfo)
	// msglog.Log(22, u.User, messagelogger.LevelInfo)
	// msglog.Log(23, u.User.Username(), messagelogger.LevelInfo)
	// fmt.Println("Scheme:", u.Scheme)
	// fmt.Println("User:", u.User)
	// fmt.Println("User.Username():", u.User.Username())
	// p, _ := u.User.Password()
	// fmt.Println("User.Password():", p)
	// if len(p) > 0 {
	// 	msglog.Log(24, "SET, redacted from logs", messagelogger.LevelInfo)
	// } else {
	// 	msglog.Log(24, "NOT SET", messagelogger.LevelInfo)
	// }

	// msglog.Log(25, u.Host, messagelogger.LevelInfo)
	// host, port, _ := net.SplitHostPort(u.Host)
	// msglog.Log(26, host, messagelogger.LevelInfo)
	// msglog.Log(27, port, messagelogger.LevelInfo)

	// msglog.Log(28, u.Path, messagelogger.LevelInfo)
	// msglog.Log(29, u.Fragment, messagelogger.LevelInfo)

	// msglog.Log(30, u.RawQuery, messagelogger.LevelInfo)
	// m, _ := url.ParseQuery(u.RawQuery)
	// msglog.Log(31, m, messagelogger.LevelInfo)
	// fmt.Println("Query:", m)

	return u
}

// ----------------------------------------------------------------------------
func Read(ctx context.Context, inputURL, engineConfigJson string, engineLogLevel, numberOfWorkers, visibilityPeriodInSeconds int, logLevel string, jsonOutput bool) error {
	// if len(logLevel) > 0 {
	// 	msglog.SetLogLevelFromString(logLevel)
	// }

	u := parseURL(inputURL)
	if len(inputURL) <= 0 {
		return fmt.Errorf("Invalid URL: %s", inputURL)
	}
	switch u.Scheme {
	case "amqp":
		rabbitmq.Read(ctx, inputURL, engineConfigJson, logLevel, jsonOutput)
	case "sqs":
		//allows for using a dummy URL with just a queue-name
		// eg  sqs://lookup?queue-name=myqueue
		sqs.Read(ctx, inputURL, engineConfigJson, engineLogLevel, numberOfWorkers, visibilityPeriodInSeconds, logLevel, jsonOutput)
	case "https":
		//uses actual AWS SQS URL.  TODO: detect sqs/amazonaws url?
		sqs.Read(ctx, inputURL, engineConfigJson, engineLogLevel, numberOfWorkers, visibilityPeriodInSeconds, logLevel, jsonOutput)
	default:
		// msglog.Log(2001, u.Scheme, messagelogger.LevelWarn)
	}
	return nil
}

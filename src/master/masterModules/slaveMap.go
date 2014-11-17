package masterModule

import (
	"net/http"
	"fmt"
	"strings"
	"net/url"
	"time"
)

// var slaveIPMap = make(map[string]string)
var slaveIPMap = initializeSlaveIPs()
var slaveHeartbeatMap = make(map[string]time.Time) // TODO: make these map to time values

func SetUp() (slaveMap map[string]string) {
	return slaveIPMap
}

func ReceiveAndMapSlaveAddress(_ http.ResponseWriter, request *http.Request) {
	slaveName := request.PostFormValue("slaveName")
	slaveIPAddress := request.PostFormValue("slaveIPAddress")
	fmt.Printf("\nNEW SLAVE RECEIVED.\n")
	fmt.Println("Slave Name: ", slaveName)
	fmt.Println("Slave IP address: ", slaveIPAddress)

	if returnedIPAddress, existsInMap := slaveIPMap[slaveName]; existsInMap == false {
		webserverIPAddressAndExtentionArray := []string{"http://localhost:4003", "/receive_slave"}

		err := sendSlaveToWebserver(webserverIPAddressAndExtentionArray, slaveName)
		printServerResponse(err, slaveName)
	} else {
		fmt.Printf("WARNING: Slave with name \"%v\" already exists with the IP address: %v. \nUpdating %v's IP address to %v.\n", slaveName, returnedIPAddress, slaveName, slaveIPAddress)
	}
	slaveIPMap[slaveName] = slaveIPAddress
	slaveHeartbeatMap[slaveName] = time.Now()
	fmt.Printf("Mapped \"%v\" to %v.\n", slaveName, slaveIPAddress)
	fmt.Println("Valid slave IDs are: ", slaveIPMap)
}

func MonitorSlaveHeartbeats(_ http.ResponseWriter, request *http.Request) {
	slaveName := request.PostFormValue("slaveName")
	heartbeatTimestamp := request.PostFormValue("heartbeatTimestamp")

	timeFormat := "2006-01-02 15:04:05.999999999 -0700 MST"
	heartbeatTime, err := time.Parse(timeFormat, heartbeatTimestamp)

	if err != nil {
		fmt.Println("Error encountered when parsing heartbeat timestamp from slave.")
		fmt.Println("ERROR: ", err)
	}

	slaveHeartbeatMap[slaveName] = heartbeatTime

}

func MonitorSlaves(timeInterval int) {
	timer := time.Tick(time.Duration(timeInterval) * time.Second)
    
    for _ = range timer {
		removeDeadSlaves(timeInterval)
    }
}

func removeDeadSlaves(deadTime int) {
	for slaveName, lastHeartbeatTime := range slaveHeartbeatMap {
		if time.Now().Sub(lastHeartbeatTime) > time.Duration(deadTime) * time.Second {
			fmt.Println("REMOVING DEAD SLAVE: ", slaveName)
			delete(slaveHeartbeatMap, slaveName)
			delete(slaveIPMap, slaveName)
			fmt.Println("Updated Slave Map: ", slaveIPMap)
			fmt.Printf("\n\n")

			// Need to delete slave from webserver, too? 
			// I think it's better to keep centralized slave list.
			// Webserver will fetch entire list from master with each refresh.
		}
	}
}

func sendSlaveToWebserver(webserverIPAddressAndExtentionArray []string, slaveName string) (err error) {
	client := &http.Client{}
	webserverReceiveSlaveAddress := strings.Join(webserverIPAddressAndExtentionArray, "")

	form := url.Values{}
	form.Set("slaveName", slaveName)
	_, err = client.PostForm(webserverReceiveSlaveAddress, form)

	printServerResponse(err, slaveName)

	return
}

func printServerResponse(error error, slaveName string) {
	if error != nil {
		fmt.Printf("Error communicating with webserver: %v\n", error)
		fmt.Printf("%v not updated on webserver.\n", slaveName)
	} else {
		fmt.Printf("Added \"%v\" to webserver slave list.\n", slaveName)
	}
}

func initializeSlaveIPs() (slaveIPMap map[string]string) {
	slaveIPs := make(map[string]string)
	slaveIPs["1"] = "http://10.0.0.122:8080"
	slaveIPs["2"] = "http://10.0.1.11:8080"

	return slaveIPs
}


package main

import (
	"flag"
	"fmt"
	"sync"

	"github.com/g0194776/tinydhcp-dockerip/providers"
	log "github.com/sirupsen/logrus"

	iris "gopkg.in/kataras/iris.v6"
	"gopkg.in/kataras/iris.v6/adaptors/httprouter"
)

var (
	provider providers.DataProvider
	lock     *sync.Mutex
)

func main() {
	listenPort := 8080
	ipRange := flag.String("i", "", "A CIDR IP range used for initializing entire IP pool.")
	envId := flag.Int("e", 255, "Business evironment ID")
	providerStr := flag.String("p", "mysql", "A storage type that where the data store.")
	needInitData := flag.Bool("n", false, "If it be set as true, the tiny-dhcp will automatically initializes IP pool at first.")
	mysqlConnectionStr := flag.String("mysql", "", "If you have decided using mysql as backend storage, this parameter MUST be a MYSQL connection string")
	flag.Parse()
	var err error
	provider, err = providers.GetDataProvider(*providerStr)
	if err != nil {
		log.Fatalf("Error occured while creating a data provider by you passed argument: provider (%s), error: %s", *providerStr, err.Error())
	}
	if *envId == 255 && *needInitData {
		log.Fatal("You have to typed another one real BIZ environment ID rather than 255.")
	}
	log.Info("Start initializing data provider...")
	ipGenerator := providers.IPCIDRGenerator{}
	var ips []string = nil
	if *needInitData {
		log.Info("Calculating IP pool...")
		ips, err = ipGenerator.Generate(*ipRange)
		if err != nil {
			log.Fatalf("Failed calculating IP range for given base IP: %s", *ipRange)
		}
	}
	err = provider.Initialize(ips, *envId, *mysqlConnectionStr, *needInitData)
	if err != nil {
		log.Fatalf("Error occured while initializing data provider, error: %s", err.Error())
	}
	lock = &sync.Mutex{}
	log.Info("Initializing Web Engine...")
	app := iris.New()
	app.Adapt(httprouter.New())
	log.Info("Registering Web Router...")
	//REG HTTP routing rules.
	app.Get("/ip", WorkProc)
	log.Infof("Registering local TCP network port: %d", listenPort)
	go app.Listen(fmt.Sprintf(":%d", listenPort))
	log.Info("Tiny DHCP server has been started!")
	select {}
}

func WorkProc(ctx *iris.Context) {
	lock.Lock()
	defer lock.Unlock()
	nodeIp := ctx.URLParam("node-ip")
	owner := ctx.URLParam("owner")
	desc := ctx.URLParam("desc")
	rsp := &HttpResponse{}
	envId, err := ctx.URLParamInt("env-id")
	if err != nil {
		rsp.ErrorID = 255
		rsp.Reason = "\"env-id\" MUST be set correctly in the HTTP query parameters."
		ctx.JSON(400, rsp)
		return
	}
	ip, err := provider.GetAvailableIP(nodeIp, envId, owner, desc)
	if err != nil {
		rsp.ErrorID = 255
		rsp.Reason = err.Error()
		ctx.JSON(400, rsp)
		return
	} else {
		rsp.ErrorID = 0
		rsp.DockerIP = ip
		ctx.JSON(iris.StatusOK, rsp)
	}
}

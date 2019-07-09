/*
 * NSSF Service
 */

package nssf_service

import (
    "github.com/gin-gonic/gin"

    "free5gc-nssf/factory"
    "free5gc-nssf/flog"
    "free5gc-nssf/nsselection"
    "free5gc-nssf/nssaiavailability"
)

type Nssf struct {
    Cfg string
}

func (n *Nssf) Initialize() {
    factory.InitConfigFactory(n.Cfg)
    flog.System.Debugf("Use configuration %s", n.Cfg)
}

func (n *Nssf) Start() {
    flog.System.Infof("Server started")

    // Running in "release" mode instead of "debug" mode
    gin.SetMode(gin.ReleaseMode)
    router := gin.Default()

    nsselection.AddService(router)
    nssaiavailability.AddService(router)

    flog.System.Fatal(router.Run(":8080"))
}

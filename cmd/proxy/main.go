package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	agentcfg "harmonycloud.cn/stellaris/pkg/agent/config"
	"harmonycloud.cn/stellaris/pkg/agent/handler"
	"time"
)

var (
	heartbeatPeriod time.Duration
	coreAddress     string
	clusterName     string
	addonPath       string
)

func init() {
	flag.DurationVar(&heartbeatPeriod, "heartbeat-send-period", 30, "The period of heartbeat send interval")
	flag.StringVar(&coreAddress, "core-address", "", "address of stellaris")
	flag.StringVar(&clusterName, "cluster-name", "", "name of agent-cluster")
	flag.StringVar(&addonPath, "addon-path", "", "path of addon config")
}
func main() {
	flag.Parse()

	cfg := agentcfg.DefaultConfiguration()
	cfg.HeartbeatPeriod = heartbeatPeriod
	cfg.ClusterName = clusterName
	cfg.CoreAddress = coreAddress
	cfg.AddonPath = addonPath
	err := handler.Register(cfg)
	if err != nil {
		logrus.Fatalf("failed register cluster: %s", err)
	}

}

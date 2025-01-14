package addons

import (
	"github.com/sirupsen/logrus"
	"harmonycloud.cn/stellaris/config"
	agentconfig "harmonycloud.cn/stellaris/pkg/agent/config"
	"harmonycloud.cn/stellaris/pkg/model"
	"harmonycloud.cn/stellaris/pkg/util/common"
	"reflect"
	"time"
)

const (
	forceSynchronization = 30
	timeOut              = 10
	HeartbeatMessage     = "ok"
)

func Load(cfg *model.PluginsConfig) (*model.RegisterRequest, *model.AddonsChannel) {
	var channels []chan *model.Addon
	var monitorChannels []chan *model.Condition
	inTreeLen := len(cfg.Plugins.InTree)
	outTreeLen := len(cfg.Plugins.OutTree)
	if inTreeLen <= 0 && outTreeLen <= 0 {
		return &model.RegisterRequest{}, &model.AddonsChannel{}
	}
	if inTreeLen > 0 {
		inTreePlugins := cfg.Plugins.InTree
		for _, name := range inTreePlugins {
			// make channel
			inTreeCh := make(chan *model.Addon)
			channels = append(channels, inTreeCh)
			inTreeMonitorCh := make(chan *model.Condition)
			monitorChannels = append(monitorChannels, inTreeMonitorCh)

			go runInTreePlugins(name.Name, inTreeCh)
			go startMonitor(name.Name, inTreeMonitorCh)
		}
	}
	if outTreeLen > 0 {
		outTreePlugins := cfg.Plugins.OutTree
		for _, name := range outTreePlugins {
			outTreeCh := make(chan *model.Addon)
			channels = append(channels, outTreeCh)
			outTreeMonitorCh := make(chan *model.Condition)
			monitorChannels = append(monitorChannels, outTreeMonitorCh)

			go runOutTreePlugins(name.Name, outTreeCh)
			go startMonitor(name.Name, outTreeMonitorCh)
		}
	}
	addonsChannel := model.AddonsChannel{Channels: channels, MonitorChannels: monitorChannels}
	addonsInfo := getAddonsInfo(addonsChannel)

	result := &model.RegisterRequest{Addons: addonsInfo}
	return result, &addonsChannel
}

func runInTreePlugins(name string, ch chan *model.Addon) {

	res := model.Addon{}
	// TODO RUN PLUGIN
	for {
		// return information
		ch <- &res
	}

}

func runOutTreePlugins(name string, ch chan *model.Addon) {
	res := model.Addon{}
	// TODO RUN PLUGIN
	for {
		// return information
		ch <- &res

	}
}

func startMonitor(name string, ch chan *model.Condition) {
	// TODO

}

func getAddonsInfo(channels model.AddonsChannel) []model.Addon {
	var addons []model.Addon

	for _, channel := range channels.Channels {
		addon := <-channel
		addons = append(addons, *addon)
	}

	return addons

}

func getAddonsCondition(channels model.AddonsChannel) []model.Condition {
	var conditions []model.Condition

	for _, channel := range channels.MonitorChannels {
		addon := <-channel
		conditions = append(conditions, *addon)
	}
	return conditions

}

func Heartbeat(channel *model.AddonsChannel, stream config.Channel_EstablishClient, cfg *agentconfig.Configuration) error {
	var heartbeatWithChange model.HeartbeatWithChangeRequest
	lastHeartbeatTime := time.Now()
	lastHeartbeat := &model.HeartbeatWithChangeRequest{}
	firstTime := true
	for {
		sendFlag := false
		var addonsInfo []model.Addon
		var addonsCondition []model.Condition
		// if plugins are specified
		if len(channel.Channels) > 0 {
			// get info
			addonsInfo = getAddonsInfo(*channel)
			addonsCondition = getAddonsCondition(*channel)
			// if not the first time,compare
			if !firstTime {
				for i, addon := range addonsInfo {
					if !reflect.DeepEqual(lastHeartbeat.Addons[i].Properties, addon.Properties) {
						sendFlag = true
					}
				}
				if sendFlag == false {
					// get conditions
					for i, condition := range addonsCondition {
						if condition != lastHeartbeat.Conditions[i] {
							sendFlag = true
						}
					}
				}
			} else {
				firstTime = false
				sendFlag = true
			}
		}
		// TODO CHECK HEALTH
		// send
		if (sendFlag) || ((!sendFlag) && ((time.Now().Sub(lastHeartbeatTime)) > forceSynchronization*time.Second)) {
			if len(channel.Channels) > 0 {
				heartbeatWithChange = model.HeartbeatWithChangeRequest{Healthy: true, Addons: addonsInfo, Conditions: addonsCondition}
			} else {
				heartbeatWithChange = model.HeartbeatWithChangeRequest{Healthy: true}
			}
			lastHeartbeat = &heartbeatWithChange
			request, err := common.GenerateRequest("HeartbeatWithChange", heartbeatWithChange, cfg.ClusterName)
			if err != nil {
				logrus.Error(err)
				continue
			}
			if err := stream.Send(request); err != nil {
				logrus.Error(err)
				continue
			}
			lastHeartbeatTime = time.Now()
			// TODO After Receive Response
			time.Sleep(cfg.HeartbeatPeriod)
		} else {
			request, err := common.GenerateRequest("Heartbeat", HeartbeatMessage, cfg.ClusterName)
			if err != nil {
				logrus.Error(err)
				continue
			}
			if err := stream.Send(request); err != nil {
				logrus.Error(err)
				continue
			}
			lastHeartbeatTime = time.Now()
			time.Sleep(cfg.HeartbeatPeriod)
		}
	}
}

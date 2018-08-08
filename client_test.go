package client_test

import (
	"github.com/go-chassis/go-sc-client"
	"github.com/go-chassis/go-sc-client/model"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/go-chassis/go-chassis/core/lager"
	"os"
)

func TestLoadbalance(t *testing.T) {
	t.Log("Testing Round robin function")
	var sArr []string

	sArr = append(sArr, "s1")
	sArr = append(sArr, "s2")

	next := client.RoundRobin(sArr)
	_, err := next()
	assert.NoError(t, err)
}

func TestLoadbalanceEmpty(t *testing.T) {
	t.Log("Testing Round robin with empty endpoint arrays")
	var sArrEmpty []string

	next := client.RoundRobin(sArrEmpty)
	_, err := next()
	assert.Error(t, err)

}

func TestClientInitializeHttpErr(t *testing.T) {
	t.Log("Testing for HTTPDo function with errors")

	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)

	hostname, err := os.Hostname()
	if err != nil {
		lager.Logger.Error("Get hostname failed.", err)
		return
	}
	microServiceInstance := &model.MicroServiceInstance{
		Endpoints: []string{"rest://127.0.0.1:3000"},
		HostName:  hostname,
		Status:    model.MSInstanceUP,
	}

	registryClient := &client.RegistryClient{}

	err = registryClient.Initialize(
		client.Options{
			Addrs: []string{"127.0.0.1:30100"},
		})
	assert.NoError(t, err)

	err = registryClient.SyncEndpoints()
	assert.NoError(t, err)

	httpHeader := registryClient.GetDefaultHeaders()
	assert.NotEmpty(t, httpHeader)

	resp, err := registryClient.HTTPDo("GET", "fakeRawUrl", httpHeader, []byte("fakeBody"))
	assert.Empty(t, resp)
	assert.Error(t, err)

	MSList, err := registryClient.GetAllMicroServices()
	assert.NotEmpty(t, MSList)
	assert.NoError(t, err)

	f1 := func(*model.MicroServiceInstanceChangedEvent) {}
	err = registryClient.WatchMicroService(MSList[0].ServiceID, f1)
	assert.NoError(t, err)

	var ms = new(model.MicroService)
	var msdepreq = new(model.MircroServiceDependencyRequest)
	var msdepArr []*model.MicroServiceDependency
	var msdep1 = new(model.MicroServiceDependency)
	var msdep2 = new(model.MicroServiceDependency)
	var dep = new(model.DependencyMicroService)
	var m = make(map[string]string)

	m["abc"] = "abc"
	m["def"] = "def"

	dep.AppID = "appid"

	msdep1.Consumer = dep
	msdep2.Consumer = dep

	msdepArr = append(msdepArr, msdep1)
	msdepArr = append(msdepArr, msdep2)

	ms.AppID = MSList[0].AppID
	ms.ServiceName = MSList[0].ServiceName
	ms.Version = MSList[0].Version
	ms.Environment = MSList[0].Environment
	ms.Properties = m

	msdepreq.Dependencies = msdepArr
	s1, err := registryClient.RegisterMicroServiceInstance(microServiceInstance)
	assert.Empty(t, s1)
	assert.Error(t, err)

	s1, err = registryClient.RegisterMicroServiceInstance(nil)
	assert.Empty(t, s1)
	assert.Error(t, err)

	msArr, err := registryClient.GetMicroServiceInstances("fakeConsumerID", "fakeProviderID")
	assert.Empty(t, msArr)
	assert.Error(t, err)

	msArr, err = registryClient.Health()
	assert.NotEmpty(t, msArr)
	assert.NoError(t, err)

	b, err := registryClient.UpdateMicroServiceProperties(MSList[0].ServiceID, ms)
	assert.Equal(t, true, b)
	assert.NoError(t, err)

	f1 = func(*model.MicroServiceInstanceChangedEvent) {}
	err = registryClient.WatchMicroService(MSList[0].ServiceID, f1)
	assert.NoError(t, err)

	f1 = func(*model.MicroServiceInstanceChangedEvent) {}
	err = registryClient.WatchMicroService("", f1)
	assert.Error(t, err)

	f1 = func(*model.MicroServiceInstanceChangedEvent) {}
	err = registryClient.WatchMicroService(MSList[0].ServiceID, nil)
	assert.NoError(t, err)

	str, err := registryClient.RegisterService(ms)
	assert.Empty(t, str)
	assert.Error(t, err)

	str, err = registryClient.RegisterService(nil)
	assert.Empty(t, str)
	assert.Error(t, err)

	ms1, err := registryClient.GetProviders("fakeconsumer")
	assert.Empty(t, ms1)
	assert.Error(t, err)

	err = registryClient.AddDependencies(msdepreq)
	assert.Error(t, err)

	err = registryClient.AddDependencies(nil)
	assert.Error(t, err)

	err = registryClient.AddSchemas(MSList[0].ServiceID, "schema", "schema")
	assert.NoError(t, err)

	getms1, err := registryClient.GetMicroService(MSList[0].ServiceID)
	assert.NotEmpty(t, getms1)
	assert.NoError(t, err)

	getms2, err := registryClient.FindMicroServiceInstances("abcd", MSList[0].AppID, MSList[0].ServiceName, MSList[0].Version)
	assert.Empty(t, getms2)
	assert.Error(t, err)

	getmsstr, err := registryClient.GetMicroServiceID(MSList[0].AppID, MSList[0].ServiceName, MSList[0].Version, MSList[0].Environment)
	assert.NotEmpty(t, getmsstr)
	assert.NoError(t, err)

	getmsstr, err = registryClient.GetMicroServiceID(MSList[0].AppID, "Server112", MSList[0].Version, "")
	assert.Empty(t, getmsstr)
	//assert.Error(t, err)

	ms.Properties = nil
	b, err = registryClient.UpdateMicroServiceProperties(MSList[0].ServiceID, ms)
	assert.Equal(t, false, b)
	assert.Error(t, err)

	err = registryClient.AddSchemas("", "schema", "schema")
	assert.Error(t, err)

	b, err = registryClient.Heartbeat(MSList[0].ServiceID, "")
	assert.Equal(t, false, b)
	assert.Error(t, err)

	b, err = registryClient.UpdateMicroServiceInstanceStatus(MSList[0].ServiceID, "", MSList[0].Status)
	assert.Equal(t, false, b)
	assert.Error(t, err)

	b, err = registryClient.UnregisterMicroService("")
	assert.Equal(t, false, b)
	assert.Error(t, err)
	services, err := registryClient.GetAllResources("instances")
	assert.NotZero(t, len(services))
	assert.NoError(t, err)
	err = registryClient.Close()
	assert.NoError(t, err)

}
func TestRegistryClient_FindMicroServiceInstances(t *testing.T) {
	lager.Initialize("", "DEBUG", "",
		"size", true, 1, 10, 7)

	hostname, err := os.Hostname()
	if err != nil {
		lager.Logger.Error("Get hostname failed.", err)
		return
	}
	ms := &model.MicroService{
		ServiceName: "Server",
		AppID:       "default",
		Version:     "0.0.1",
	}
	var sid string
	registryClient := &client.RegistryClient{}

	err = registryClient.Initialize(
		client.Options{
			Addrs: []string{"127.0.0.1:30100"},
		})
	assert.NoError(t, err)
	sid, err = registryClient.RegisterService(ms)
	if err == client.ErrMicroServiceExists {
		sid, err = registryClient.GetMicroServiceID("default", "Server", "0.0.1", "")
		assert.NoError(t, err)
		assert.NotNil(t, sid)
	}

	microServiceInstance := &model.MicroServiceInstance{
		ServiceID: sid,
		Endpoints: []string{"rest://127.0.0.1:3000"},
		HostName:  hostname,
		Status:    model.MSInstanceUP,
	}

	iid, err := registryClient.RegisterMicroServiceInstance(microServiceInstance)
	assert.NotNil(t, iid)
	_, err = registryClient.FindMicroServiceInstances(sid, "default", "Server", "0.0.1")
	assert.NoError(t, err)
	_, err = registryClient.FindMicroServiceInstances(sid, "default", "Server", "0.0.1")
	//todo revision function is conflicted with go-chassis cache module, disable it until rewrite go-chassis cache module
	//assert.Equal(t, client.ErrNotModified, err)
	assert.Equal(t, nil, err)
	t.Log(err)
	microServiceInstance2 := &model.MicroServiceInstance{
		ServiceID: sid,
		Endpoints: []string{"rest://127.0.0.1:3001"},
		HostName:  hostname + "1",
		Status:    model.MSInstanceUP,
	}
	iid, err = registryClient.RegisterMicroServiceInstance(microServiceInstance2)
	_, err = registryClient.FindMicroServiceInstances(sid, "default", "Server", "0.0.1")
	assert.NoError(t, err)

}

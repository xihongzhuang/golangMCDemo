package api_service

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/stretchr/testify/require"
)

func TestAPIServiceInstance_StartService(t *testing.T) {
	port := 3000
	ep := fmt.Sprintf("http://localhost:%d/api/appmetadata", port)
	instance := NewAPIServiceInstance(port)
	instance.StartService()

	yfile, err := ioutil.ReadFile("../data/valid_metadata1.yaml")
	require.NoError(t, err)
	request, err := http.NewRequest(http.MethodPost, ep, bytes.NewReader(yfile))
	require.NoError(t, err)
	request.Header.Add("Content-Type", "application/x-yml")
	client := &http.Client{Timeout: time.Second * 5}
	response, err := client.Do(request)
	require.NoError(t, err)
	respdata, err := ioutil.ReadAll(response.Body)
	t.Log("request response:", string(respdata))
	response.Body.Close()

	yfile, err = ioutil.ReadFile("../data/valid_metadata2.yaml")
	require.NoError(t, err)
	request, err = http.NewRequest(http.MethodPost, ep, bytes.NewReader(yfile))
	require.NoError(t, err)
	request.Header.Add("Content-Type", "application/x-yml")
	response, err = client.Do(request)
	require.NoError(t, err)
	respdata, err = ioutil.ReadAll(response.Body)
	t.Log("request response:", string(respdata))
	response.Body.Close()

	request, err = http.NewRequest(http.MethodGet, ep, nil)
	require.NoError(t, err)
	response, err = client.Do(request)
	require.NoError(t, err)
	respdata, err = ioutil.ReadAll(response.Body)
	t.Log("request response:", string(respdata))
	response.Body.Close()

	var arr []AppMetaData
	yaml.Unmarshal(respdata, &arr)
	require.Len(t, arr, 2)
	for i, x := range arr {
		t.Log(i, ":", x)
	}
	epFaked := fmt.Sprintf("%s/fakeId1", ep)
	request, err = http.NewRequest(http.MethodDelete, epFaked, nil)
	require.NoError(t, err)
	response, err = client.Do(request)
	require.NoError(t, err)
	respdata, err = ioutil.ReadAll(response.Body)
	t.Log("deleteFakedId request response:", string(respdata))
	response.Body.Close()

	ep1 := fmt.Sprintf("%s/%s", ep, arr[0].Id)
	request, err = http.NewRequest(http.MethodDelete, ep1, nil)
	require.NoError(t, err)
	response, err = client.Do(request)
	require.NoError(t, err)
	respdata, err = ioutil.ReadAll(response.Body)
	t.Log("request response:", string(respdata))
	response.Body.Close()

	//Delete 1 record, there should be only 1 left
	request, err = http.NewRequest(http.MethodGet, ep, nil)
	require.NoError(t, err)
	response, err = client.Do(request)
	require.NoError(t, err)
	respdata, err = ioutil.ReadAll(response.Body)
	t.Log("after delete 1, request response:", string(respdata))
	response.Body.Close()
	var arr2 []AppMetaData
	yaml.Unmarshal(respdata, &arr2)
	require.Len(t, arr2, 1)
	for i, x := range arr2 {
		t.Log(i, ":", x)
	}

	instance.Shutdown(context.Background())
}

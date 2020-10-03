package operator

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"

	"github.com/appscode/jsonpatch"
	"github.com/elek/flekszible/api/v2/data"
	"github.com/elek/flekszible/api/v2/processor"
	"github.com/elek/flekszible/api/v2/yaml"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type AdmissionReviewResult struct {
	Kind       string   `json:"kind"`
	ApiVersion string   `json:"apiVersion"`
	Response   Response `json:"response"`
}

type Response struct {
	Uid       string `json:"uid"`
	Allowed   bool   `json:"allowed"`
	PatchType string `json:"patchType,omitempty"`
	Patch     string `json:"patch,omitempty"`
}
type AdmissionReview struct {
	Kind       string  `json:"kind"`
	ApiVersion string  `json:"apiVersion"`
	Request    Request `json:"request"`
}

type Request struct {
	Uid    string                 `json:"uid"`
	Object map[string]interface{} `json:"object"`
}

//start Kubernetes mutation webhook endpoint
func StartServer(workingDir string) error {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
		})
	})
	r.POST("/", func(c *gin.Context) {
		content, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		result, err := handleRequest(workingDir, content)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(200, result)
	})
	return r.RunTLS("0.0.0.0:8443", "server.crt", "server.key")
}

func handleRequest(dir string, request []byte) (AdmissionReviewResult, error) {
	result := AdmissionReviewResult{
		ApiVersion: "admission.k8s.io/v1beta1",
		Kind:       "AdmissionReview",
	}
	reviewRequest := AdmissionReview{}
	err := json.Unmarshal(request, &reviewRequest)
	if err != nil {
		return result, err
	}
	patch, err := processResource(dir, reviewRequest.Request.Object)
	if err != nil {
		return result, err
	}
	result.Response = Response{
		Uid:     reviewRequest.Request.Uid,
		Allowed: true,
	}
	if len(patch) > 0 {
		logrus.Info("Applying patch:" + patch)
		encodedString := base64.StdEncoding.EncodeToString([]byte(patch))
		result.Response.PatchType = "JSONPatch"
		result.Response.Patch = encodedString
	}
	return result, nil
}

func processResource(dir string, resource map[string]interface{}) (string, error) {
	if resource["metadata"] != nil {
		delete(resource["metadata"].(map[string]interface{}), "managedFields")
		if resource["metadata"].(map[string]interface{})["annotations"] != nil {
			delete(resource["metadata"].(map[string]interface{})["annotations"].(map[string]interface{}), "kubectl.kubernetes.io/last-applied-configuration")
		}
	}
	jsonOrigin, err := json.Marshal(resource)
	if err != nil {
		return "", errors.Wrap(err, "Can't write parsed json resource to json")
	}
	//convert resource to yaml
	yamlResource, err := yaml.Marshal(resource)
	if err != nil {
		return "", errors.Wrap(err, "Can't write parsed json resource to YAML")
	}
	flekszibleResources, err := data.LoadResourceFromByte(yamlResource)
	if err != nil {
		return "", errors.Wrap(err, "Can't read YAML as flekszible Resource")
	}
	//create context
	ctx := processor.CreateRenderContext("k8s", dir, "/tmp/out")
	err = ctx.Init()
	ctx.RootResource.Resources = append(ctx.RootResource.Resources, flekszibleResources...)
	if err != nil {
		return "", err
	}
	//apply transformation
	err = ctx.Render()
	if err != nil {
		return "", err
	}

	//convert yaml to json
	result, err := json.Marshal(ctx.RootResource.Resources[0].Content)
	if err != nil {
		return "", errors.Wrap(err, "Rendered resource file is not a valid json")
	}

	patch, err := jsonpatch.CreatePatch(jsonOrigin, result)
	if err != nil {
		return "", errors.Wrap(err, "Couldn't create json patch")
	}
	patchString, err := json.Marshal(patch)
	if err != nil {
		return "", errors.Wrap(err, "Couldn't marshall json patch")
	}
	return string(patchString), nil
}

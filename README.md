# fle[ksz]ible

Flekszible is a Kubernetes resource manager.

Features:

 1. Based on k8s resource fragments (yaml) and mixin rules
 2. Can generate final k8s resources
 3. Or can generate helm charts.

## Recipes

### Getting started

Put one kubernetes resource to the directory:

./nginx.yaml:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  selector:
    matchLabels:
      app: nginx
  replicas: 2
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.7.9
        ports:
        - containerPort: 80
```

Create an override file to the `./transformations/mylabels.yaml`
```yaml
- type: Add
  path: 
    - metadata
    - annotations
  value: 
    felkszible: generated
```

And execute the generation:

```bash
./flekszible k8s . ./out
```

The result will be something like this:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    felkszible: generated
  annotations: {}
spec:
  selector:
    matchLabels:
      app: nginx
  replicas: 2
  template:
    metadata:
      labels:
        app: nginx
      annotations: {}
    spec:
      containers:
        - name: nginx
          image: nginx:1.7.9
          ports:
            - containerPort: 80
          volumeMounts: []
          env: []
          envFrom: []
      volumes: []
```

As you can see the original k8s source is modified base on the transformation rules.


### Source other dirs

Let's imagine that you would like to run the same nginx as in the previous section but you need 10 replicas for production and 2 for dev.

You can do it with creating 3 directories:

You need the following files:

 * common
   * nginx.yaml (same as before but with replicas = 10) 
 * dev
   * flekszible.yaml (include common)
   * transformations
     * replicas.yaml (override replicas with 2)
 * prod
   * flekszible.yaml (include common)

You can include all the resource files and transformations from common with using the following `flekszible.yaml` in both the `dev` and `prod` folder:

```yaml
import:
  - ../../common
```

And you need a the transformation for `dev/transformations/replicas.yaml`

```yaml
- type: Change
  path:
    - spec
    - replicas
  pattern: .*
  replacement: 2
```


### Set image or namespace

Setting the namespace or an image are very typical task. Therefore they could be activated without creating separated transformations. You can use `--namespace` or `--image` cli arguments which are equivalent with the following transformation files:

`transformations/image.yaml`:

```yaml
- type: ImageSet
  image: elek/flokkr:devbuild
```

`transformations/ns.yaml`:

```yaml
- type: Namespace
  namespace: ozone

```

### Deplloy dev build (the skaffold use case)

Skaffold is a tool which could be used to deploy a specific dev build to the kubernetes cluster. While skaffold has many functionality (automatic redeploy, coud build) the basic functionality (local build, simple deploy) could be replaced with the following 4 lines:

```bash
export IMAGE=elek/ozone:$(git describe --tag)
docker build -t $IMAGE .
docker push $IMAGE
flekszible k8s --image=$IMAGE --namespace=mynamespace k8s/resources/ - | kubectl apply -f 
```

Notes:

 * `flekszible.yaml` configuration file is optional
 * You can generate the k8s resources files to the standard output instead of directory (all the additional log lines are suppressed)
 * image and namespace could be changed without any config file

## Definitions

### Path

Many transformations requires a `Path` definition to address a specific node in the yaml tree.

Path is a string array where each element represents a new level. Array elements are indexed from zero (eg ["foo","barr",0]) __except_ if the elements in the array are maps with _name_ key. In this case the index is this name.

For example:

```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    felkszible: generated
spec:
  selector:
    matchLabels:
      app: nginx
  replicas: 2
  template:
    metadata:
      labels:
        app: nginx
      annotations: {}
    spec:
      containers:
        - name: nginx
          image: nginx:1.7.9
          ports:
            - containerPort: 80
          env: 
            - name:  KEY
              value: XXX
```

Here tha path of the environment variable is `[spec, template, spec, containers, nginx, env, KEY ]` and not `[ spec, template, spec, containers, 0, env, 0]`

For matching, path segments are used as regular expressions. Therefore the following path matches for both the init and main containers:

```yaml
path:
  - spec
  - template
  - spec
  - (initC|c)ontainers
  - .*
```

### Trigger

Most of the processors also use a `trigger` parameter. With trigger you can specify any k8s fragments and only the matched resources will be transformed.

For example:

```yaml
- type: Add
  trigger:
     metadata:
        name: datanode
  path:
    - metadata
    - labels
  value:
     flokkr.github.io/monitoring: false
```

This definition will apply only to the k8s resources where the value of `metadata.name` is `datanode`.

You can use multiple values in the trigger. All the key nodes will be collected and should be the same in the target resource.


## Directory structure

TLDR;

 * `*.yaml`: used as k8s resources
 * `transformations/*.yaml`: will be applied to all the resources according to the specified rules
 * `definitions/*.yaml`: composit definitions which could be used in `transformations.yaml`. Won't be applied by default.
 * `flekszible.yaml`: configuration file. Could include other directories (all the resources + transformations + definitions will be added from that directory.)
 
In the output directory all the `yaml` files (except the `flekszible.yaml` configuration file) are considered to be a k8s resource. One file could contain multiple resources.

All the yaml files from the `transformations` subdirectory are parsed as transformation definitions. Each file should contain an array of objects. The object should have a `type` definition (see the available transformations below) 

Example: `./transformations/label.yaml`

```
- type: Add
  path: 
    - metadata
    - annotations
  value: 
    felkszible: generated
```

All the yaml files from the `definitions` directory will be parsed as composit transformation type. You can define multiple transoformation and name it. It may be used form other transformation files.

## Available transformation types

### Add

Extends existing k8s resources with additional elements.

#### Parameters

| Name    | Type     | Value 
|---------|----------|-------
| path    | []string | Path in the resource file to 
| value   | yaml     | Any yaml fragment which should be added to the node which is selected by the `Path`

#### Example

```
- type: Add
  path: 
    - metadata
    - annotations
  value: 
    felkszible: generated
```

#### Supported add methods

| Type of the destination node (selected by `Path`) | Type of the `Value` | Supported
|---------------------------------------------------|---------------------|------------
| Map                                               | Map                 | Yes
| Array                                             | Array               | Yes


#### ImageSet

Replaces the docker image definition eveywhere

| Name    | Type     | Value 
|---------|----------|-------
| image   | string   | Full name of the required docker image. Such as 'elek/ozone:trunk'


Note: This transformations could also added with the `--image` CLI argument.

#### PublishStatefulSet

Creates additional NodeType service for StatefulSet internal services.


#### DaemonToStatefulset

Converts daemonset to statefulset.

Useful for minikube based environments where you may not have enough node to run a daemonset based cluster.

#### Composit

You can create additional transformations with grouping existing transformations. For example the following definition register a new transformation type:


```yaml
type: flokkr.github.io/prometheus
transformations:
  - type: Add
    path:
      - spec
      - template
      - metadata
      - annotations
    value:
      prometheus.io/scrape: "true"
      prometheus.io/port: "28942"
  - type: Add
    path:
      - spec
      - template
      - spec
      - containers
      - "*"
      - env
    value:
      - name: "PROMETHEUSJMX_ENABLED"
        value: "true"
      - name:  "PROMETHEUSJMX_AGENTOPTS"
        value: "port=28942"

```

Put the previous file to the `definitions/prometheus.yaml`. By default it won't be applied to any k8s resources but from now you can use the `flokkr.github.io/prometheus` type in your transformations:

For example in `transformations/monitor.yaml` you can write:

```yaml
- type: flokkr.github.io/prometheus
```
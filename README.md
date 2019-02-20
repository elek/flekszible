# fle[ksz]ible

Flekszible is a Kubernetes resource manager. It's composition based (like kustomize) instead of templates (like helm).

Compared to kustomize or ksonnet:

 * It's almost as powerful as them but it's more simple to use. The key challenge here is to find the balance between simplicity and usability.

Features:

 1. Zero-config: it can work without any external files
 2. Mixins: you can define additional transformations to change k8s resources
 3. Imports: You can compose resources from multiple sources
 4. Multi-tenancy: With imports you can manage multiple environments (dev,prod,...)
 5. Multi-instance: You can import the same template (eg. zookeeper resources) with different flavour. With this approach you can create two different zookeeper ring from a template to your cluster.
 
## Recipes (Features)

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
./flekszible generate . ./out
```

__Note__: Before v0.4.0 you should use `flekszible k8s` instead of `flekszible generate`

The result will be something like this:

```yaml
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
    spec:
      containers:
        - name: nginx
          image: nginx:1.7.9
          ports:
            - containerPort: 80

```

As you can see the original k8s source is modified base on the transformation rules.


### Import/Source other dirs

Let's imagine that you would like to run the same nginx as in the previous section but you need 10 replicas for production and 2 for dev.

You can do it with creating 3 directories:

You need the following files:

 * `common`
   * `nginx.yaml` (same as before but with replicas = 10) 
 * `dev`
   * `flekszible.yaml` (include common)
   * `transformations`
     * `replicas.yaml` (override replicas with 2)
 * `prod`
   * `flekszible.yaml` (include common)

You can include all the resource files and transformations from common with using the following `flekszible.yaml` in both the `dev` and `prod` folder:

```yaml
import:
  - path: ../../common
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
- type: Image
  image: elek/flokkr:devbuild
```

`transformations/ns.yaml`:

```yaml
- type: Namespace
  namespace: ozone

```

### Deploy dev build (the skaffold use case)

Skaffold is a tool which could be used to deploy a specific dev build to the kubernetes cluster. While skaffold has many functionality (automatic redeploy, coud build) the basic functionality (local build, simple deploy) could be replaced with the following 4 lines:

```bash
export IMAGE=elek/ozone:$(git describe --tag)
docker build -t $IMAGE .
docker push $IMAGE
flekszible generate --image=$IMAGE --namespace=mynamespace k8s/resources/ - | kubectl apply -f 
```

Notes:
 
 * Before v0.4.0 you should use `flekszible k8s` instead of `flekszible generate`
 * `flekszible.yaml` configuration file is optional
 * You can generate the k8s resources files to the standard output instead of directory (all the additional log lines are suppressed)
 * image and namespace could be changed without any config file

### Instantiate 

During the import of an external resource set you can apply additinonal transformations just for the imported resources. 

Example `flekszible.yaml`:

```yaml
import:
  - path: ./zookeeper
    transformations:
      - type: Prefix
        prefix: zk1
```

Here the resources from the zookeeper dir will be imported to the current kubernetes resource set with an additional prefix.

With this method you can import the same resource more than once. For example if you need a separated zookeeper instance for Hadoop HDFS HA and an for HBase you can import it twice:

```yaml
import:
  - path: ./zookeeper
    transformations:
      - type: Prefix
        prefix: zk1
  - path: ./zookeeper
    transformations:
      - type: Prefix
        prefix: zk2
```

As a result the zookeeper instances will be imported twice with different prefixes:

```bash
ls -lah 
Permissions Size User Date Modified Name
.rw-r--r--   184 elek 29 Dec 12:15  zk1-zookeeper-service.yaml
.rw-r--r--   749 elek 29 Dec 12:15  zk1-zookeeper-statefulset.yaml
.rw-r--r--   184 elek 29 Dec 12:15  zk2-zookeeper-service.yaml
.rw-r--r--   749 elek 29 Dec 12:15  zk2-zookeeper-statefulset.yaml
```

## Reference

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

## Import

You can import other directory structures with adding references to the `flekszible.yaml`

For example

```
import:
  - path: ../../hadoop
```

All the transformations + definitions + k8s resources will be added and applied. Note: the transformations from the imported directory will be applied only to the imported resources.

The imported resources could be generated to a subdirectory:

```
import;
  - path: ../../hadoop
  - path: ../../prometheus
    destination: monitoring
```

With this approach the prometheus related resources will be saved to the `monitoring` subdirectory of the destination path.

Transforamtions also can be applied to the imported resources:

```
import:
  - path: ../../hadoop
    transformations:
       - type: Image
         image: elek/ozone
```

The hadoop resources are imported here and the image reference is changed during the import (only for the imported resources).

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
| Array                                             | Map                 | Yes


#### Image

Replaces the docker image definition everywhere

| Name    | Type     | Value 
|---------|----------|-------
| image   | string   | Full name of the required docker image. Such as 'elek/ozone:trunk'


Note: This transformations could also added with the `--image` CLI argument.

#### Namespace

Similar to the image namespace also can be changed with simple transformation:


| Name    | Type     | Value 
|---------|----------|-------
| namespace   | string   | Name of the used kubernetes namespace.

Note: This transformations could also added with the `--namespace` CLI argument.

Example (`transformations/set.yaml`):

```yaml
- type: Namespace
  namespace: myns
```

#### Prefix

Add a specific prefix for all of the names.


| Name     | Type     | Value 
|----------|----------|-------
| prefix   | string   | Prefix which will be added to all the names.

Example (`transformations/set.yaml`):

```yaml
- type: Namespace
  namespace: myns
```

#### ConfigHash

Add a kubernetes annotation with the hash of the used configmap. With 
this approach you can force to re-create the k8s resources in case of config change. 
In case of configmap change the annotation value will be different and the resource
will be recreated.

As of now it supports only one configmap per resource and only the top-level
resource will be annotated (in case of statefulset this is the statefulset 
not the pod).

Example (`transformations/config.yaml`):

```yaml
- type: ConfigHash
```

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
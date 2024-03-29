---
title: Recipes
menu: main
---

# Getting started

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
./flekszible generate -d ./out
```

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

In the transformation file we defined:

 1. The type of the transformation (`Add`). You can see the available transformations below or you can define your own composite transformations
 2. The path to modify (this is defined by the `Add` transformation)
 3. The value which should be added (the is also `Add` specific)

# Import/Source other dirs

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

# Set image or namespace

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

# Deploy dev build (the skaffold use case)

Skaffold is a tool which could be used to deploy a specific dev build to the kubernetes cluster. While skaffold has many functionality (automatic redeploy, coud build) the basic functionality (local build, simple deploy) could be replaced with the following 4 lines:

```bash
export IMAGE=elek/ozone:$(git describe --tag)
docker build -t $IMAGE .
docker push $IMAGE
flekszible generate --image=$IMAGE --namespace=mynamespace -s k8s/resources/ -d - | kubectl apply -f 
```

Notes:

 * `flekszible.yaml` configuration file is optional.
 * You can generate the k8s resources files to the standard output instead of directory (all the additional log lines are suppressed)
 * image and namespace could be changed without any config file 
 
# Instantiate 

During the import of an external resource set you can apply additinonal transformations just for the imported resources. 

Example `flekszible.yaml`:

```yaml
import:
  - path: zookeeper
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

# Destination dir support

Not all of the resources are equal. Sometimes it's better to use a hierarchy for the generated resources:

```
import:
    - path: ozone
      transformations:
        - type: ozone/prometheus
    - path: ozone-csi
      destination: csi
    - path: prometheus
      destination: monitoring

```

Here we imported 3 subcomponent from the subdirectories. But the resources from the `prometheus` and `csi` subdirectories will be generated to the `monitoring` and `csi` subdirectory of the destination directory. 

# External sources

Resources can be imported from external sources with the help of the [go-getter](https://godoc.org/github.com/hashicorp/go-getter) library.

```
source:
    - url: github.com/flokkr/k8s
import:
    - path: ozone
    - path: prometheus
      destination: monitoring
```

This `flekszible.yaml` works out of the box: It downloads the `k8s` repository to a cache folder and imports `ozone` and `prometheus` subfolders from there. 

# FLEKSZIBLE_PATH

The source directory also can be defined with the `FLEKSZIBLE_PATH` environment variable.

The previous example can work if you have the `k8s` project in your home directory.
```
export FLEKSZIBLE_PATH=~/k8s
```

```
import:
    - path: ozone
    - path: prometheus
      destination: monitoring
```

# Define transformations

The default directory structure of a component is:

 * `flekszible.yaml` (optional)
 * `transformations` directory which contains 
 * `definitions` reusable definitions
 * `configmaps` raw resources to import them as configmaps (name convention: configmapname_keyname.ext)
 * `...*.yaml` Any kubernetes sresource file.
 
 In the transformation we can create any number files which can contain transformations:
 
 For example the `transformations/fixhosts.yaml`
 
 ```
- type: Add
  trigger:
      metadata:
          name: scm
  path:
    - spec
    - template
    - spec
  value:
    nodeSelector:
       kubernetes.io/hostname: node1.flekszible.com
- type: Add
  trigger:
      metadata:
          name: om
  path:
    - spec
    - template
    - spec
  value:
    nodeSelector:
       kubernetes.io/hostname: node1.flekszible.com
```

This file contains two transformation rule. The first one modifies the `scm` statefulset (see the `trigger` condition) and it adds (type is `Add`) a custom nodeSelector (`value`) to the spec/template/spec part of the kubernetes yaml (defined by the `path`)

This transformation will be executed an all the resource files, but:
 
  * Only if the trigger condition is matched
  * Only if the patch exists


# List of definitions

You can show the available processor definitions (transformations templates) which can be used in the `type` field with executing:

```
> flekszible processor

+---------------------+-----------------------------------------------------------------------+
| name                | description                                                           |
+---------------------+-----------------------------------------------------------------------+
| Image               | Replaces the docker image definition                                  |
| K8sWriter           | Internal transformation to print out k8s resources as yaml            |
| Namespace           | Use explicit namespace                                                |
| ConfigHash          | Add labels to the k8s resources with the hash of the used configmaps  |
| DaemonToStatefulSet | Converts daemonset to statefulset                                     |
| Prefix              | Add same prefix to all the k8s names                                  |
| PublishStatefulSet  | Creates additional NodeType service for StatefulSet internal services |
| ozone/prometheus    | Enable prometheus monitoring in Ozone                                 |
| Add                 | Extends yaml fragment to an existing k8s resources                    |
| Change              | Replace existing value literal in the yaml struct                     |
+---------------------+-----------------------------------------------------------------------+

```

And you can check the available variable for one definition:

```
> flekszible processor show Add

### Add

Extends yaml fragment to an existing k8s resources

#### Parameters

+-------+---------+--------------------------------------------+
| name  | default | description                                |
+-------+---------+--------------------------------------------+
| value |         | A yaml struct to replace the defined value |
+-------+---------+--------------------------------------------+


#### Supported value types

| Type of the destination node (selected by 'Path') | Type of the 'Value' | Supported
|---------------------------------------------------|---------------------|------------
| Map                                               | Map                 | Yes
| Array                                             | Array               | Yes
| Array                                             | Map                 | Yes

#### Example

'''
- type: Add
  path:
  - metadata
  - annotations
  value:
     flekszible: generated
'''

```

# Define reusable processor definitions

As you saw in the previous section there are a few predefined definition which can be used for transformation type. But you can also define your own one. Put the transformation

# Template support

Put this file to your `definitions/prometheus.yaml`

```
name: ozone/prometheus
description: Enable prometheus monitoring in Ozone
---
- type: Add
  trigger:
      metadata:
          name: config
  path:
    - data
  value:
    OZONE-SITE.XML_hdds.prometheus.endpoint.enabled: true
```

This file won't be applied to any resource file but it's a reusable definition type. If any of the imported directory or your current directory contains such definition, it will be available as a transformation type. For example in your `flekszible.yaml` you can use it during the import:

```
import:
  - path: ozone
    transformations: 
      - type: ozone/prometheus

```

Note: you can use it both in `transformation/...yaml` or in the `flekszible.yaml` under the transformation section. The syntax is the same but the first one will be activated only for the imported resources while the second one will be executed all the resources.

Note: the imported definitions also can be listed together with the build in processors:

```
> flekszible processor

+---------------------+-----------------------------------------------------------------------+
| name                | description                                                           |
+---------------------+-----------------------------------------------------------------------+
| Namespace           | Use explicit namespace                                                |
| PublishStatefulSet  | Creates additional NodeType service for StatefulSet internal services |
| ozone/prometheus    | Enable prometheus monitoring in Ozone                                 |
| Add                 | Extends yaml fragment to an existing k8s resources                    |
| Change              | Replace existing value literal in the yaml struct                     |
| ConfigHash          | Add labels to the k8s resources with the hash of the used configmaps  |
| Prefix              | Add same prefix to all the k8s names                                  |
| DaemonToStatefulSet | Converts daemonset to statefulset                                     |
| Image               | Replaces the docker image definition                                  |
| K8sWriter           | Internal transformation to print out k8s resources as yaml            |
+---------------------+-----------------------------------------------------------------------+

> flekszible processor show prometheus

### ozone/prometheus

Enable prometheus monitoring in Ozone

#### Parameters

No parameters.

```

# Templating in the definition files

The reusable definitions can have additional parameter variables.

See this example (`zookeeper/definitions/scale.yaml`):

```
name: zookeeper/scale
description: Set the number of the zookeeper replicas.
parameter:
  - name: replicas
    type: int
    default: 1
    description: Number of the required replicas
---
- type: Change
  path:
    - spec
    - replicas
  pattern: .*
  replacement: 3
- type: Change
  trigger:
     metadata:
        name: zookeeper-config
  path:
    - data
    - "zoo.cfg"
  pattern: .*server.0.*
  replacement: |-
    {{range $val := Iterate .replicas}}
    server.{{.}}=zookeeper-{{.}}.zookeeper:2888:3888
    {{end}}
```

This definition can be activated to scale up Zookeeper instances. It requires the replacement of the replica number in the statefulset (`Change`) and and additional line in the Configmap. 

All the defintions contain two parts. The first part is the metadata: name + variable definition. The second part is the transformation definition but before the usage it will be rendered as a go template.

The previous example can be used:

```
source:
    - url: github.com/flokkr/k8s
import:
    - path: zookeeper
      transformations:
          - type: zookeeper/scale
            replicas: 3
```

This `flekszible.yaml` imports all the resources (and pre created definitions) from the `k8s` git repository and will use the previous definition (which is defined in the imported repository).


# Package management

Import path is transitive. If you import a directory you will see all the imported kubernetes resources and definitions. To make it easier to follow what can be used you can list the available directories (Directories which contain `flekszible.yaml` with valid metadata).

Let's create a new empty directory:

```
mkdir /tmp/flekszible-demo && cd /tmp/flekszible-demo
```

Add a remote source (this is equivalent to create a `flekszible.yaml` and add a `source` tag:

```
flekszible source add github.com/flokkr/k8s
```

Now your `flekszible.yaml` contains the following:

```
source:
- url: github.com/flokkr/k8s
```

Now you can search for the available directories which can be imported:

```
 flekszible app search
+-------------+--------------------------------------------------------+
| path        | description                                            |
+-------------+--------------------------------------------------------+
| grafana     | Grafana dashboard server                               |
| hdfs        | Apache Hadoop HDFS base setup                          |
| hdfs-ha     | Apache Hadoop HDFS, HA setup                           |
| jaeger      | Jaeger tracing server                                  |
| kafka       | Apache Kafka                                           |
| ozone/csi   | CSI server to use Apache Hadoop Ozone via s3           |
| ozone       | Apache Hadoop Ozone                                    |
| ozone/freon | Load test tool for Apache Hadoop Ozone                 |
| prometheus  | Prometheus monitoring                                  |
| pv-test     | Nginx example deployment with persistent volume claim. |
| zookeeper   | Scalable Apache Zookeeper setup                        |
+-------------+--------------------------------------------------------+

```

Well, it looks good. Let's import `ozone` together with `prometheus`:

```
flekszible app add ozone
flekszible app add prometheus
```

Now your `flekszible.yaml` is modified to: 

```
source:
- url: github.com/flokkr/k8s
import:
- path: ozone
- path: prometheus
```

Let's generate the kubernetes resources:

```
flekszible generate -d output
```

And you have all the kubernetes resource files:

```
.
├── flekszible.yaml
└── output
    ├── config-configmap.yaml
    ├── datanode-daemonset.yaml
    ├── om-service.yaml
    ├── om-statefulset.yaml
    ├── prometheus-clusterrole.yaml
    ├── prometheusconf-configmap.yaml
    ├── prometheus-deployment.yaml
    ├── prometheus-operator-clusterrolebinding.yaml
    ├── prometheus-operator-serviceaccount.yaml
    ├── prometheus-service.yaml
    ├── s3g-service.yaml
    ├── s3g-statefulset.yaml
    ├── scm-service.yaml
    └── scm-statefulset.yaml
```

# Package management: custom repository

You can search for the available repositories with 

```
flekszible search source
```

You don't need to create any PR to show up your own repository. Just tag your repository with `flekszible` topic. The previous command is just a simple github search.

(As there is no moderation, including a repository is just as safe as executing a downloaded script.)

# Service-mesh support

There are multiple built-in processor definition which can be parameterized in the `transformations` file. One notable is the `pipe` type which can execute any command to transform the given kubernetes manifest file.

For example to inject istio service-mesh fragments, you can define the following `flekszible.yaml`:

```
import:
- path: kube
  transformations:
  - type: pipe
    command: istioctl
    args:
    - kube-inject
    - "-f"
    - "-"

```

Note: you should put the `istioctl` cli to your path and define the istio sample dir as a source of imports. For example with environment vaiable:

```
example FLEKSZIBLE_PATH=/home/elek/prog/istio-1.0.5/samples/bookinfo/platform/
```

Please note, that istio injection could be slow as it requires a live connection to the existing istio services.

# Kustomize integration

`Flekszible` definitions can be defined in k8s resource file which is compatible with Kustomize: 

```yaml
apiVersion: flekszible.github.io/v1
kind: Flekszible
metadata:
  name: cluster1
  annotations:
    config.kubernetes.io/function: |
      exec: 
        path: flekszible

spec:
  transformations:
    - type: image
      image: xxx
```

If you use this file as a transformer plugin, flekszible can modify the resources managed by kustomize (as long as you have the latest binary on your path):

```
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - simple.yaml
  - simple2.yaml
metadata:
  name: kustomize-integration
transformers:
  - flekszible.yaml
```

Please note that `kustomize` plugin can receive all the resources on standard output, but you couldn't easily use external directories. Therefor it's better to use only inline transformation in the `Flekszible` file itself.
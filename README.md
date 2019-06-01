# fle[ksz]ible

Flekszible is a Kubernetes configuration/manifest manager. It helps to manage your kubernetes yaml files before the deployment.

 * Flekszible generates the final k8s yaml files based on source yaml file + easy to use transformation rules
 * Everything can be versioned as the final files are generated and saved (GitOps)
 * Transformation rules can be reused and shared (It has a powerful but extremly simple package management)
 * It combines the best part of the composition based and template based approach (it's composition based but in reusable transformation rule definition you can use templates)

## Features:

  1. Zero-config: it can work without any external files
  2. Mixins: you can define additional transformations to change k8s resources
  3. Imports: You can compose resources from multiple sources (github repositories)
  4. Multi-tenancy: With imports you can manage multiple environments (dev,prod,...)
  5. Multi-instance: You can import the same template (eg. zookeeper resources) with different flavour. With this approach you can create two different zookeeper ring from a template to your cluster.
  6. Reusable transformations: you can define transformations and reuse them later (or provide them to the user as optional flags.)
  7. Package management:  Simple github based package management. No separated repository format just github repositories and tags.
  8. Side-car pattern friendly design: additional side-car containers can be injected
  10. __build time transformation_: compared to the oprerator pattern, here everything is visible as the transformations are applied on build time before using `kubectl apply -f`. It's more safe to switch between environments.
  11. GitOps friendy: generates all the final resources to static files
  12. Supports external processors like service-mesh injectors

For more information and for comparison with Helm and Kustomize: [Check the docs](https://flekszible.netlify.com/) 

## Install

On macOS, you can install flekszible with Homebrew package manager:

```brew 
brew install elek/brew/flekszible
```

For linux: download the binary from the [Release page](https://github.com/elek/flekszible/releases)

## Documentation

Latest docs are available from [HERE](https://flekszible.netlify.com/)
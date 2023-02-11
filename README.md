## Installation
```sh
go get github.com/LyridInc/cluster-api-go-sdk
```

## Examples

### Setup OpenStack environment variables
```go
yamlByte, _ := os.ReadFile("./test/data/clouds.yaml")
cloudsYaml := model.CloudsYaml{}
cloudsYaml.Parse(yamlByte)
opt := option.OpenstackGenerateClusterOptions{
  ControlPlaneMachineFlavor: "SS2.2",
  NodeMachineFlavor:         "SM8.4",
  ExternalNetworkId:         "79241ddc-c51b-4677-a763-f48c60870923",
  ImageName:                 "ubuntu-2004-kube-v1.24.8",
  SshKeyName:                "kube-key",
  DnsNameServers:            "8.8.8.8",
  FailureDomain:             "az-01",
  IgnoreVolumeAZ:            true,
}
cloudsYaml.SetEnvironment(opt)
```

### Initialize infrastructure OpenStack

```go
infrastructure := "openstack"
capi := api.NewClusterApiClient("", "<YOUR KUBECONFIG FILE PATH>")

capi.InitInfrastructure(infrastructure)
```

### Authenticate OpenStack client

```go
cl := api.OpenstackClient{
  NetworkEndpoint: os.Getenv("OS_NETWORK_ENDPOINT"),
  AuthEndpoint:    os.Getenv("OS_AUTH_ENDPOINT"),
  AuthToken:       os.Getenv("OS_TOKEN"),
  ProjectId:       os.Getenv("OS_PROJECT_ID"),
}

err = cl.Authenticate(api.OpenstackCredential{
  ApplicationCredentialName:   os.Getenv("OS_APPLICATION_CREDENTIAL_NAME"),
  ApplicationCredentialId:     os.Getenv("OS_APPLICATION_CREDENTIAL_ID"),
  ApplicationCredentialSecret: os.Getenv("OS_APPLICATION_CREDENTIAL_SECRET"),
})
if err != nil {
  log.Fatal("Openstack authentication error:", err)
}
```

### Generate workload cluster
```go
clusterName := "capi-local-2"
clusterOpt := option.GenerateWorkloadClusterOptions{
  ClusterName:              clusterName,
  KubernetesVersion:        "v1.24.8",
  WorkerMachineCount:       3,
  ControlPlaneMachineCount: 1,
  InfrastructureProvider:   infrastructure,
  Flavor:                   "external-cloud-provider",
}
yaml, err = capi.GenerateWorkloadClusterYaml(clusterOpt)
if err != nil {
  log.Fatal("Generate workload cluster error:", err)
}

// update cidr blocks and allow all in cluster traffic to support flannel cni
yaml, _ = cl.UpdateYamlManifest(yaml, option.ManifestOption{
  ClusterKindSpecOption: option.ClusterKindSpecOption{
    CidrBlocks: []string{"10.244.0.0/16"},
  },
  InfrastructureKindSpecOption: option.InfrastructureKindSpecOption{
    AllowAllInClusterTraffic: true,
  },
})

if err := capi.ApplyYaml(yaml); err != nil {
  log.Fatal(err)
}
```

### Get workload cluster kubeconfig
```go
conf, err := capi.GetWorkloadClusterKubeconfig(clusterName)
if err == nil {
  if err := os.WriteFile(fmt.Sprintf("./%s.kubeconfig", clusterName), []byte(*conf), 0644); err != nil {
    log.Fatal("Write kubeconfig error:", err)
  }
}
```

### Create secret
```go
cloudConf := os.Getenv("OPENSTACK_CLOUD_PROVIDER_CONF")
if cloudConf == "" {
  log.Fatal("Error reading cloud conf: OPENSTACK_CLOUD_PROVIDER_CONF is not set")
}

if err := capi.SetKubernetesClientset("<YOUR WORKLOAD CLUSTER KUBECONFIG FILE PATH>"); err != nil {
  log.Fatal("Error set kubeconfig:", error.Error(err))
}
secret := v1.Secret{
  TypeMeta: metav1.TypeMeta{
    APIVersion: v1.SchemeGroupVersion.String(),
    Kind:       "Secret",
  },
  ObjectMeta: metav1.ObjectMeta{
    Name:      "cloud-config",
    Namespace: "kube-system",
  },
  Data: map[string][]byte{
    "cloud.conf": []byte(cloudConf),
  },
}

if _, err := capi.CreateSecret(secret); err != nil {
  fmt.Println("Error create secret:", error.Error(err))
  fmt.Println("Trying again..")
  time.Sleep(3 * time.Second)
} else {
  break
}
```

### Installing Flannel CNI
```go
yaml, err = model.ReadYamlFromUrl(option.FLANNEL_MANIFEST_URL)
if err != nil {
  log.Fatal(error.Error(err))
}

if err := capi.ApplyYaml(yaml); err != nil {
  log.Fatal("Error apply flannel cni yaml:", error.Error(err))
}

secretName := fmt.Sprintf("%s-csi-cloud-secret", clusterName)

for _, url := range option.OPENSTACK_CLOUD_CONTROLLER_MANIFEST_URLS {
  yaml, err := model.ReadYamlFromUrl(url)
  if err != nil {
    log.Fatal(error.Error(err))
  }

  if err := capi.ApplyYaml(yaml); err != nil {
    log.Fatal("Error apply yaml:", url, " - ", error.Error(err))
  }
}
```

### Installing Cinder CSI Driver on OpenStack provider cluster
```go
yaml, err = model.ReadYamlFromUrl(option.OPENSTACK_CINDER_DRIVER_MANIFEST_URLS["secret"].(string))
if err != nil {
  log.Fatal(error.Error(err))
}

cloudConfB64 := os.Getenv("OPENSTACK_CLOUD_PROVIDER_CONF_B64")
if cloudConfB64 == "" {
  log.Fatal("Error reading cloud conf: OPENSTACK_CLOUD_PROVIDER_CONF_B64 is not set")
}

yamlResult, err := cl.UpdateYamlManifest(yaml, option.ManifestOption{
  SecretKindOption: option.SecretKindOption{
    Data: map[string]interface{}{
      "cloud.conf": cloudConfB64,
    },
    Metadata: map[string]interface{}{
      "name":      secretName,
      "namespace": "kube-system",
    },
  },
})
if err != nil {
  log.Fatal("Update yaml from url error:", error.Error(err))
}

if err := capi.ApplyYaml(yamlResult); err != nil {
  log.Fatal("Error create cinder csi secret:", error.Error(err))
}

pluginUrls := option.OPENSTACK_CINDER_DRIVER_MANIFEST_URLS["plugins"].([]string)
for i, url := range pluginUrls {
  yaml, err := model.ReadYamlFromUrl(url)
  if err != nil {
    log.Fatal(error.Error(err))
  }

  if url == "https://raw.githubusercontent.com/kubernetes/cloud-provider-openstack/master/manifests/cinder-csi-plugin/cinder-csi-controllerplugin.yaml" {
    // deployment
    yamlResult, err := cl.UpdateYamlManifest(yaml, option.ManifestOption{
      DeploymentKindOption: option.DeploymentKindOption{
        VolumeSecretName: secretName,
      },
    })
    if err != nil {
      log.Fatal("Update yaml from url error:", error.Error(err))
    }

  } else if url == "https://raw.githubusercontent.com/kubernetes/cloud-provider-openstack/master/manifests/cinder-csi-plugin/cinder-csi-nodeplugin.yaml" {
    // daemonset
    yamlResult, err := cl.UpdateYamlManifest(yaml, option.ManifestOption{
      DaemonSetKindOption: option.DaemonSetKindOption{
        VolumeSecretName: secretName,
      },
    })
    if err != nil {
      log.Fatal("Update yaml from url error:", error.Error(err))
    }
  }

  if err := capi.ApplyYaml(yaml); err != nil {
    log.Fatal("Error apply yaml:", url, " - ", error.Error(err))
  }
}
```

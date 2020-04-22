# LogServer Daemon for Viewing Application Logs

LogServer DaemonSet runs a daemon for fetching application logs from shared volume.
It also provides viewing STDOUT & STDERR logs for a given container.

## URL For Viewing Log
`http://<NODE_IP>:11000/logs/<POD_UID>/<FILE_NAME>`

e.g `http://localhost:11000/logs/d2649681-d971-4d47-93d9-cc2a7cba2c1b/out.log`

## List Files
`http://<NODE_IP>:11000/logs/<POD_UID>`

e.g `http://localhost:11000/logs/d2649681-d971-4d47-93d9-cc2a7cba2c1b`

* **NODE_IP**: IP Address of Node where POD is running
* **POD_UID**: UID of the POD
* **FILE_NAME**: Log file name for viewing

# How to run locally

## 1) Build & Publish LogServer Docker Image

```shell
  $ cd logserver-daemon
  $ docker build . -t mahkath/logserver:latest
  $ docker push mahkath/logserver:latest
```

## 2) Start Local Kubernetes via Kind

```shell
  $ kind create cluster
```

## 3) Deploy DaemonSet to k8s cluster

```shell
  $ kubectl apply -f daemon/logserver-daemonset.yaml
```

## 4) Deploy Sample App to Cluster(emptyDir logs folder)

```shell
  $ kubectl apply -f examples/sample-deployment.yaml
```

## 5) Port Forward (Only in Local Cluster)

```shell
  # Get Pod Name for LogServer DaemonSet
  $ logserverpodname=$(k get pods -o jsonpath='{ .items[?(@.metadata.ownerReferences[0].name=="logserver")].metadata.name}')
  $ kubectl port-forward $logserverpodname 11000:11000
```
## 6) Access Logs via LogServer WebView

```shell
  # Get UID for Running MyApp POD
  $ podname=$(k get pods -o jsonpath='{ .items[?(@.metadata.labels.app=="myapp")].metadata.name}')
  $ uid=$(k get pod $podname -o jsonpath='{.metadata.uid}')

  # List Logs Name in EmptyDir of MyApp
  $ curl http://localhost:11000/logs/$uid
```


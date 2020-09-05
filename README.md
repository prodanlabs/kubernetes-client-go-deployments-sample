# kubernetes-client-go-deployments-sample

使用场景:

Kubernetes集群没有公网IP,网络可以出公网,但无法从外部发起更新版本请求,该工具通过轮询harbor仓库的最新镜像进行更新deploy.

```shell
root@client-go-deploy-68857c6778-d7h9j:/# ./app -v
Version=v1.0.0
GitCommitLog=91d7f33dcc7d8081d526015a9113f45655b114ec init
BuildTime=2020-09-04 13:39:01
GoVersion=go version go1.14.4 linux/amd64
runtime=linux/amd64
root@client-go-deploy-68857c6778-d7h9j:/# ./app -h
Usage of ./app:
  -apps string
        namespace_deploy_container.Examples: kube-system.busybox.busybox,default.for.test. (default "kube-system.busybox.alpine,default.busybox2.alpine")
  -cluster string
        out-of or in cluster. (default "in-cluster")
  -pass string
        passwd for harbor. (default "Harbor12345")
  -project string
        app harbor projectName. (default "os")
  -registry string
        url for harbor. (default "harbor.xxx.com")
  -user string
        username for harbor. (default "admin")
  -v    version
root@client-go-deploy-68857c6778-d7h9j:/# 
```

| 变量名            | 值                                                        |
| :---------------- | :-------------------------------------------------------- |
| CLUSTER_ENV       | 默认运行在K8s,如需要运行在集群外,传入变量"out-of-cluster" |
| HARBOR_USER       | harbor仓库用户                                            |
| HARBOR_PASS       | harbor仓库密码                                            |
| HARBOR_URL        | harbor的地址,只支持https,如:harbor.xxx.com        |
| HARBOR_PROJECT    | harbor仓库的项目名                                        |
| APPLICATION_GROUP | 要更新的deploy信息, "命名空间.deploy名.容器名"            |

### 示例

```shell
[root@k8s-master01 ]# kubectl create -f  busybox2.yaml
deployment.apps/busybox2 created
[root@k8s-master01 ]# kubectl create -f  busybox.yaml
deployment.apps/busybox created
[root@k8s-master01 ]# 
[root@k8s-master01 ]# kubectl get po --all-namespaces -o wide | grep busybox
default         busybox2-6d55cbb9db-sz8w7                   1/1     Running     0          35s     172.16.85.218    k8s-node01     <none>           <none>
kube-system     busybox-6656b56686-f9s9x                    1/1     Running     0          24s     172.16.85.220    k8s-node01     <none>           <none>
[root@k8s-master01 ]# kubectl create -f  deploy-sample.yaml
clusterrolebinding.rbac.authorization.k8s.io/default-admin created
deployment.apps/client-go-deploy created
[root@k8s-master01 ]#  kubectl get po -n kube-system -owide
NAME                                       READY   STATUS              RESTARTS   AGE     IP               NODE           NOMINATED NODE   READINESS GATES
busybox-6656b56686-f9s9x                   1/1     Running             0          58s     172.16.85.220    k8s-node01     <none>           <none>
calico-kube-controllers-6d7f4d76c7-tltfq   1/1     Running             0          5d19h   192.168.1.106    k8s-master01   <none>           <none>
calico-node-bvnqm                          1/1     Running             0          5d19h   192.168.1.106    k8s-master01   <none>           <none>
calico-node-n7dv6                          1/1     Running             0          5d19h   192.168.1.107    k8s-master02   <none>           <none>
calico-node-twwzr                          1/1     Running             0          5d19h   192.168.1.108    k8s-node01     <none>           <none>
client-go-deploy-68857c6778-d7h9j          0/1     ContainerCreating   0          1s      <none>           k8s-master02   <none>           <none>
coredns-76945f5f5c-lh98g                   1/1     Running             0          2d21h   172.16.85.223    k8s-node01     <none>           <none>
coredns-76945f5f5c-t2hjc                   1/1     Running             0          2d21h   172.16.122.188   k8s-master02   <none>           <none>
elasticsearch-logging-0                    1/1     Running             0          5d19h   172.16.32.141    k8s-master01   <none>           <none>
elasticsearch-logging-1                    1/1     Running             0          5d19h   172.16.122.139   k8s-master02   <none>           <none>
filebeat-29bzl                             1/1     Running             0          17h     192.168.1.106    k8s-master01   <none>           <none>
filebeat-2mv6l                             1/1     Running             0          17h     192.168.1.108    k8s-node01     <none>           <none>
filebeat-brxf8                             1/1     Running             0          17h     192.168.1.107    k8s-master02   <none>           <none>
kibana-logging-765ddf7f6-mscfs             1/1     Running             4          5d19h   172.16.32.142    k8s-master01   <none>           <none>
metrics-server-v0.3.6-6b55f4b546-jkp26     2/2     Running             0          5d19h   172.16.32.131    k8s-master01   <none>           <none>
pushgateway-bfbd7995-26ch2                 1/1     Running             0          44h     172.16.32.132    k8s-master01   <none>           <none>
[root@k8s-master01 ]# kubectl get po --all-namespaces -o wide | grep busybox
default         busybox2-6d55cbb9db-sz8w7                   1/1     Terminating   0          74s     172.16.85.218    k8s-node01     <none>           <none>
default         busybox2-7fcf754467-x29bh                   1/1     Running       0          4s      172.16.122.147   k8s-master02   <none>           <none>
kube-system     busybox-6656b56686-f9s9x                    1/1     Terminating   0          63s     172.16.85.220    k8s-node01     <none>           <none>
kube-system     busybox-9d78c7f79-wz7d6                     1/1     Running       0          5s      172.16.122.148   k8s-master02   <none>           <none>
[root@k8s-master01 ]#
```

如果集群开启了RBAC,还需要授权

```shell
kubectl create -f rbac.yaml
```

```shell
[root@k8s-master01 ]# kubectl logs -f client-go-deploy-68857c6778-d7h9j -n kube-system 
[INFO] 2020/09/04 06:54:12 mian.go:42: namespace: kube-system  deployment: busybox container: alpine
[INFO] 2020/09/04 06:54:12 harbor.go:48: The new tag of the os/alpine is: 3.12, The corresponding index is: 1 , Image is: harbor.junengcloud.com/os/alpine:3.12
[INFO] 2020/09/04 06:54:12 deployment.go:24: Found deployment
[INFO] 2020/09/04 06:54:12 deployment.go:26: container name -> busybox
[INFO] 2020/09/04 06:54:12 deployment.go:34: Current container image -> harbor.xxx.com/os/debian:stretch-slim
[INFO] 2020/09/04 06:54:12 deployment.go:35: New image of harbor registry -> harbor.xxx.com/os/alpine:3.12
[INFO] 2020/09/04 06:54:12 mian.go:42: namespace: default  deployment: busybox2 container: alpine
[INFO] 2020/09/04 06:54:12 harbor.go:48: The new tag of the os/alpine is: 3.12, The corresponding index is: 1 , Image is: harbor.junengcloud.com/os/alpine:3.12
[INFO] 2020/09/04 06:54:12 deployment.go:24: Found deployment
[INFO] 2020/09/04 06:54:12 deployment.go:26: container name -> busybox2
[INFO] 2020/09/04 06:54:12 deployment.go:34: Current container image -> harbor.xxx.com/os/debian:stretch-slim
[INFO] 2020/09/04 06:54:12 deployment.go:35: New image of harbor registry -> harbor.xxx.com/os/alpine:3.12
```

<img src="https://github.com/ProdanLabs/Golang-practice-project/blob/master/image/qrcode_for_gh.jpg" width="120">

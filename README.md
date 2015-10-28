## Kameni

Listen to Marathon hooks and update vulcand backends.

### Example config:

```toml
# kameni listening port
addr = ":7373"

# marathon url
marathon = "http://localhost:8080"

# etcd servers used for kameni and vulcand
etcd_servers = [
  "http://localhost:4001"
]

kameni_prefix = "kameni"
vulcand_prefix = "vulcand"
```

### Running

Configure marathon for http callbacks, and set kameni address callback, eg(`http://localhost:7373/marathon_callback`) in marathon http callbacks param (`http_endpoints`)

```
$ kameni -f kameni.toml

# Or setup upstart file from example kameni-upstart.conf
$ sudo service kameni start
```

### Setup apps in etcd:

```
# my-app => app id in marathon
# myapp => backend id in vulcand
$ etcdctl set /kameni/apps/my-app '{"backend_id": "myapp"}'
{"backend_id": "my-app"}

$ # scale apps in marathon

$ ./etcdctl get /vulcand/backends/myapp/servers/my-app.3faf04e6-674f-11e5-bee1-56847afe9799
{"URL":"http://mesos-slave05.myapp.com:31509"}
```


### Setup using confd for nginx

```
# nginx-app => app id in marathon
# /nginx/apps/nginx-app/upstreams => nginx upstream
$ etcdctl set /kameni/apps/my-app '{"backend_id": "/nginx/apps/nginx-app/upstreams", "type": "confd"}'
{"backend_id": "/nginx/apps/nginx-app/upstreams", "type": "confd"}

$ # scale apps in marathon

$ ./etcdctl get /nginx/apps/nginx-app/upstreams/my-app.3faf04e6-674f-11e5-bee1-56847afe9799
mesos-slave05.myapp.com:31509
```

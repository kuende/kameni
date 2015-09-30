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

```
$ kameni -f kameni.toml

# Or setup upstart file from example kameni-upstart.conf
$ sudo service kameni start
```

### Setup apps in etcd:

```
# my-app => app id in marathon
# myapp => backend id for vulcand
$ etcdctl set /kameni/apps/my-app '{"backend_id": "myapp"}'
{"backend_id": "my-app"}

$ # scale apps in marathon

$ ./etcdctl get /vulcand/backends/ng/servers/my-app.3faf04e6-674f-11e5-bee1-56847afe9799
{"URL":"http://my.marathon.hostname:31509"}
```

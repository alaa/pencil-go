# Pencil

Pencil is a simple service-discovery tool that meant to work with Docker and Consul.
It basically syncronize the "diff" between the local state (Running Docker Containers)
and the remote state on Consul registry every (n) seconds. the default is set to 5.

Pencil never does a bulk syncing but it only syncs the changes wheather they are (additions or deletions)
on the consul registry which is important for external Load-balancing or service-monitoring.

## Running Pencil (Recommended way)

```
$ docker run -d \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -e CONSUL_HTTP_ADDR=[consul-api-address] alaa/pencil-go
```

## Testing Pencil locally

- ``` https://github.com/alaa/pencil-go.git ```

- Run Consul cluster (4 nodes) on your machine:

- ``` scripts/start_consul_cluster ```

- Start Pencil

- ``` CONSUL_HTTP_ADDR=[consul-api-address] go run main.go ```

- Run few Nginx containers:

- ``` docker run -P nginx ```

Watch the changes on the stdout and on the consul web-ui.

## TODO

- Add customizable consul health-checks

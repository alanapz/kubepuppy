# CLUSTERPUPPY

- Git: https://github.com/alanapz/clusterpuppy
- Docker: https://hub.docker.com/repository/docker/alanmpinder/gomaster/general

Test Go application, used to test the following technologies:

- Basic Go REST API (using Gin framework)
- Calls Kubernetes API via Go client library
- Kubernetes user/role/certificate management
- Front-end using Knockout
- Testing statically-compiled application with a `FROM scratch` image
- Testing Helm subchart handling


### Testing statically-compiled application

Image is statically built: Only 1 binary and assets (CSS/HTML).
(We use `from scratch` to force no base image).

To compile: `CGO_ENABLED=0 go build -ldflags="-extldflags=-static"`

In Dive:

```
Permission     UID:GID       Size  Filetree
drwxr-xr-x         0:0      11 MB  └── app
drwxr-xr-x         0:0      11 kB      ├── assets
-rwxr-xr-x         0:0      237 B      │   ├── gomaster.css
-rwxr-xr-x         0:0      11 kB      │   └── index.html
-rwxr-xr-x         0:0      11 MB      └── gomaster
```

For a total image size of 11MB.

### [Testing Helm subchart organisation](https://helm.sh/docs/chart_template_guide/subcharts_and_globals/)

The parent chart (gomaster) is completely empty, with no templates.

- `gomaster-deployment` subchart stores deployment.yaml
- `gomaster-service` subchart stores service.yaml


All settings are stored in values.yaml of parent chart.

See `chart\gomaster\values.yaml` for more details.

# To Run

Via Docker:

```
$ docker run --init -p 8080:8080 docker.io/alanmpinder/gomaster
```

Then visit http://localhost:8080 in your browser.


Via K8S:

(Assuming cluster is setup and kubeconfig available)

```
$ helm install release1 chart\clusterpuppy -n clusterpuppy --create-namespace
kubectl port-forward -n gomaster service/gomaster-release1-service 8080:http
```

Then visit http://localhost:8080 in your browser.

# To Build (Docker)

From the root folder:

```
$ cd app
$ docker build -f Dockerfile gomaster
```

# To Build (from source)

From the root folder:

```
$ cd app\gomaster
$ go build .

# Or, to build a static binary
# (Dont forget to copy assets folder if necessary)
$ CGO_ENABLED=0 go build -ldflags="-extldflags=-static"
```


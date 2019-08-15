# pubsub-push

**pubsub-push** is a small utility for receiving **Pull** messages from a subscription in **Google PubSub** and redirect them as an **HTTP Push** to a **local endpoint**.

It is necessary because using the normal **PubSub Push** you need a *public* and *verified* endpoint. Using this tool for **development purposes** you will be able to receive **HTTP Push** locally.


## Installation

You can just download the latest binary from [Releases](https://github.com/klassmann/pubsub-push/releases) tab and put it inside your *PATH*.

If you have a recent **go** installation, you can use:

```
$ go get github.com/klassmann/pubsub-push
$ go install github.com/klassmann/pubsub-push
```

## Usage

Before running, you have to set `GOOGLE_APPLICATION_CREDENTIALS` environment variable with the path of your service account key created with the correct Roles for getting messages from the desired subscription. [See more here](https://cloud.google.com/pubsub/docs/access-control) and [here](https://cloud.google.com/iam/docs/service-accounts).

Example in Linux/Mac:
```
$ export GOOGLE_APPLICATION_CREDENTIALS=$HOME/project/key.json
```

After that you should run the following command with your configurations.
**These arguments are required:**
```
$ pubsub-push -project gcloud-project -sub subscription_name -endpoint http://localhost
```

If you want to add some headers to request, you can use the `-header` option.

    Note that -header can be used multiple times and you must use the format "key=value".


Adding headers:
```
$ pubsub-push -project proj-id -sub name -endpoint localhost -header "Content-type=application/json" -header "Auth-Key=my-secret-key"
```


### Arguments
```
-project        The ID of your project inside Google Cloud Platform.
                    Eg: -project my-cloud-project

-sub            The subscription name, not including namespace, only the name.
                    Eg: -sub topic_sub

-endpoint       The complete URL, including schema, domain, port and the path.
                    Eg: -endpoint http://localhost:5000/services/sync

-header         A string like "key=value" for request headers. It can be used
                multiplem times.
                    Eg: -header "Content-type=application/json"
```



## Requirements and Go version

It was built with `go1.12.5 linux/amd64` and you can see dependencies inside [go.mod](go.mod).


## License
[Apache 2.0](LICENSE)
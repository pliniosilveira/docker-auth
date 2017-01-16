Docker Authorization plugin Example

Building
-
To build this plugin at least Go 1.7 is required
```sh
$ git clone https://github.com/daehyeok/docker-auth
$ cd docker-auth
$ make
```

Running & testing
-
`docker daemon` have to restart after plugin binary executed, with  `--authorization-plugin=docker-auth-plugin` command line flags.

```sh
$ docker-auth-plugin &
$ docker daemon --authorization-plugin=docker-auth-plugin
```

This plugin allow/deny create new container via web server

```sh
$ curl localhost:8080/status
New Container create is Unblocked
$ docker run -it --name helloworld hello-world #work

$ curl localhost:8080/block
New Container create is Blocked

$ docker run -it --rm hello-world #create new container is blocked
docker: Error response from daemon: authorization denied by plugin docker-auth-plugin: Create New Container Blocked.
See 'docker run --help'.

$ docker start helloworld #existing container is still work
helloworld

$ curl localhost:8080/unblock
New Container create is Unblocked
```

License
-
MIT

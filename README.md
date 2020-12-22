# simpleLB
A simple and lite load balancer(support automatic access to certificates from Let's Encrypt )

## Install
```shell script
go get -u github.com/si9ma/simpleLB
```

## Usage

```shell script
A simple and lite load balancer(support automatic access to certificates from Let's Encrypt )

Usage:
  simpleLB [command]

Available Commands:
  help        Help about any command
  lb          start a load balancer

Flags:
      --config string   config file (default is $HOME/.simpleLB.yaml)
  -h, --help            help for simpleLB
  -t, --toggle          Help message for toggle

Use "simpleLB [command] --help" for more information about a command.
```

## TODO

- [ ] Support Round-robin load balancing strategy
- [ ] Downstream health check
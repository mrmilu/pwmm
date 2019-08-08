# PWMM - Pingdom window maintenance manager

## Description
Simple program written in GO to simulate windows maintenances pausing and unpausing checks of pingdom. 

It is my first program in GO and I have make it rushed, so less complaints and more PR's!!

## install
You can clone this repo, install go and build it or you can download [the linux / mac precompiled binary](https://github.com/mrmilu/pwmm/releases/tag/v0.0.1).

Also you can use a docker image

## build
#### required
You need go > 1.11 to resolve dependencies automatically or docker

You can use older versions but you must install the dependencies manually before build. 

to build run:
```
bash build.sh
```

## how to use

You must define a yaml configuration with pingdom credentials and events (windows maintenance)
check _example.yml_

```yaml
credentials:
  apikey: pingdom api key
  user: pingdom user
  password: pingdom password
```

The windows maintenances are defined like events with three params: the exact name of pingdom check, and the start date & end dates. 
Always use the timezone of the OS that execute the program. 

**IMPORTANT respect the time format!**

```yaml
events:
  - name: "my website - home"
    startdate: "08-08-2019 03:55"
    finishdate: "08-08-2019 08:10"
```

## docker

You can build your own image
```
# build image
docker build -t pwmm .

# run command with custom config in current dir
docker run --rm -d -v $(pwd):/workdir/config pwmm -f myconfig.yml
```
or use the image uploaded to official public repository of docker
```
docker run --rm -d -v $(pwd):/workdir/config mrmiludevops/pwmm -f myconfig.yml 
```

# TODO
- add tests
- automate compilation

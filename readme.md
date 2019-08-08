# PWMM - Pingdom window maintenance manager

## Description
Simple program written in GO to simulate windows maintenances pausing and unpausing checks of pingdom. 

It is my first program in GO and I have make it rushed, so less complaints and more PR's!!

## install
You can clone this repo, install go and build it or you can download the linux precompiled binary.


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

# TODO
- docker image
- add tests
- automate compilation

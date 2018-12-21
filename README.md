[![Build Status](https://travis-ci.org/GiedriusS/SAM.svg?branch=master)](https://travis-ci.org/GiedriusS/SAM)

# SAM
Similar Alerts Manager - provides a way to see what alerts were firing at the same time historically

# Purpose
The goal of SAM is to implement a similar kind of functionality that New Relic has when it shows what other alerts were firing around the same time.

# Features
* Can show all related alerts of any alert
* Can show the information an alert according to its hash
* Shows when the internal database was last updated
* Has a persistence layer through Redis


# Command line arguments
| Name            | Default value | Purpose                                           | Example               |
|-----------------|---------------|---------------------------------------------------|-----------------------|
| --elasticsearch |               | ES instance                                       | http://127.0.0.1:1234 |
| --addr          | :9888         | API listen address                                | 0.0.0.0:1111          |
| --redis         |               | Redis instance                                    | 127.0.0.1:5555        |
| --cacheinterval | 5             | Interval between cache uploads in seconds         | 25                    |
| --esinterval    | 10            | ES check interval in seconds                      | 33                    |
| --esindex       | alertmanager  | ES index name                                     | foobar                |

# Architecture
![architecture](https://github.com/GiedriusS/SAM/raw/master/SAM.png "SAM architecture")


# Main use-case
* Go to http://127.0.0.1:9888/alert/DiskEmpty?dev=sda2&instance=example1 to retrieve related alerts to that alert and its labels;
* Then find out: `{"d3c96c8359aa53c97154786d3b3fad3df6ac8ba557a64e4f1b782a586c105e6f":2}`. This means that the alert with that hash has been firing two times in the past at the same time as this alert;
* Go to http://127.0.0.1:9888/hash/d3c96c8359aa53c97154786d3b3fad3df6ac8ba557a64e4f1b782a586c105e6f and you will get back the information about it.

# Running
A Docker image is available with SAM. Pull it down and run it with:

```
$ docker pull stag1e/sam:latest
$ docker run --rm -it -p 9888:9888 stag1e/sam --elasticsearch 'http://127.0.0.1:1234' --redis '127.0.0.1:3333'
```
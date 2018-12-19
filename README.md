# SAM
Similar Alerts Manager - provides a way to see what alerts were firing at the same time historically

# Example deployment
![architecture](https://github.com/GiedriusS/SAM/raw/master/SAM.png "SAM architecture")

# Purpose
The goal of SAM is to implement a similar kind of functionality that New Relic has when it shows what other alerts were firing around the same time.

# Main use-case
* Start up SAM
* Go to http://127.0.0.1:9888/alert/DiskEmpty?dev=sda2&instance=example1 to retrieve related alerts to that label set
* Then find out: {"d3c96c8359aa53c97154786d3b3fad3df6ac8ba557a64e4f1b782a586c105e6f":2}
* Go to http://127.0.0.1:9888/hash/d3c96c8359aa53c97154786d3b3fad3df6ac8ba557a64e4f1b782a586c105e6f and you will get back the information about it

# Command line arguments
| Name            | Default value | Purpose                                           | Example               |
|-----------------|---------------|---------------------------------------------------|-----------------------|
| --elasticsearch |               | Specify ES instance                               | http://127.0.0.1:1234 |
| --addr          | :9888         | Specify API listen address                        | 0.0.0.0:1111          |
| --redis         |               | Specify Redis instance                            | 127.0.0.1:5555        |
| --cacheinterval | 5             | Specify interval between cache uploads in seconds | 25                    |
| --esinterval    | 10            | Specify ES check interval in seconds              | 33                    |
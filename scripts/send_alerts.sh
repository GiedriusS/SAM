#!/bin/bash
alerts1='[
  {
    "labels": {
       "alertname": "DiskRunningFull",
       "dev": "sda1",
       "instance": "example1"
     },
     "annotations": {
        "info": "The disk sda1 is running full",
        "summary": "please check the instance example1"
      }
  }
]'
alerts2='[
  {
    "labels": {
       "alertname": "DiskEmpty",
       "dev": "sda2",
       "instance": "example1"
     },
     "annotations": {
        "info": "The disk sda2 is running empty",
        "summary": "please check the instance example1",
        "runbook": "the following link http://test-url should be clickable"
      }
  }
]'
curl -XPOST -d"$alerts1" http://localhost:9093/api/v1/alerts
sleep 3
curl -XPOST -d"$alerts2" http://localhost:9093/api/v1/alerts

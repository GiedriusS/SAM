#!/bin/bash
data='
{
  "template": "alertmanager-2*",
  "settings": {
    "number_of_shards": 1,
    "number_of_replicas": 1,
    "index.refresh_interval": "1s",
    "index.query.default_field": "groupLabels.alertname"
  },
  "mappings": {
    "_default_": {
      "_all": {
        "enabled": false
      },
      "properties": {
        "@timestamp": {
          "type": "date",
          "doc_values": true
        }
      },
      "dynamic_templates": [
        {
          "string_fields": {
            "match": "*",
            "match_mapping_type": "string",
            "mapping": {
              "type": "string",
              "index": "not_analyzed",
              "ignore_above": 1024,
              "doc_values": true
            }
          }
        }
      ]
    }
  }
}
'

curl -XPUT -d "${data}" -H "Content-Type: application/json" 'http://127.0.0.1:9200/_template/alertmanager-2*'

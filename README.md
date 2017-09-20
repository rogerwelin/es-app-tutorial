
## Create index
curl -XPUT 'localhost:9200/cartoon?pretty' -H 'Content-Type: application/json' -d' { "settings" : { "index" : { "number_of_shards" : 3, "number_of_replicas" : 1 } } } '

## Import data
curl -XPOST localhost:9200/cartoon/_bulk?pretty=true --data-binary @dataset/anime.json

## View index status after import
curl -XGET localhost:9200/_cat/indices?v

## Make a simple search
curl -XGET 'localhost:9200/cartoon/_search?q=name:berserk&pretty'

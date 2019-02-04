
## Pull/run image
```shell
$ docker run -d -p 9200:9200 -p 9300:9300 -e "discovery.type=single-node" docker.elastic.co/elasticsearch/elasticsearch:6.2.3
```

## Create index
```shell
$ curl -XPUT 'localhost:9200/cartoon?pretty' -H 'Content-Type: application/json' -d' { "settings" : { "index" : { "number_of_shards" : 3, "number_of_replicas" : 1 } } } '
```

## Import data
```shell
$ curl -XPOST -H "Content-Type: application/json" -XPOST --data-binary @dataset/anime.json http://localhost:9200/cartoon/_bulk?pretty=true
```

## View index status after import
```shell
$ curl -XGET localhost:9200/_cat/indices?v
```

## Make a simple search
```shell
$ curl -XGET 'localhost:9200/cartoon/_search?q=name:berserk&pretty'
```

## Start the Go backend
```shell
$ cd backend && go build && ./main
```

## Open the frontend and start typing anime titles
**On Mac**
```shell
$ open index.html
```

**On Linux**
```shell
$ google-chrome index.htnl
```

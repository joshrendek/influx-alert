## How to Use

* `name`: the name of the alert ( will be used in notifier )
* `interval`: how often to check influxdb (in seconds)
* `timeshift`: how far back to go (query is like: `where time > now() - TIMESHIFT`
* `limit`: the max number of results to return
* `type`: influxdb (the only option for now)
* `function`: min/max/average are the only supported functions for now
* `query`: the influxb query to run (omit any limit or where clause on the time)
* `trigger`: the type of trigger and value that would trigger it
  * `operator`: gt/lt
  * `value`: value to compare against (note all values are floats internally)
* `notifiers`: an array of notifiers, possible options are slack and hipchat

Example: ( see example.yml for more )

``` yml
- name: Not Enough Foo
  type: influxdb
  function: average
  timeshift: 1h
  limit: 10
  interval: 10
  query: select * from "foo.counter"
  notifiers:
      - slack
      - hipchat
      - foobar
  trigger:
  operator: lt
  value: 10
```


## Environment Variables
```
  * SLACK_API_TOKEN
  * SLACK_ROOM
  * HIPCHAT_API_TOKEN
  * HIPCHAT_ROOM_ID
  * HIPCHAT_SERVER (optional)
  * DEBUG (optional)
```

## Supported Notifiers

* HipChat ( hosted and private servers )
* Slack ( Generate slack token: https://api.slack.com/web )

## Supported Backends

* InfluxDB v0.9

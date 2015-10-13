## Influx Alert

This is a tool to alert on data that is fed into
InfluxDB (for example, via statsd) so you can get alerted on it.

## How to get it

Go to releases, or download the latest here: [v0.1](https://github.com/joshrendek/influx-alert/releases/download/0.1/influx-alert)

## How to Use

* `name`: the name of the alert ( will be used in notifier )
* `interval`: how often to check influxdb (in seconds)
* `timeshift`: how far back to go (query is like: `where time > now() - TIMESHIFT`
* `limit`: the max number of results to return
* `type`: influxdb (the only option for now)
* `function`: min/max/average are the only supported functions for now
* `query`: the influxdb query to run (omit any limit or where clause on the time)
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
  * INFLUX_HOST
  * INFLUX_PORT (8086 is default)
  * INFLUX_DB
  * INFLUX_USER
  * INFLUX_PASS
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

## TODO

* [ ] Tests
* [ ] Email Notifier
* [ ] Web notifier (POST?)


## License

```
The MIT License (MIT)

Copyright (c) 2015 Josh Rendek

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
```

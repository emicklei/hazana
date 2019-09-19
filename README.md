# hazana - package for creating load tests of services

[![Build Status](https://travis-ci.org/emicklei/hazana.png)](https://travis-ci.org/emicklei/hazana)
[![GoDoc](https://godoc.org/github.com/emicklei/hazana?status.svg)](https://godoc.org/github.com/emicklei/hazana)

Hazana is created for load tests that use (generated) clients in Go to communicate to services (in any supported language).
By providing the Attack interface, any client and protocol could potentially be tested with this package.
This package was created to load test gRPC services.

Compared to existing HTTP load testing tools (e.g. tsenart/vegeta) that can send raw HTTP requests, this package requires the use of client code to send the requests and receive the response.

## Attack

        // Attack must be implemented by a service client.
        type Attack interface {
                // Setup should establish the connection to the service
                // It may want to access the config of the runner.
                Setup(c Config) error

                // Do performs one request and is executed in a separate goroutine.
                // The context is used to cancel the request on timeout.
                Do(ctx context.Context) DoResult

                // Teardown can be used to close the connection to the service.
                Teardown() error

                // Clone should return a fresh new Attack
                // Make sure the new Attack has values for shared struct fields initialized at Setup.
                Clone() Attack
        }
The **hazana** runner will spawn goroutines to meet this load.
Each goroutine will use one Attack value to perform the communication ( see **Do()** ).
Typically each Attack value uses its own connection but your implementation can use another strategy.

### Rampup

The **hazana** runner will use a rampup period in which the RPS is increased (every second) during the rampup time. In this phase, new goroutines are spawned up to the given maximum. This package has two strategies for adding new attackers to meet the rps.
The **linear** rampup strategy will create exactly the maximum number of goroutines within the rampup period. 
The **exp2** strategy spawn goroutines as needed (exponential with max factor of 2) to match the current rps load during.

![](hazana_profile.png)

### Flags

Programs that use the **hazana** package will have several flags to control the load runner.

    Usage of <<your load test program>>:
        -attack int
                duration of the attack in seconds (default 60)
        -max int
                maximum concurrent attackers (default 10)
        -timeout int
                timeout in seconds for an attack call (default 5)
        -o string
                output file to write the metrics per sample request index (use stdout if empty)
        -ramp int
                ramp up time in seconds (default 10)
        -s string
                set the rampup strategy, possible values are {linear,exp2}
        -rps int
                target number of requests per second, must be greater than zero (default 1)
        -t int
                test your attack implementation with a number of sample calls. Your program exits after this.
        -verbose
                produce more verbose logging

#### Example

After creating your implementation type **YourAttack** then this would be the minimal program to run a load test.

        func main() {
                r := hazana.Run(new(YourAttack), hazana.ConfigFromFlags())

                // inspect the report and compute whether the test has failed
                // e.g by looking at the success percentage and mean response time of each metric.
                r.Failed = false

                hazana.PrintReport(r)
        }

### Configuration

In addition to using flags, you can load the configuration from a JSON file. Values set with flags will override those from the configuration file.

        {
                "RPS": 10,
                "AttackTimeSec": 20,
                "RampupTimeSec": 10,
                "MaxAttackers": 10,
                "OutputFilename": "myreport.json",
                "Verbose": true,
                "Metadata": {
                        "service" : "happiness.services.com",
                        "environment" : "staging",
                        "version": "v1.42",
                        "apiToken*": "your-secret-token"
                }
        }

_Note that metadata keys that end with * will be obfuscated when reporting_. 

#### Example

        func main() {
                r := hazana.Run(YourAttack{}, hazana.ConfigFromFile("myconfig.json"))
                hazana.PrintReport(r)
                hazana.PrintSummary(r)
        }

See **examples/zombie.go** for a complete minimal example.

See **examples/clock** for an actual gRPC service that can tell time under load.

### Sample verbose output from one of our services

        2017/08/17 10:26:32 hazana - load runner
        2017/08/17 10:26:32 begin rampup of [10] seconds
        2017/08/17 10:26:32 setup and spawn new attacker [1]
        2017/08/17 10:26:34 current rate [1.998330838214118], target rate [1], attackers [1], mean response time [229.576185ms]
        2017/08/17 10:26:34 setup and spawn new attacker [2]
        2017/08/17 10:26:36 current rate [2.9997667261403085], target rate [2], attackers [2], mean response time [52.674453ms]
        2017/08/17 10:26:36 setup and spawn new attacker [3]
        2017/08/17 10:26:37 current rate [3.9896134642448655], target rate [3], attackers [3], mean response time [52.805457ms]
        2017/08/17 10:26:37 setup and spawn new attacker [4]
        2017/08/17 10:26:39 current rate [4.973899789133906], target rate [4], attackers [4], mean response time [54.645931ms]
        2017/08/17 10:26:39 setup and spawn new attacker [5]
        2017/08/17 10:26:40 current rate [5.995198996681865], target rate [5], attackers [5], mean response time [53.096359ms]
        2017/08/17 10:26:40 setup and spawn new attacker [6]
        2017/08/17 10:26:41 current rate [6.973412713416999], target rate [6], attackers [6], mean response time [55.183152ms]
        2017/08/17 10:26:41 setup and spawn new attacker [7]
        2017/08/17 10:26:42 current rate [7.982915515035891], target rate [7], attackers [7], mean response time [45.521208ms]
        2017/08/17 10:26:42 setup and spawn new attacker [8]
        2017/08/17 10:26:43 current rate [8.953025436248573], target rate [8], attackers [8], mean response time [42.844917ms]
        2017/08/17 10:26:43 setup and spawn new attacker [9]
        2017/08/17 10:26:44 current rate [9.982111816054946], target rate [9], attackers [9], mean response time [42.128101ms]
        2017/08/17 10:26:44 setup and spawn new attacker [10]
        2017/08/17 10:26:45 current rate [10.99217377013479], target rate [10], attackers [10], mean response time [37.483798ms]
        2017/08/17 10:26:45 end rampup ending up with [10] attackers
        2017/08/17 10:26:45 begin full attack of [10] remaining seconds
        2017/08/17 10:26:55 end full attack
        2017/08/17 10:26:55 stopping attackers [10]
        2017/08/17 10:26:55 tearing down attackers [10]
        {
                "startedAt": "2017-08-17T10:26:32.976273638+02:00",
                "finishedAt": "2017-08-17T10:26:55.961789195+02:00",
                "configuration": {
                        "rps": 10,
                        "attackTimeSec": 20,
                        "rampupTimeSec": 10,
                        "maxAttackers": 10,
                        "rampupStrategy" : "linear",
                        "verbose": true,
                        "doTimeoutSec": 5,
                        "metadata": {
                                "service" : "happiness.services.com",
                                "environment" : "staging",
                                "version": "v1.42",
                                "apiToken*": "***---***---***"
                        }
                },
                "metrics": {
                        "POST item.xml": {
                                "latencies": {
                                        "total": 3817277924,
                                        "mean": 37794830,
                                        "50th": 32147032,
                                        "95th": 46125381,
                                        "99th": 71243508,
                                        "max": 422720083
                                },
                                "earliest": "2017-08-17T10:26:45.924789988+02:00",
                                "latest": "2017-08-17T10:26:55.929145547+02:00",
                                "end": "2017-08-17T10:26:55.961659257+02:00",
                                "duration": 10004355559,
                                "wait": 32513710,
                                "requests": 101,
                                "rate": 10.095602800636126,
                                "success": 1,
                                "status_codes": null,
                                "errors": null
                        }
                },
                "failed":false
        }

### Stackdriver integration

The [hazana-stackdriver-monitoring](https://github.com/emicklei/hazana-stackdriver-monitoring) project offers a tool to send the results of a loadtest to a Google Stackdriver account. The metrics from the load test are sent as custom metrics to Stackdriver Monitoring. The report itself is sent as a log entry to Stackdriver Logging.

## Graph visualization

The [hazana-report-visualizer](https://github.com/robertalpha/hazana-report-visualizer) is a tool that produces a diagram served by a local webapp that visualizes a set of reports. It parses the JSON documents to collect the data points.

The [hazana-grafana-monitoring](https://github.com/emicklei/hazana-grafana-monitoring) package sends data to a Graphite server which data can be visualised using a Grafana dashboard. Using the "-m" flag you can tell your running loadtest to send this data in realtime to the dashboard (via Graphite).

© 2017-2019, [ernestmicklei.com](http://ernestmicklei.com).  Apache v2 License. Contributions welcome.

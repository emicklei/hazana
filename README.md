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

#### linear

The **linear** rampup strategy will create exactly the maximum number of goroutines within the rampup period.

#### exp2

The **exp2** strategy spawn goroutines as needed (exponential with max factor of 2) to match the current rps load during. You can change the time in seconds to measure the rate (default=1) using the `keep` parameter. The factor can be changed with the `max-factor` parameter.
Using the configuration:

        “RampupStrategy”: “exp2 keep=5 max-factor=1.1",

As a command line flag:

        -s “exp2 keep=5 max-factor=1.1"

![profile](hazana_profile.png)

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
        -csv string
                CSV output file to write the metrics
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

#### Example from flags

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
                "RampupStrategy": "linear",
                "MaxAttackers": 10,
                "DoTimeoutSec": 5,
                "OutputFilename": "myreport.json",
                "Verbose": true,
                "Debug": false,
                "Metadata": {
                        "service" : "happiness.services.com",
                        "environment" : "staging",
                        "version": "v1.42",
                        "apiToken*": "your-secret-token"
                }
        }

_Note that metadata keys that end with * will be obfuscated when reporting_. 

#### Example from file

        func main() {
                r := hazana.Run(YourAttack{}, hazana.ConfigFromFile("myconfig.json"))
                hazana.PrintReport(r)
                hazana.PrintSummary(r)
        }

See **examples/zombie.go** for a complete minimal example.

See **examples/clock** for an actual gRPC service that can tell time under load.

### Sample verbose output from one of our services

        +1s - *** Hazana load runner ready to attack ***
        +1s - rps [20] attack [90] rampup [30] strategy [exp2 keep=5 max-factor=1.1] max [10] timeout [60] JSON [] CSV [report.csv]
        +1s - [8] available logical CPUs
        +1s - ||| rampup of [30] seconds to RPS [20] within attack of [90] seconds
        +1s - setup and spawn new attacker [1]
        +10s - rate [0.999929 -> 1], mean response [4.200527174s], requests [2], attackers [1], success [50 %]
        +10s - setup and spawn new attacker [2]
        +17s - rate [1.496180 -> 1], mean response [2.144609037s], requests [3], attackers [2], success [66 %]
        +22s - rate [0.624308 -> 2], mean response [2.294694897s], requests [7], attackers [2], success [100 %]
        +22s - setup and spawn new attacker [3]
        +28s - rate [1.912042 -> 2], mean response [726.317443ms], requests [10], attackers [3], success [100 %]
        +28s - setup and spawn new attacker [4]
        +33s - rate [0.895940 -> 3], mean response [2.401489741s], requests [5], attackers [4], success [100 %]
        +33s - setup and spawn new attacker [5]
        +39s - rate [1.468714 -> 4], mean response [2.27204947s], requests [16], attackers [5], success [93 %]
        +39s - setup and spawn new attacker [6]
        +44s - rate [3.861015 -> 4], mean response [1.179716954s], requests [22], attackers [6], success [95 %]
        +44s - setup and spawn new attacker [7]
        +49s - rate [4.998839 -> 5], mean response [1.186741846s], requests [25], attackers [7], success [100 %]
        +49s - setup and spawn new attacker [8]
        +55s - rate [3.441637 -> 6], mean response [1.799351987s], requests [18], attackers [8], success [100 %]
        +55s - setup and spawn new attacker [9]
        +1m0s - rate [3.976226 -> 6], mean response [2.055811128s], requests [25], attackers [9], success [100 %]
        +1m0s - setup and spawn new attacker [10]
        +1m5s - rate [3.555044 -> 7], mean response [2.143760751s], requests [24], attackers [10], success [95 %]
        +1m10s - rate [3.146333 -> 8], mean response [2.618258609s], requests [15], attackers [10], success [100 %]
        +1m15s - rate [4.270554 -> 8], mean response [2.266823792s], requests [28], attackers [10], success [100 %]
        +1m21s - rate [4.259347 -> 9], mean response [2.315611968s], requests [22], attackers [10], success [100 %]
        +1m26s - rate [5.568545 -> 10], mean response [2.139091741s], requests [20], attackers [10], success [95 %]
        +1m31s - rate [2.886864 -> 10], mean response [2.436964999s], requests [25], attackers [10], success [100 %]
        +1m36s - rate [5.345010 -> 11], mean response [1.871945022s], requests [26], attackers [10], success [100 %]
        +1m41s - rate [3.830200 -> 12], mean response [1.887682054s], requests [23], attackers [10], success [100 %]
        +1m47s - rate [3.146157 -> 12], mean response [2.172577651s], requests [27], attackers [10], success [96 %]
        +1m52s - rate [4.188442 -> 13], mean response [2.130385868s], requests [25], attackers [10], success [100 %]
        +1m57s - rate [3.782267 -> 14], mean response [2.001996278s], requests [27], attackers [10], success [100 %]
        +2m2s - rate [4.101091 -> 14], mean response [1.654729749s], requests [24], attackers [10], success [95 %]
        +2m7s - rate [3.609923 -> 15], mean response [2.414714412s], requests [24], attackers [10], success [100 %]
        +2m13s - rate [3.681610 -> 16], mean response [2.316453726s], requests [23], attackers [10], success [100 %]
        +2m18s - rate [4.531789 -> 16], mean response [1.65280437s], requests [34], attackers [10], success [100 %]
        +2m23s - rate [4.410283 -> 17], mean response [2.037735812s], requests [24], attackers [10], success [95 %]
        +2m28s - rate [3.834626 -> 18], mean response [2.006833405s], requests [22], attackers [10], success [95 %]
        +2m34s - rate [3.846774 -> 18], mean response [2.391088483s], requests [26], attackers [10], success [100 %]
        +2m40s - rate [3.027653 -> 19], mean response [1.998581075s], requests [21], attackers [10], success [95 %]
        +2m46s - rate [1.628317 -> 20], mean response [3.36120323s], requests [18], attackers [10], success [94 %]
        +2m46s - ||| rampup ending up with [10] attackers
        +2m46s - begin full attack of [60] remaining seconds
        +3m46s - end full attack
        +3m46s - stopping attackers [10]
        +3m46s - tearing down attackers [10]
        +3m46s - CSV report written to [report.csv]
        ---------
        category-c
        - - - - -
        requests: 4
        errors: 0
        rps: 4.1229347790648765
        mean: 4.571788432s
        50th: 4.222491097s
        95th: 4.883595542s
        99th: 4.883595542s
        avg kB >: 0
        avg kB <: 19
        max: 5.184646443s
        success: 100 %
        ---------
        product-p
        - - - - -
        requests: 2
        errors: 0
        rps: 2.319250297795507
        mean: 3.39725241s
        50th: 2.817941596s
        95th: 2.817941596s
        99th: 2.817941596s
        avg kB >: 0
        avg kB <: 11
        max: 3.976563224s
        success: 100 %
        ---------
        search
        - - - - -
        requests: 4
        errors: 0
        rps: 0.8711446902078789
        mean: 3.802252054s
        50th: 4.087427907s
        95th: 4.897371961s
        99th: 4.897371961s
        avg kB >: 0
        avg kB <: 8
        max: 5.792558602s
        success: 100 %

### Stackdriver integration

The [hazana-stackdriver-monitoring](https://github.com/emicklei/hazana-stackdriver-monitoring) project offers a tool to send the results of a loadtest to a Google Stackdriver account. The metrics from the load test are sent as custom metrics to Stackdriver Monitoring. The report itself is sent as a log entry to Stackdriver Logging.

## Graph visualization

The [hazana-report-visualizer](https://github.com/robertalpha/hazana-report-visualizer) is a tool that produces a diagram served by a local webapp that visualizes a set of reports. It parses the JSON documents to collect the data points.

The [hazana-grafana-monitoring](https://github.com/emicklei/hazana-grafana-monitoring) package sends data to a Graphite server which data can be visualised using a Grafana dashboard. Using the "-m" flag you can tell your running loadtest to send this data in realtime to the dashboard (via Graphite).

© 2017-2022, [ernestmicklei.com](http://ernestmicklei.com).  Apache v2 License. Contributions welcome.

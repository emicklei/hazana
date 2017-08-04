# hazana - package for creating load tests of services

[![Build Status](https://travis-ci.org/emicklei/hazana.png)](https://travis-ci.org/emicklei/hazana)

Hazana was created to create load tests that use (generated) gRPC clients in Go to communicate to gRPC services (in any language). However, by providing the Attack interface, any client and protocol could potentially be tested with this package.

Compared to existing HTTP load testing tools (e.g. tsenart/vegeta) that can send raw HTTP requests, this package requires the use client code to perform the request. Consequently, time to send a request and receive a response includes time spent on marshalling that request and unmarshalling a response.

### Attack

    // Attack must be implemented by a service client.
    type Attack interface {
        // Setup should establish the connection to the service
        // It may want to access the config of the runner.
        Setup(c Config) error
        // Do performs one request and is executed in one fixed goroutine.
        Do() DoResult
        // Teardown should close the connection of the service
        Teardown() error
        // Clone should return a new fresh Attack
        Clone() Attack
    }
    
Depending on the target RPS (requests per second), the **hazana** runner will spawn goroutines to meet this load.
Each goroutine will use one Attack value to perform the communication ( see **Do()** ). Typically each Attack value uses its own connection but your implementation can use another strategy.

### Rampup
The **hazana** runner can use a rampup period in which the RPS is increased (per second) during the rampup time. In this phase, new goroutines could be spawned if the actual rate is lower than the increasing target.

© 2017, [ernestmicklei.com](http://ernestmicklei.com).  Apache v2 License. Contributions welcome.

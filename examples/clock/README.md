# load testing a clock service

## start the gRPC service

    make server

## run the load test

    make load

## example output

    2017/11/02 14:42:09 hazana - load runner
    2017/11/02 14:42:09 begin rampup of [10] seconds using strategy [exp2]
    2017/11/02 14:42:09 setup and spawn new attacker [1]
    2017/11/02 14:42:11 current rate [1.9926321390489952], target rate [1], attackers [1], mean response time [2.016872ms], requests [2], errors [0]
    2017/11/02 14:42:13 current rate [1.9914542734133827], target rate [1], attackers [1], mean response time [434.36µs], requests [2], errors [0]
    2017/11/02 14:42:15 current rate [1.9974484653227422], target rate [1], attackers [1], mean response time [448.466µs], requests [2], errors [0]
    2017/11/02 14:42:17 current rate [1.9988431555318338], target rate [1], attackers [1], mean response time [380.016µs], requests [2], errors [0]
    2017/11/02 14:42:19 current rate [1.9989110811950523], target rate [1], attackers [1], mean response time [427.215µs], requests [2], errors [0]
    2017/11/02 14:42:21 current rate [1.9993368259728352], target rate [1], attackers [1], mean response time [413.588µs], requests [2], errors [0]
    2017/11/02 14:42:23 current rate [1.991473474455543], target rate [1], attackers [1], mean response time [433.84µs], requests [2], errors [0]
    2017/11/02 14:42:25 current rate [1.9920863771043618], target rate [1], attackers [1], mean response time [501.35µs], requests [2], errors [0]
    2017/11/02 14:42:27 current rate [1.9977612727580443], target rate [1], attackers [1], mean response time [484.231µs], requests [2], errors [0]
    2017/11/02 14:42:29 current rate [2.997741246936552], target rate [2], attackers [1], mean response time [504.39µs], requests [3], errors [0]
    2017/11/02 14:42:29 end rampup ending up with [1] attackers
    2017/11/02 14:42:29 begin full attack of [0] remaining seconds
    2017/11/02 14:42:29 end full attack
    2017/11/02 14:42:29 stopping attackers [1]
    2017/11/02 14:42:29 stopped attacker  0
    2017/11/02 14:42:29 tearing down attackers [1]
    2017/11/02 14:42:29 teared down attacker  0
    {
        "startedAt": "2017-11-02T14:42:29.351258666+01:00",
        "finishedAt": "2017-11-02T14:42:29.351294837+01:00",
        "configuration": {
            "rps": 2,
            "attackTimeSec": 10,
            "rampupTimeSec": 10,
            "rampupStrategy": "exp2",
            "maxAttackers": 10,
            "verbose": true,
            "doTimeoutSec": 5
        },
        "error": "",
        "metrics": {}
    }

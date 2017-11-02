# load testing a clock service

## start the gRPC service

    make server

## run the load test

    make load

## example output

    2017/11/02 14:45:26 hazana - load runner
    2017/11/02 14:45:26 begin rampup of [10] seconds using strategy [exp2]
    2017/11/02 14:45:26 setup and spawn new attacker [1]
    2017/11/02 14:45:28 current rate [1.991313607278519], target rate [1], attackers [1], mean response time [2.108787ms], requests [2], errors [0]
    2017/11/02 14:45:30 current rate [1.9919558905032486], target rate [1], attackers [1], mean response time [508.721µs], requests [2], errors [0]
    2017/11/02 14:45:32 current rate [1.9965004202628391], target rate [1], attackers [1], mean response time [473.721µs], requests [2], errors [0]
    2017/11/02 14:45:34 current rate [1.999651920590183], target rate [1], attackers [1], mean response time [425.944µs], requests [2], errors [0]
    2017/11/02 14:45:36 current rate [1.9962680905919197], target rate [1], attackers [1], mean response time [442.761µs], requests [2], errors [0]
    2017/11/02 14:45:38 current rate [1.9993932241443368], target rate [1], attackers [1], mean response time [466.836µs], requests [2], errors [0]
    2017/11/02 14:45:40 current rate [1.9991322266787044], target rate [1], attackers [1], mean response time [581.987µs], requests [2], errors [0]
    2017/11/02 14:45:42 current rate [1.9997057992833978], target rate [1], attackers [1], mean response time [529.287µs], requests [2], errors [0]
    2017/11/02 14:45:44 current rate [1.99803303438127], target rate [1], attackers [1], mean response time [442.995µs], requests [2], errors [0]
    2017/11/02 14:45:46 current rate [2.9929693145511753], target rate [2], attackers [1], mean response time [443.452µs], requests [3], errors [0]
    2017/11/02 14:45:46 end rampup ending up with [1] attackers
    2017/11/02 14:45:46 begin full attack of [10] remaining seconds
    2017/11/02 14:45:56 end full attack
    2017/11/02 14:45:56 stopping attackers [1]
    2017/11/02 14:45:56 stopped attacker  0
    2017/11/02 14:45:56 tearing down attackers [1]
    2017/11/02 14:45:56 teared down attacker  0
    {
        "startedAt": "2017-11-02T14:45:46.255293779+01:00",
        "finishedAt": "2017-11-02T14:45:56.260565813+01:00",
        "configuration": {
            "rps": 2,
            "attackTimeSec": 20,
            "rampupTimeSec": 10,
            "rampupStrategy": "exp2",
            "maxAttackers": 10,
            "verbose": true,
            "doTimeoutSec": 5
        },
        "error": "",
        "metrics": {
            "now": {
                "latencies": {
                    "total": 17256594,
                    "mean": 821742,
                    "50th": 473803,
                    "95th": 651139,
                    "99th": 710850,
                    "max": 7151028
                },
                "earliest": "2017-11-02T14:45:46.255303826+01:00",
                "latest": "2017-11-02T14:45:56.259998927+01:00",
                "end": "2017-11-02T14:45:56.26047273+01:00",
                "duration": 10004643101,
                "wait": 473803,
                "requests": 21,
                "rate": 2.0990254013060174,
                "success": 1,
                "status_codes": {},
                "errors": null
            }
        }

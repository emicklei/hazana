# attacking a site using HTTP requests

    go run loadrun.go -v

## example output

2017/11/07 21:41:08 hazana - load runner
2017/11/07 21:41:08 begin rampup of [10] seconds using strategy [exp2]
2017/11/07 21:41:08 setup and spawn new attacker [1]
2017/11/07 21:41:10 current rate [1.9955431798042864], target rate [1], attackers [1], mean response time [319.071381ms], requests [2], errors [0]
2017/11/07 21:41:12 current rate [1.9945783650123206], target rate [1], attackers [1], mean response time [80.455101ms], requests [2], errors [0]
2017/11/07 21:41:14 current rate [1.9963918491409411], target rate [1], attackers [1], mean response time [79.463091ms], requests [2], errors [0]
2017/11/07 21:41:16 current rate [1.9991594773876462], target rate [1], attackers [1], mean response time [80.79111ms], requests [2], errors [0]
2017/11/07 21:41:18 current rate [1.994817311025024], target rate [1], attackers [1], mean response time [81.079856ms], requests [2], errors [0]
2017/11/07 21:41:20 current rate [1.999034336478249], target rate [1], attackers [1], mean response time [79.284879ms], requests [2], errors [0]
2017/11/07 21:41:22 current rate [1.99660324478755], target rate [1], attackers [1], mean response time [81.743933ms], requests [2], errors [0]
2017/11/07 21:41:24 current rate [1.9968518811384581], target rate [1], attackers [1], mean response time [78.591967ms], requests [2], errors [0]
2017/11/07 21:41:26 current rate [1.9986990328086516], target rate [1], attackers [1], mean response time [76.761118ms], requests [2], errors [0]
2017/11/07 21:41:28 current rate [1.9991410890225125], target rate [1], attackers [1], mean response time [81.871875ms], requests [2], errors [0]
2017/11/07 21:41:28 end rampup ending up with [1] attackers
2017/11/07 21:41:28 begin full attack of [50] remaining seconds
2017/11/07 21:42:18 end full attack
2017/11/07 21:42:18 stopping attackers [1]
2017/11/07 21:42:18 tearing down attackers [1]
2017/11/07 21:42:18 ---------
2017/11/07 21:42:18 ubanita.org
2017/11/07 21:42:18 - - - - -
2017/11/07 21:42:18 requests: 51
2017/11/07 21:42:18      rps: 1.019954188043676
2017/11/07 21:42:18     mean: 80.030592ms
2017/11/07 21:42:18     95th: 87.842548ms
2017/11/07 21:42:18      max: 94.756054ms
2017/11/07 21:42:18   errors: 0
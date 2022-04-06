# attacking a site using HTTP requests

    go run loadrun.go -verbose

## example output

```
     +0s - *** Hazana load runner ready to attack ***
     +0s - rps [1] attack [60] rampup [10] strategy [exp2 keep=1 factor=2.0] max [10] timeout [5] JSON [] CSV []
     +0s - [12] available logical CPUs
     +0s - ||| BEGIN rampup of [10] seconds to RPS [1] within attack of [60] seconds
     +0s - setup and spawn new attacker [1]
     +2s - rate [1.999777 -> 1], mean response [36.58042ms], requests [2], attackers [1], success [100 %]
     +4s - rate [1.993485 -> 1], mean response [15.253384ms], requests [2], attackers [1], success [100 %]
     +6s - rate [1.998952 -> 1], mean response [16.089334ms], requests [2], attackers [1], success [100 %]
     +8s - rate [1.999890 -> 1], mean response [20.862207ms], requests [2], attackers [1], success [100 %]
    +10s - rate [1.999784 -> 1], mean response [15.822663ms], requests [2], attackers [1], success [100 %]
    +12s - rate [1.999822 -> 1], mean response [24.013889ms], requests [2], attackers [1], success [100 %]
    +14s - rate [1.999661 -> 1], mean response [16.7976ms], requests [2], attackers [1], success [100 %]
    +16s - rate [1.999736 -> 1], mean response [16.842378ms], requests [2], attackers [1], success [100 %]
    +18s - rate [1.999700 -> 1], mean response [19.061471ms], requests [2], attackers [1], success [100 %]
    +20s - rate [1.996123 -> 1], mean response [16.840931ms], requests [2], attackers [1], success [100 %]
    +20s - ||| rampup ENDing up with [1] attackers
    +20s - BEGIN full attack of [50] remaining seconds
  +1m10s - END full attack
  +1m10s - stopping attackers [1]
  +1m10s - tearing down attackers [1]
---------
github.com
- - - - -
requests: 51
  errors: 0
     rps: 1.0199964515547446
    mean: 16.668581ms
    50th: 15.2158ms
    95th: 24.331013ms
    99th: 33.707276ms
avg kB >: 0
avg kB <: 0
     max: 33.707276ms
 success: 100 %
```

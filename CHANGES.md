[2022-04-06] v1.10:

- replace quantile package

[2020-06-11] v1.9:

- add CSV output
- add parameters for rampup strategy exp2
- logging show relative time
- able to abort test, rampup or full attack, preserving report
- add debug option
- stop attackers on quit signal
- set label for timeout result
- update uber ratelimit
- report bytes transferred in/out

v1.6.1:
- 50th and 99th percentile response times to the result summary (Geurtje)

1.6.0:
- rename verbose flag "v" to "verbose"

1.5.2:
- add ReadFile utility

1.5.1:
- improve summary

1.5.0 (and 1.3.0):
- split RunLifecycle into two interface. Pass report to AfterRun

1.2.0:
- added RunLifecycle interface for BeforeRun and AfterRun hooks

v1.1.0:
- BREAKING API: change signature of Attacker.Do to include context

v1.0.0:
- initial

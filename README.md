# promlint

This project provides a `metriclint` package used to make standard metrics for your applications. 

The rules and implementation mostly based on [Prometheus promlint](github.com/prometheus/client_golang/prometheus/testutil/promlint).

## Why
[Prometheus promlint](github.com/prometheus/client_golang/prometheus/testutil/promlint) used for check metrics at the end, 
always in the E2E test. It's hard to cover all metrics.

`metriclint` intended to be a supplement of [Prometheus promlint](github.com/prometheus/client_golang/prometheus/testutil/promlint),
it helps check your metric at the development phase, especially when the metric registering to a registry.

## Future
Reserve a place to donate it to [Prometheus promlint](github.com/prometheus/client_golang/prometheus/testutil/promlint) if
 it works fine after some experiment.
  
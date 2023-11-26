faultgauge
==========

`faultgauge` creates a new gauge which samples the fail rate over the
window length. Be sure to call IncrementFail() or IncrementSuccess() to
inform the gauge of the success or failure. A windowLength of 10 seconds
or more is recommended, and 60 seconds is even better.

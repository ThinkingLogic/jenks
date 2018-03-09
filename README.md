# Jenks Natural Breaks

A Golang implementation of the [Jenks natural breaks optimization](http://en.wikipedia.org/wiki/Jenks_natural_breaks_optimization) algorithm.

Ported from a [javascript version](https://gist.github.com/tmcw/4977508)
\- itself ported from Fortran and described here by
[Tom MacWright](https://macwright.org/2013/02/18/literate-jenks.html).

[![Build Status](https://travis-ci.org/ThinkingLogic/jenks.svg?branch=master)](https://travis-ci.org/ThinkingLogic/jenks)
[![Coverage Status](https://coveralls.io/repos/github/ThinkingLogic/jenks/badge.svg)](https://coveralls.io/github/ThinkingLogic/jenks)

## Usage

```
import 	"github.com/ThinkingLogic/jenks"

//...

data := []float64{1, 2, 3,  12, 13, 14,  21, 22, 23,  27, 28, 29}

breaks := jenks.NaturalBreaks(data, 4)
// [1, 12, 21, 27]
```

## License
This software is Licenced under the [MIT License](LICENSE.md).


## Changes from the original JS version
Please note that the javascript version this was ported from is broken,
specifically this line:
```
kclass[countNum - 1] = data[lower_class_limits[k][countNum] - 2];
```
should be:
```
kclass[countNum - 1] = data[lower_class_limits[k][countNum] - 1];
```
\- it has been fixed in this version, along with a number of minor improvements
(such as not returning the upper bound so that the length of the returned slice
 matches the requested number of classes).


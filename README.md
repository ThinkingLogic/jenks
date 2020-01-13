# Jenks Natural Breaks

A Golang implementation of the [Jenks natural breaks optimization](http://en.wikipedia.org/wiki/Jenks_natural_breaks_optimization) algorithm.
It is a data clustering algorithm designed to determine the best arrangement of values into different classes,
seeking to reduce the variance within classes and maximize the variance between classes.

Ported from a [javascript version](https://gist.github.com/tmcw/4977508)
\- itself ported from Fortran and described here by
[Tom MacWright](https://macwright.org/2013/02/18/literate-jenks.html).

[![Build Status](https://travis-ci.org/ThinkingLogic/jenks.svg?branch=master)](https://travis-ci.org/ThinkingLogic/jenks)
[![Coverage Status](https://coveralls.io/repos/github/ThinkingLogic/jenks/badge.svg)](https://coveralls.io/github/ThinkingLogic/jenks)

## Usage

```
import 	"github.com/ThinkingLogic/jenks"

//...

data := []float64{1.1, 2.1, 3.1,  12.1, 13.1, 14.1,  21.1, 22.1, 23.1,  27.1, 28.1, 29.1}

breaks := jenks.NaturalBreaks(data, 4)
// [1.1, 12.1, 21.1, 27.1]

rounded := jenks.Round(breaks, data)
// [0, 10, 20, 27]

allBreaks := jenks.AllNaturalBreaks(data, 4)
// [ [1.1, 21.1]
//   [1.1, 12.1, 21.1]
//   [1.1, 12.1, 21.1, 27.1] ]
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
\- it has been fixed here, along with a number of minor improvements
(such as not returning the upper bound so that the length of the returned slice
 matches the requested number of classes).


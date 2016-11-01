# grayt

[![Build Status](https://travis-ci.org/peterstace/grayt.svg?branch=master)](https://travis-ci.org/peterstace/grayt)

Go RAY Tracer

## Known Issues

* Elapsed time follows clock time rather than execution time. This causes
  elapsed time to be misleading when the computer sleeps.

* Only supports triangles. This leads to slow runtimes. At least should also
  have sphere, aligned box, aligned square, and aligned plane.

# grayt

[![Build Status](https://travis-ci.org/peterstace/grayt.svg?branch=master)](https://travis-ci.org/peterstace/grayt)

Go RAY Tracer

## Features

- [X] Path tracing via rendering equation simulation (Monte Carlo method).
- [X] Diffuse reflections (matte surfaces).
- [ ] Specular reflections (mirror surfaces).
- [ ] Light transmission (transparent surfaces).
- [X] Depth of field effects.
- [ ] Multithreading support.
- [X] Fast acceleration structure.
- [ ] Automatic Lambda detection for grid data structure.

## TODO

- Use terminology for focal distance / focal ratio correctly.
- Use Ginkgo and Gomega for testing.

## Gallery

### Cornell Box

![Cornell Box](/gallery/out_q100000.png)

### Split Box

![Split Box](/gallery/splitbox[wnNQ9molxVk]_720x720_q5000.png)

### Sphere Tree

![Sphere Tree](/gallery/sphere_tree[1Fq0zJkUpGk]_820x410_q5000.png)

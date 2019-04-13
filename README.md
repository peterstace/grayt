# grayt

Go RAY Tracer

## Features

- [X] Path tracing via rendering equation simulation (Monte Carlo method).
- [X] Diffuse reflections (matte surfaces).
- [X] Specular reflections (mirror surfaces).
- [ ] Light transmission (transparent surfaces).
- [X] Depth of field effects.
- [X] Multithreading support.
- [X] Fast acceleration structure.
- [X] Web UI.
- [ ] Persistent storage of partial renders.

## TODO

- Ability to delete renders.
- Use fixed space font for data in UI.
- Default scene display ratio.
- Load passes statistics from accumulator.
- Calculate resolutions in backend.
- Allow to downsample resolution.
- Allow to choose exposure level.
- Try different lambda values for grid.
- Bounding Volume Hierarchy
- Use pointer to material instead of copying in each object.

## Gallery

### Split Box

![Split Box](/gallery/splitbox[KQRdZO3e8KI]_1024x1024_q100000.png)

### Sphere Tree

![Sphere Tree](/gallery/sphere_tree[T35GCh3Lpj4]_1024x512_q100000.png)

### Cornell Box

![Cornell Box](/gallery/out_q100000.png)

## Scene Ideas

- 3D patchwork: https://mattdesl.svbtle.com/pen-plotter-2

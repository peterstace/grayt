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
- [X] Persistent storage of partial renders.

## Architecture

Package layout:

| package name     | purpose                                                                                           |
| scene            | hold types that define the scene, e.g. Scene, Object, Material, Surface                           |
| scene/dsl        | DSL to make build scenes easier, intended as a dot import                                         |
| scene/library    | holds a registry for all of the scenes                                                            |
| scene/cornellbox | holds cornellbox specific scene factories. Can have other packages similarly to hold other scenes |
| api              | API/server code                                                                                   |
| control          | orchestrates tracing, accepts commands from the API, handles persistence                          |
| trace            | tracing logic, acceleration structures, geometry                                                  |
| colour           | colour type                                                                                       |
| xmath            | math types such as vect and ray                                                                   |


## TODO

- Use terminology for focal distance / focal ratio correctly.
- Voxel Geometry
- Ability to delete renders.
- Don't use sync xhr when posting new scene.
- Use fixed space font for data in UI.
- Put a link to the github repo in the source code.
- Don't load scenes into memory unless they have active workers.

## Gallery

### Split Box

![Split Box](/gallery/splitbox[KQRdZO3e8KI]_1024x1024_q100000.png)

### Sphere Tree

![Sphere Tree](/gallery/sphere_tree[T35GCh3Lpj4]_1024x512_q100000.png)

### Cornell Box

![Cornell Box](/gallery/out_q100000.png)

## Scene Ideas

- 3D patchwork: https://mattdesl.svbtle.com/pen-plotter-2

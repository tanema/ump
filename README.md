# Ump [![](https://godoc.org/github.com/tanema/ump?status.svg)](http://godoc.org/github.com/tanema/ump)

Go collision-detection library for axis-aligned rectangles. Its main features are:

* `ump` only does axis-aligned bounding-box (AABB) collisions. If you need anything
  more complicated than that (circles, polygons, etc.) give [box2dlite](https://github.com/neguse/go-box2d-lite/box2dlite) a look.
* Handles tunnelling - all items are treated as "bullets". The fact that we only
  use AABBs allows doing this fast.
* Strives to be fast while being economic in memory
* It's centered on *detection*, but it also offers some (minimal & basic) *collision response*
* Can also return the items that touch a point, a segment or a rectangular zone.
* `ump` is _gameistic_ instead of realistic.

The demos are Amore based, but this library can be used in any Go program.

`ump` is ideal for:

* Tile-based games, and games where most entities can be represented as axis-aligned
  rectangles.
* Games which require some physics, but not a full realistic simulation - like
  a platformer.
* Examples of genres: top-down games (Zelda), Shoot-em-ups, fighting games (Street
  Fighter), platformers (Super Mario).

`ump` is not a good match for:

* Games that require polygons for the collision detection
* Games that require highly realistic simulations of physics - things "stacking
  up", "rolling over slides", etc.
* Games that require very fast objects colliding reallistically against each other
  (in ump, being _gameistic_, objects are moved and collided _one at a time_)
* Simulations where the order in which the collisions are resolved isn't known.

## Example

For full example usage please see the [amore](https://github.com/tanema/amore) example
[platformer](https://github.com/tanema/amore-examples/tree/master/platformer)

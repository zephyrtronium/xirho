# xirho/encoding/flame

Package flame implements parsing the XML-based Flame format used by flam3 and Apophysis.

By default, the parser recognizes the following variations and their parameters:

- linear (always treated as 3D)
- linear3D
- bipolar
- blur
- pre_blur
- bubble
- elliptic (not identical to Apophysis's elliptic)
- curl
- cylinder
- disc
- exp
- expo
- flatten
- foci
- gaussian_blur
- post_heat
- hemisphere
- julia (not tested and probably wrong)
- julian
- lazysusan
- log
- mobius
- mobiq
- noise
- polar
- rod (rod_blur parameter is ignored)
- scry
- spherical
- spherical3D
- pre_spherical
- splits
- splits3D
- unpolar

## Adding variations

The decoder ignores any xform attributes it doesn't recognize. To teach it to understand more variation types, create a Parser function and add it to the Funcs map. If it has variables, add them to the KnownAttrs map.

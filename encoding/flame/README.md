# xirho/encoding/flame

Package flame implements parsing the XML-based Flame format used by flam3 and Apophysis.

By default, the parser recognizes the following variations and their parameters:

- linear (always treated as 3D)
- linear3D
- blur
- pre_blur
- bubble
- elliptic (not identical to Apophysis's elliptic)
- disc
- flatten
- julia (not tested and probably wrong)
- julian
- mobius
- mobiq
- polar
- spherical
- spherical3D
- pre_spherical
- splits
- splits3D

## Adding variations

The decoder ignores any xform attributes it doesn't recognize. To teach it to understand more variation types, create a Parser function and add it to the Funcs map.

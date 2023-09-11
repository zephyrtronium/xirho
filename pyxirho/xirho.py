#!/usr/bin/env python3

"""
Utilities for manipulating xirho histograms.

pyxirho is intended for experimentation with xirho image conversion utilities.
It contains facilities to read raw histogram dumps as numpy arrays, as well as
implementations of xirho's color conversion and tone mapping routines.

As pyxirho exists for experimentation, it is not guaranteed to stay in sync
with xirho.

"""

import math
import struct
import typing

import numpy as np

class ToneMap(typing.NamedTuple):

    """Tone mapping parameters.

    Attributes:
        brightness: Multiplicative scaling factor.
        gamma: Gamma correction factor.
        gamma_min: Threshold at which to apply gamma correction.
    """

    brightness: float = 1
    gamma: float = 1
    gamma_min: float = 0

class Histogram:

    """Histogram object.
    
    Wraps the functionality in module xirho in a manner aware of the
    histogram's shape and tone mapping parameters.
    """

    hist: np.ndarray
    tone_map: ToneMap
    lqa: float
    osa: int

    def __init__(self,
                 hist: np.ndarray,
                 osa: int = 1,
                 iters: int = 25000,
                 proj_area: float = 1,
                 tone_map: typing.Optional[ToneMap] = None):
        self.hist = hist
        self.tone_map = tone_map or ToneMap()
        self.lqa = lqa(hist.size//4, osa, area(hist.shape[:2], proj_area), iters)
        self.osa = osa
    
    def brightest(self):
        """Find the brightest bin, i.e. the highest count."""
        return brightest(self.hist)
    
    def pixel_region(self, x, y):
        """Return the oversampled bins corresponding to an image pixel."""
        osa = self.osa
        return self.hist[osa*x:osa*(x+1), osa*y:osa*(y+1), :]
    
    def pixel(self, bin_x, bin_y):
        """Calculate the color corresponding to a single bin."""
        bin = self.hist[bin_x,bin_y,:]
        return pixel(bin, 65535*self.tone_map.brightness, self.lqa, self.tone_map.gamma, self.tone_map.gamma_min)

# log10(0xffff). Xirho palettes have channels in [0, 0xffff], but the Flame
# algorithm is based on colors in [0, 1]. Subtracting this from log counts
# performs the conversion.
clscale = 4.81647330376524970778

def read(f: typing.BinaryIO) -> np.ndarray:
    """Read a xirho histogram from a file.

    Args:
        f: File containing a xirho histogram dump.
    
    Returns:
        WxHx4 uint64 array containing all histogram counts.
        Pages are channels: in order, R, G, B, N.
    """
    b = f.read(16)
    w, h = struct.unpack('<QQ', b)
    a = np.fromfile(f, dtype=np.uint64, count=w*h*4) # type: np.ndarray
    return np.reshape(a, (w, h, 4), order='F')

def area(hist_shape: tuple[int, ...], proj_area: float) -> float:
    """Calculates Cartesian histogram area.

    Args:
        hist_shape: Shape of the histogram, or at least the width and height.
        proj_area: Projective area of the renderer's linear camera. This can
            be calculated as the determinant of the upper-left 2x2 submatrix
            of the camera matrix.
    """
    w, h = hist_shape
    aspect = w / h
    if aspect > 1:
        aspect = 1 / aspect
    return aspect / proj_area

def lqa(hist_size: int, osa: int, area: float, iters: int) -> float:
    """Calculate log quality-area coefficient.

    Args:
        hist_size: Total number of bins in the histogram.
        osa: Oversampling factor used when rendering.
        area: Cartesian histogram area, as calculated by area.
        iters: Total iterations (usually not hits) during rendering.

    Returns:
        Log quality-area coefficient including adjustment for 16-bit color.
    """
    o = 4 * math.log10(osa)
    a = math.log10(area)
    q = math.log10(hist_size) - math.log10(iters)
    return o - a + q - 2*clscale

def ascale(n: np.uint64, br: float, lqa: float) -> float:
    """Alpha channel scaling.

    Args:
        n: Bin N channel.
        br: Scaled brightness factor.
        lqa: Log quality-area coefficient.

    Returns:
        Scaled alpha channel. The nominal range is [0, 1], but actual values
        may be larger or negative depending on brightness and lqa.
    """
    a = br * (math.log10(n) - lqa)
    return a / n

def gamma(a: float, gamma: float, threshold: float) -> float:
    """Gamma correction.

    Args:
        a: Sample value in [0, 1] (approximately).
        gamma: Gamma correction factor.
        threshold: Gamma threshold.

    Returns:
        Gamma-corrected sample value.
    """
    exp = 1/gamma
    if a >= threshold:
        return a ** exp
    p = a / threshold
    return p * a**exp + (1-p) * threshold**(exp - 1)

def pixel(bin: np.ndarray, br: float, lqa: float, gf: float, thresh: float) -> np.ndarray:
    """Calculate a single pixel value.

    Args:
        bin: 4-element uint64 numpy array containing the sample.
        br: Scaled brightness factor.
        lqa: Log quality-area coefficient.
        gf: Gamma factor.
        thresh: Gamma threshold.

    Returns:
        R, G, B, A channels scaled to a nominal range of [0, 1].
    """
    r, g, b, n = bin
    if n == 0:
        return [0., 0., 0., 0.]
    a = ascale(n, br, lqa)
    ag = gamma(a, gf, thresh)
    if ag <= 0:
        return [0., 0., 0., 0.]
    s = a / n
    return np.array((r*s, g*s, b*s, ag), dtype=np.float64)

def brightest(hist: np.ndarray) -> np.ndarray:
    """Find the brightest bin, i.e. the one with the highest count.

    Args:
        hist: WxHx4 histogram as returned by read.

    Returns:
        4-element array containing the bin counts.
    """
    k = np.argmax(hist[:,:,3])
    x, y = divmod(k, hist.shape[1])
    return hist[x,y,:]

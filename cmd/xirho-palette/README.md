# xirho-palette

xirho-palette manipulates the palette format used in xirho JSON files.
When invoked without arguments, it reads formatted alpha-premultiplied color values from standard input and outputs the encoded string on standard output.
When invoked with -d, it reads the encoded string from standard input and outputs formatted alpha-premultiplied color values on standard output.

In both cases, colors are given one per line in R G B A order with whitespace separating each channel.
Channels are formatted as floating point values between 0 and 1.

# Steganography in Go

A toy package to explore [steganography](http://en.wikipedia.org/wiki/Steganography) in Go.

Currently supported:
* embed text in PNG images
* retrieve surreptiously-embedded text from PNG images

## Implementation details

The least significant bit of each color channel is used to store the text data. Greater storage efficiency could be gained, still with no visible artifacts, by using 64-bit PNGs, rather than 32-bit.

The input text is transformed into a channel of 3-bit chunks; each pixel in the image can host one of these chunks.

I decided to not touch the alpha-channel for two reasons: firstly, to make the task more challenging (4-bit chunks fitting UTF-8 better); and secondly because most PNGs have 0xFF throughout for their alpha value: small, random-seeming variations here might be suspicious.

I also choose sub-optimal approaches in a couple of places so as to be able to play with more Go feature. For example, implementing the 3-bit chunks using a channel of 1-bit integers, which is quite wasteful!

One interesting complication I found was that the pixels bordering the edge of the image seemed to be encoded incorrectly by Go's `image/png` encoder. The pixels are slightly lightened during the encoding process. For that reason, I only use the pixels one row and column in from the edges of the image.

## Examples

This _shouldn't_ be interesting, because the text hidden in the image should be imperceptable.

Original:

![](https://github.com/goodgravy/stegano/blob/master/original.png)

With hidden text:

![](https://github.com/goodgravy/stegano/blob/master/output.png)

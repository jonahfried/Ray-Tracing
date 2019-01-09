# Ray-Tracing

rayTracing.go holds the main function.
Build all the .go files and run to render the image.

In a basic sense, the algorithm works by casting various rays through a plane, or "screen". This "screen" is effectively the computer screen. Sending rays through this "screen" they enter the scene and, using intersection algebra, it determines what (if anything) the ray hits. Based on what color object the ray hits, the place where it intersected the screen is colored, and the corresponding pixel on screen is colored.

TODO:
- add complexity to the coloring system (i.e. reflection and diffusion)
- improve ability to translate the camera
- add triangle intersection
- allow for oversampling per pixel to reduce blockiness
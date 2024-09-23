# autogo
The silly Raylib cargame I made in C but in Go


### note Garage texture 
changed saturation to 10% to make Raylib tint work well enough, using ImageMagick:
```convert garage.png -resize 200x200! -channel A -evaluate multiply 0.5 +channel -modulate 100,25,100 garage_200px.png```
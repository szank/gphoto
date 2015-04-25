# gphoto
Package gphoto let you control your digital camera or webcam from your PC.

Gphoto is a go wrapper for libgphoto2 library (http://www.gphoto.org/proj/libgphoto2/).

It supports most of modern cameras supporting PTP protocol and connecting to PC with an USB cable.
Gphoto lets you inspect and modify camera settings ( shutter speed, aperture, iso, image quality etc.),
capture photos and download them to the PC.

You can also capture image previews and by doing it in a loop you can have live view on your computer screen

You will need to have installed libgphoto2 on your system before you wil be able to use this package. To do it on Ubuntu and other debian derivatives you need to type in the terminal : 
```sh
sudo apt get install libgphoto2-6 libgphoto2-port10
```
You might also need to change the  library path in source files to  tell go compiler where to find library files. By default the compiler expects to find libgphoto in 
```sh
/usr/lib/x86_64-linux-gnu
```
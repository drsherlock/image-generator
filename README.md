# Image Generator

### Tool to generate images with overlay and text.
#### Inspired by [this post](https://pace.dev/blog/2020/03/02/dynamically-generate-social-images-in-golang-by-mat-ryer.html) by Mat Ryer.

To use - 

```
go run image_generator.go
```
Arguments - 

```
-file-name
-title
-title-color
```

Example - 
```
go run image_generator.go -file-name=input.jpg -title="TITLE TEXT" -title-color=#aaa
```

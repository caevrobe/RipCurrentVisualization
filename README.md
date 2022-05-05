# Rip Current Visualization

This project provides various functions for visualizing rip currents given .raw scalar volumes and 2D time varying vector fields(???) of shoreline footage.

 - [Results](#results)
 - [Scalar Volume Format](#scalar-volume-format)
 - [Vector Volume Format](#vector-volume-format)
 - [Visualization Methods](#visualization-methods)

<br><br>

## Scalar Volume Format
The scalar volume files contain binary data which includes the RGB values for each pixel for each frame of a given video. Each RGB value is 4 bits long. For example, the following 1 frame video (upscaled below) will produce a scalar volume file as follows:

<div class="column">
   <img src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAIAAAACCAYAAABytg0kAAAAGElEQVQYV2P8z8Dwn5HhPwMjA8P//yAmADgEBf/+zQusAAAAAElFTkSuQmCC" />

   ```
   0xFF 0x00 0x00 0x00 0x00 0xFF
   0x00 0x00 0xFF 0x00 0xFF 0x00
   ```
</div>

<br>

## Vector Volume Format
The vector volume should be a csv file, with the x,y components for each pixel for each frame appearing on their own lines. This csv file can be created by opening the provided vector file in Paraview and exporting.

## Visualization Methods
Scalar data: frame extraction, frame averageing, timestacking.<br>
Vector data: timelines.
The process for initializing scalar, vector, and composite volumes and accessing their functions is provided in main.go.



<style>
   .column {
      display: flex;
      flex-direction: row;
      gap: 1em;
   }

   img {
      width: 200px;
      image-rendering: pixelated;
      filter: brightness(0.95);
   }
</style>
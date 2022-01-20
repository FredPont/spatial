<h1><p align="center">
<br/>
<br/>
<img alt="logo" src="./src/web/logo.png">

</p>
</h1>


Single-Cell Spatial Explorer is a software for user-friendly and versatile exploration of spatial transcriptomic datasets. It is compatible with 
[Single-Cell Signature Explorer](https://doi.org/10.1093/nar/gkz601) (available 
[here](https://sites.google.com/site/fredsoftwares/products/single-cell-signature-explorer)) and 
[Single-Cell Virtual Cytometer](https://doi.org/10.1093/nargab/lqaa025) (available 
[here.](https://github.com/FredPont/single-cell-virtual-cytometer))



For more details see the <b><a href="doc/Manual_Single_Cell_Spatial_Explorer.pdf" target="_blank">Manual</a></b>
inside the doc folder

<!---How to [Cite](https://doi.org/10.1093/nargab/lqaa025) --->

Demo videos can be found in the supplemental data of the [reference article](https://doi.org/10.1093/nargab/lqaa025)

<br/>

## ScreenShots



Single-Cell Spatial Explorer can display clusters :
![clusters](./src/web/clusters.png)

In combination with [Single-Cell Signature Explorer](https://doi.org/10.1093/nar/gkz601), Single-Cell Spatial Explorer can display about 30,000 Human pathways or gene signatures as tunable overlays on a microscopy image. Both transparency and image contrast can be adjusted by the user  :
![expression](./src/web/expression.png)



Single-Cell Spatial Explorer is compatible with [FIJI](https://fiji.sc/) and [ImageJ](https://imagej.nih.gov/ij/). In this picture the contours of tumoral cells have been obtained by FIJI analysis and were imported into Single-Cell Spatial Explorer:
![FIJI](./src/web/fiji.png)

## Prerequisites
A computer with at least 8 Go RAM and a HD display with at least a HD resolution of 1920×1080.

## Installing

To install the software, download the zip archive and unzip it.

 Single-Cell Spatial Explorer is written in pure 
 <a href="https://go.dev/">
    <img src="./src/web/go.jpg" height="30"> 
 </a>
  with the Fyne library 
 <a href="https://github.com/fyne-io">
    <img src="./src/web/fyne.png" height="30">
</a>
statically linked. Precompiled static  binaries are areavailable for Linux, Mac, and Windows. Installation of Single-Cell Spatial Explorer is immediate.

## Documentation

The documentation is available as a PDF file in the "doc" directory

## Single-Cell Spatial Explorer features.

1.  Single-Cell  Spatial  Explorer  is  ready  to  use  in  a  pre-compiled binary, no installation required2.  cross-platform (the interface and the software are codedin pure Go)3.  low memory usage4.  compatible  with  any  PNG  image  associated  with  anyCSV  file  (TAB  separated)  containing  XY  coordinatesof the image.5.  compatible  with  any  numeric  data  :  gene  expression,pathway scores, antibody expression etc...6.  unlimited number of gates.7.  import/export gates in ImageJ/FIJI format.8.  extract cells and sub-tables delimited by the gates on anunlimited number of tables. Exportation is done in CSVfile for great interoperability.9.  2D  plots  of  the  cells  inside  the  gates  with  any  XYcoordinates : t-SNE, UMAP, gene expression, pathwayscores, antibody expression etc...10.  interactive 2D plot to show the selected cells on a t-SNE,UMAP  or  any  other  coordinates  on  the  image  and  tofilter the data tables into sub-tables.11.  cluster display with custom dot opacity and custom dotsize.12.  display  any  kind  of  cell  expression  (genes,  pathways,antibodies...)  with  7  preset  gradients,  custom  legendcolor,  dot  opacity  and  custom  dot  size.  The  gradientsare  simple  two  colors  maps  and  rainbow  colors  mapsTurbo,  Viridis  and  Inferno  to  optimize  accuracy  anddetails visualisation.


## Acknowledgements
Special thanks to Andrew Williams, CTO and technical leader of [Fyne](https://github.com/fyne-io) project, for  his  useful  technical  advises  about  the usage  of  the  [Fyne  graphical  library](https://fyne.io/).

The GO developpment team is aknowledged as well as the contributors of the following projects :

[gonum](https://www.gonum.org/), [fc](https://github.com/ajstarks/fc), [colorgrad](https://github.com/mazznoer/colorgrad), [go-colorful](github.com/lucasb-eyer/go-colorful), [stats](github.com/aclements/go-moremath/stats)

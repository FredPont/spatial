<h1><p align="center">
<br/>
<br/>
<img alt="logo" src="./src/web/logo.png">

</p>
</h1>


Single-Cell Spatial Explorer (scSpatial Explorer) is a software for user-friendly and versatile exploration of spatial transcriptomic datasets. It is compatible with 
[Single-Cell Signature Explorer](https://doi.org/10.1093/nar/gkz601) (available 
[here](https://sites.google.com/site/fredsoftwares/products/single-cell-signature-explorer)) and 
[Single-Cell Virtual Cytometer](https://doi.org/10.1093/nargab/lqaa025) (available 
[here.](https://github.com/FredPont/single-cell-virtual-cytometer))



For more details see the <b><a href="doc/Manual_Single_Cell_Spatial_Explorer.pdf" target="_blank">Manual</a></b>
inside the doc folder

<!---[Contact](mailto:frederic.pont@inserm.fr)--->
<a href="mailto:frederic.pont@inserm.fr?"><img src="src/web/Email.png" height="50"></a>

<!---How to [Cite](https://doi.org/10.1093/nargab/lqaa025) --->

<!---Demo videos can be found in the supplemental data of the [reference article](https://doi.org/10.1093/nargab/lqaa025) --->


## How to cite
Pont, F., Cerapio, J.P., Gravelle, P. et al. Single-cell spatial explorer: easy exploration of spatial and multimodal transcriptomics. BMC Bioinformatics 24, 30 (2023). 
https://doi.org/10.1186/s12859-023-05150-1

## Demo videos

[Quick 2' overview](https://youtu.be/mId538e5JDk)

[Overview](https://youtu.be/dqudL36Dg1M)

[How to install scSpatial Explorer in less than 2'on Windows](https://youtu.be/LBBkN_rJHuc)

[How to install scSpatial Explorer in less than 2'on MacOS](https://youtu.be/tm8lzlP3m_4)

[High throughput data extraction with scSpatial Explorer](https://youtu.be/rSefd0pmc7g)

## <span style="color:red">New !</span>
The new version of Single-Cell Spatial Explorer under development is able to process millions of spots. Thus, it is possible to use the software with technologies such as Akoya PhenoCycler.
The strategy to how to use such datasets is explained in the documentation, read it carefully before starting.

Ki67 expression (2.2 x 10<sup>6</sup> spots)
![Ki67](./src/web/akoyaKi67.jpg)
Keratin 14 expression (2.2 x 10<sup>6</sup> spots)
![K14](./src/web/akoyaKeratin14.jpg)
Cells Clusters (2.2 x 10<sup>6</sup> spots)
![AkoClu1](./src/web/akoyaCluster1.jpg)
Cells cluster layers (2.2 x 10<sup>6</sup> spots)
![AkoClu2](./src/web/akoya_clusters.png)
Cells characterization with a 2D plot (2.2 x 10<sup>6</sup> spots)
![Ako2D1](./src/web/akoya2D1.png)
Cells in gates characterization with a 2D plot (2.2 x 10<sup>6</sup> spots)
![Ako2D1](./src/web/akoyaPlotGate.png)

## ScreenShots

Single-Cell Spatial Explorer is a compiled software, easy to install and to use. It produces beautiful pictures from microcopy images and single cell data sets. Many information can be visualized on top on the microscopy image, such as gene expression, clusters, pathways (screenshot), antibodies... :
![clusters](./src/web/overview.png)


Single-Cell Spatial Explorer can display clusters with preset color gradients. In this example with 11 clusters there is 159667200 possible representations :
![clusters](./src/web/clusters6.png)

In combination with [Single-Cell Signature Explorer](https://doi.org/10.1093/nar/gkz601), Single-Cell Spatial Explorer can display about 29,000 Human pathways or gene signatures (+ genes, antibodies...) as tunable overlays on a microscopy image. Both transparency and image contrast can be adjusted by the user  :
![expression](./src/web/expression6.png)

Single-Cell Spatial Explorer displaying a biological function with different settings of contrast and opacity. A) The original image. B) Contrast was increased with the min/max slider. C) Same image as B) with a wide opacity gradient. D) Same image as C) with a narrow opacity gradient it is possible to display only relevant cells.
![opacity](./src/web/opacity.png)


Screenshot of Single-Cell Spatial Explorer gate comparison. It is possible to compare gates (or groups of gates) across the whole dataset. A volcano plot is obtained and relevant dots can be clicked to reveal the differential genes, antibodies, biological functions or gene signatures. :
![volcano](./src/web/screenshot_1.png)

To localize cells of interest on the microscopy image, Single-Cell Spatial Explorer can draw an interactive 2D plot with any coordinates (t-SNE, UMAP, genes, gene signatures, antibodies ...). On this 2D plot it is possible to select dots, display them on the microscopy image and filter datasets to extract sub-tables and cell names corresponding to these dots:
![FIJI](./src/web/2Dinter.png)

To better characterize cells, Single-Cell Spatial Explorer can plot the gates content on a 2D scatter plot with any coordinates (t-SNE, UMAP, genes, gene signatures, antibodies ...). It is possible to filter datasets to extract sub-tables and cell names corresponding to these dots. Selected cells can be exported to [Single-Cell Virtual Cytometer](https://doi.org/10.1093/nargab/lqaa025) (available 
[here](https://github.com/FredPont/single-cell-virtual-cytometer)) to refine the analysis :
![2Dplot](./src/web/2Dplot.png)


Single-Cell Spatial Explorer can take advantage of years of development in image analysis. It is compatible with [FIJI](https://fiji.sc/) and [ImageJ](https://imagej.nih.gov/ij/). In this picture the contours of tumoral cells have been obtained by FIJI analysis and were imported into Single-Cell Spatial Explorer:
![FIJI](./src/web/fiji.png)


Single-Cell Spatial Explorer is compatible with [Single-Cell Signature Explorer](https://doi.org/10.1093/nar/gkz601) available [here](https://sites.google.com/site/fredsoftwares/products/single-cell-signature-explorer). It is able to display classical (non spatial) scRNAseq data with the help of [Spatial Background Builder](https://github.com/FredPont/Spatial_Background_Builder).![signature explorer](./src/web/signatureExplorer.png)

## Prerequisites
A computer with at least 8 Go RAM and a HD display with at least a HD resolution of 1920×1080.

## Installing

[scSpatial Explorer is installed in less than 2'](https://youtu.be/dqudL36Dg1M). To install the software, download the zip archive and unzip it.

 Single-Cell Spatial Explorer is written in pure 
 <a href="https://go.dev/">
    <img src="./src/web/go.jpg" height="30"> 
 </a>
  with the Fyne library 
 <a href="https://github.com/fyne-io">
    <img src="./src/web/fyne.png" height="30">
</a>
statically linked. Precompiled static  binaries are are available for Linux, Mac, and Windows.

## Documentation

The documentation is available as a PDF file in the "doc" directory

## Single-Cell Spatial Explorer features.

1.  Single-Cell Spatial Explorer has a graphical interface is ready to use in a pre-compiled
    binary, no installation required

2.  Cross-platform (the interface and the software are coded in pure Go)

3.  Low memory usage

4.  Compatible with any PNG image associated with any TAB separated files containing XY coordinates of the image.

5.  Compatible with any numeric data : gene expression, pathway scores, antibody expression etc...

6.  An unlimited number of gates can be drawn in the microscopy image or in the 2D interactive plot.

7.  Import/export gates in ImageJ/FIJI format.

8.  Extract spots and sub-tables delimited by the gates on an unlimited number of tables. Exportation is done in TAB separated files for great interoperability. 

9.  2D plots of the cells inside the gates with any XY coordinates :
    t-SNE, UMAP, gene expression, pathway scores, antibody expression
    etc...

10. Interactive 2D plot to show the selected cells on a t-SNE, UMAP or
    any other coordinates on the image and to filter the data tables
    into sub-tables.

11. Cluster display with 3 color gradients, custom color palette, custom dot opacity and custom dot size. Shuffle color option change color positions on the map leading to almost 2 billions of possible images with 12 clusters.

12. Display any kind of cell expression (genes, pathways, antibodies...)
    with 7 preset gradients, custom legend color, dot opacity and custom
    dot size. The gradients are simple two colors maps and rainbow
    colors maps Turbo, Viridis and Inferno to optimize accuracy and
    details visualization.

13. Min/Max intensity sliders to tune image contrast or remove artifacts due to outliers.

14. Expression opacity gradient with min/max threshold.

15. Slide show to review many cell expression maps without need of
    repetitive click.

16. Screenshot or native resolution image exportation.

17. import and display an unlimited number of cells list by repetitive click on the "import cells" button. The format is directly compatible with Single-Cell Virtual Cytometer .

18. Comparison of two groups of gates together across the whole dataset.

19. Comparison of one group of gates against all the remaining spots.

20. Draw an interactive volcano plot after gate comparison.

21. Plot cell expression of a selected dot in the volcano plot.

22. Export volcano plot and the corresponding data table.

23. Image zoom 10-200%

## Licence
[GNU GPL 3.0](https://www.gnu.org/licenses/gpl-3.0.en.html)

## Acknowledgments
Special thanks to [Andrew Williams](https://andy.xyz), creator of the [Fyne](https://github.com/fyne-io) project and CEO at [Fyne Labs](https://fynelabs.com/), for his useful technical advice about the usage of the [Fyne graphical toolkit](https://fyne.io/).

Thanks to [Miguel Madrid](https://github.com/mimame) for his advices on Github.


The GO developpment team is aknowledged as well as the contributors to Go's GUI ([Fyne](https://github.com/fyne-io)), data and color ecosystem, especially the following projects :

[Go graphics](https://github.com/fogleman/gg), [gonum](https://www.gonum.org/), [fc](https://github.com/ajstarks/fc), [colorgrad](https://github.com/mazznoer/colorgrad), [go-colorful](https://github.com/lucasb-eyer/go-colorful), [stats](https://pkg.go.dev/github.com/aclements/go-moremath/stats),  [pogreb](https://github.com/akrylysov/pogreb), 
[progressbar](https://github.com/schollz/progressbar)

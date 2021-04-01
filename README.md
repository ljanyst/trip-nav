
trip-nav
========

For a reason that I cannot really understand there seems to be no way for Google
Maps navigation function on my iPhone to follow a route defined in the MyMaps
part of Google maps. The mobile app can display the route just fine, but it
cannot follow it.

![Screenshot][screenshot]

This is a small tool that can take the KML file exported from Google Maps and
generate a static website with buttons representing the interest points on the
map The buttons are big enough to press with a motorcycling glove. Pressing a
button fires up the Waze up and sets it up such that it navigates to the
selected interest point.

Example usage:

    ./trip-nav -kml ~/Temp/map/Around-Zürisee.kml -out around-zurisee.html

[screenshot]: screenshot.png

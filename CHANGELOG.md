# 1.0.0-beta4
**Maintainer**: Raffael Sahli <public@raffaelsahli.com>\
**Date**: Tue Mar 24 15:57:21 CET 2020

## Bugfixes
* Do not abort if lookup config in home fails
* Do not abort pull listeners if an error occurs #2
* Do not panic if MongoDB is unreachable during bootstrap

## Changes
* Execute docker container rootless by default
* introduce two config changes: cacheTime => interval / realtime => mod [pull|push]

## Packaging
* Dont build on travis and in dockerfile, only build docker image on travis which contains everything else


# 1.0.0-beta3
**Maintainer**: Raffael Sahli <public@raffaelsahli.com>\
**Date**: Fri Mar 13 16:33:22 CET 2020

## Bugfixes 
* Metrics never updated within startPullListeners? #1
* Proper usage of defaultCacheTime from mongodb.defaultCacheTime
* fixes logging output during testing

## Changes
* go 1.13 update


## Packaging
* correctly push build to github releases from travis
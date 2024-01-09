#!/bin/bash

# URL to be curled
url="https://xddm5bzxkt5zzvw3xr4tv2rbtm0rempb.lambda-url.sa-east-1.on.aws/"

# Looping 100 times
for i in {1..100}
do
   curl "$url"
done

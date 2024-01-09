#!/bin/bash

# URL to be curled
url="<FUNCTION_URL_HERE>"

# Looping 100 times
for i in {1..100}
do
   curl "$url"
done

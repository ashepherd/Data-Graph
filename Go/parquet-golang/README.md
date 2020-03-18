# parquet-golang
Repository to show samples of processing parquet files using Golang

https://github.com/UFOKN/Knowledge-Graph/tree/master/ontologies


So we are really trying to tackle 2 issues right? 

The first is cataloging the data that sits between the current day, and the reanalysis product. This comes from the google bucket youre hitting. The second is storing the operational forecasts.


I think for this we can and should start with a 6 hour forecasts from  from the short range version of the NWM. All operational output can be accessed here:

https://nomads.ncep.noaa.gov/pub/data/nccf/com/nwm/prod/

Today the data paths would be:

https://nomads.ncep.noaa.gov/pub/data/nccf/com/nwm/prod/nwm.20200213/

Lets say its midnight now (eg t**z == t00z). 

I think we want 2 things. The first is the best estimate of the current system which is the analysis and assimilation run:

nwm.t00z.analysis_assim.channel_rt.tm00.conus.nc
where again the t**z is the UTC hour and the tm** is hours ahead
We also want the future forecast which I think we can start with six files (there are 18 available so up to you). Please note the new top directory

https://nomads.ncep.noaa.gov/pub/data/nccf/com/nwm/prod/nwm.20200213/short_range/nwm.t00z.short_range.channel_rt.f001.conus.nc
https://nomads.ncep.noaa.gov/pub/data/nccf/com/nwm/prod/nwm.20200213/short_range/nwm.t00z.short_range.channel_rt.f002.conus.nc
https://nomads.ncep.noaa.gov/pub/data/nccf/com/nwm/prod/nwm.20200213/short_range/nwm.t00z.short_range.channel_rt.f003.conus.nc
https://nomads.ncep.noaa.gov/pub/data/nccf/com/nwm/prod/nwm.20200213/short_range/nwm.t00z.short_range.channel_rt.f004.conus.nc
https://nomads.ncep.noaa.gov/pub/data/nccf/com/nwm/prod/nwm.20200213/short_range/nwm.t00z.short_range.channel_rt.f005.conus.nc
https://nomads.ncep.noaa.gov/pub/data/nccf/com/nwm/prod/nwm.20200213/short_range/nwm.t00z.short_range.channel_rt.f006.conus.nc

And store them as you are, but when the new forecast is released at 1 AM I think there is wisdom in dropping all of those and getting the new set. But we want to keep the A&A run and append the new A&A run.

Is that clear?

Mike


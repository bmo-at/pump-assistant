# Pump Assistant

This project is meant to be a tiny server keeping a record of your fuel consumption.

This is a personal project, therefore your mileage may vary (pun intended). I use it as follows:

My phone has a voice assistant shortcut to record the relevant pieces of data right after I fill up the car (I always brim the tank. Alternatively, I wait with data collection until the tank is brimmed, at which point I add together all the data manually before using the voice shortcut).
This shortcut sends the data via a POST request to the Pump Assistant Server.
Here, it is stored in a local sqlite database.
The data can also be retrieved again via GET Request, with additional features coming up.

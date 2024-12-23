# ATHAN CLI
A simple CLI script that fetches athan times based off your current location

## Setup
- Download athan from release section
- Copy `athan` file into `/usr/local/bin` or its equivalent 

## Todo
- [x] Force refresh of cached files
- [ ] Set home location if location is not found
  - [ ] Manual passing of location
- [ ] Show the location of the Athan times in the table and string
- [ ] Add styling with colours
- [-] Write unit tests for code
  - [ ] Find out how to store a test file 
- [ ] Refractor functions -> no need to pass through `athanCacheJson`


## Notes
Flow
  - If location file doesn't exsit
    - ask user to pass location 
    - default to getting location manually
  
  - If location automatic update daily
  - If manual 
    - Make user update -> provide a prompt for it
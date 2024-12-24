# ATHAN CLI
A simple CLI script that fetches athan times based off your current location

## Setup
- Download athan from release section
- Copy `athan` file into `/usr/local/bin` or its equivalent 

## Todo
- [x] Force refresh of cached files
- [ ] Change from raw json to -> SQLite implementation
- [ ] Set home location if location is not found
  - [x] Manual passing of location
- [ ] Update Athan function to switch between lat-long / city-country
- [ ] Show the location of the Athan times in the table and string
- [ ] Add styling with colours
- [-] Write unit tests for code
  - [ ] Find out how to store a test file 
- [ ] Refractor functions -> no need to pass through `athanCacheJson`
- [ ] Store a list of invalid locations in another json file


## Testing Checklist 
Cache Location
- Location struct is correct
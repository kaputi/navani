package filesystem

// all that has to do with searching/filtering, it goes hand in hand with the watcher,
// the watecher reads all metadata/snipett files and this creates a struct used for search
// it then writes to a json file that is used for searching/filtering that way we dont have to read all files each time the app starts

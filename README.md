# folder-list-go
This sample project lists all folders and files on a path sorted based on their size. there are some optimization that can be made:

### Cache implementation
 Since getting the actual size for big folders is time consuming, a proper caching might help. We need to refresh cache if the folder is modified.

### Implement different sorting
Different sorting can be added

### Better error caching on go routines
We should implement the pipeline to signal go routines to stop in case an error happens. we can save on response time to user but it depends on what is desired.

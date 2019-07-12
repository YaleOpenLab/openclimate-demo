## ipfs

Contains the ipfs interacting parts. When we are adding a file to ipfs, we either could use the javascript handler to call the ipfs api and then use the hash ourselves to decrypt it; otherwise, we need to process a pdf file (i.e. build an xref table) and then convert that into an ipfs file.

### Folder structure

 - ipfs.go

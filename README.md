# cdn.acollier.dev
My personal image hosting server using DigitalOcean Spaces and Go w/ Mux.

To-Do:  
- [x] Add routes to display image links using open graph  
- [ ] Add dashboard with authentication to view, delete and search/filter images  
    - [ ] Add routes to add and delete images to s3 (handled via authentication)  
    - [ ] Add ShareX capabilities for adding new images to s3 (handled via authentication)  
- [ ] An actual CDN? Probably through o

Ideas:  
- [ ] Add support for additional users
    - [ ] Admin users to manage other users, files, account status, etc
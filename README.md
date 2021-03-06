# cdn
My personal image proxy using DigitalOcean Spaces and Go w/ Mux.

To-Do:  
- [x] Add routes to display image links using open graph (ex. [Discord embeds](https://i.imgur.com/otjv4zE.png))  
- [ ] Update README with instructions to configure, build and deploy the project  
    - [ ] Give proper credits to @AlbertSPedersen -  [AlbertSPedersen/s3-discord-embedder](https://github.com/AlbertSPedersen/s3-discord-embedder)  
- [ ] Implement caching from DigitalOcean  
- [ ] Add dashboard with authentication to view, delete and search/filter images  
	- [x] Add authentication  
	- [x] Add view for current images  
    - [x] Add routes to add and delete images to s3 (handled via authentication)  
        - [x] Upload  
        - [x] Delete   
    - [x] Add ShareX capabilities for adding new images to s3 (handled via authentication)  
- [ ] Write a blog post about the development experience  
- [ ] Add Renovate bot  
- [ ] Review dependenices in `go.mod` because of IDE installs and the `// indirect` comments  
- [ ] Add Sentry DSN  

Ideas:  
- [ ] Add support for additional users  
    - [ ] Admin users to manage other users, files, account status, etc  
- [ ] Customize embed (colors, filename, etc)  
- [x] (Re)-Generate access token for API usage (ex. ShareX)  
- [ ] Garbage collection  
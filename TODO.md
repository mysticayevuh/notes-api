# TODO

-   [ ] Implement DELETE /notes/:id endpoint
-   [ ] Add PUT/PATCH endpoint for updating notes
-   [ ] Replace in-memory storage with actual database (postgres)
-   [ ] Add proper error handling and error response types
-   [ ] Write some tests (unit tests at least)

-   [ ] Figure out ID generation strategy (UUIDs are fine for now but might want something else)
-   [ ] Add request validation middleware
-   [ ] Add logging middleware (structured logging would be nice)
-   [ ] Rate limiting
-   [ ] CORS handling if needed

-   [ ] Decide on pagination for GET /notes (probably need it eventually)
-   [ ] Consider adding tags or categories to notes
-   [ ] Maybe add search functionality
-   [ ] Dockerfile for deployment
-   [ ] CI/CD setup

## Project structure

- `cmd`: contains the application-specific code for the executable applications in the project. For now we’ll have just one executable application — the web application — which will live under the cmd/web directory.

- `pkg`: contains the ancillary non-application-specific code used in the project. We’ll use it to hold potentially reusable code like validation helpers and the SQL database models for the project.

- `ui`: contains the user-interface assets used by the web application. Specifically, the ui/html directory will contain HTML templates, and the ui/static directory will contain static files (like CSS and images).

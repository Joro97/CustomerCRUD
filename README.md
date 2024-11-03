# Description:
This is basic CRUD application for a Customer entity that should look
something like this: <br>
Entity: Customer<br>
Attribute Type Constraints Notes<br>
Id UUID PK<br>
First Name String<br>
Middle Name String Null is acceptable<br>
Last Name String<br>
Email Address String Unique<br>
Phone Number String Can be composite<br>

# Prerequisites:
Make sure you have created a .env file in the root directory, where this readme resides. If you decide
to not use local db (recommended) you should also register for Neon managed Postgres DBaaS:
1. If you do not have a Neon account, [click here](https://console.neon.tech/realms/prod-realm/protocol/openid-connect/auth?client_id=neon-console&redirect_uri=https%3A%2F%2Fconsole.neon.tech%2Fauth%2Fkeycloak%2Fcallback&response_type=code&scope=openid+profile+email&state=6tSQVbgMQ2Al1q6GRWXHqA%3D%3D%2C%2C%2C) to sign up for an account.
2. Log in to your Neon account.
3. On the Console page, click Create project.
4. On the Create project page, the highest postgres version is selected by default. Name your project and db. Select the desired Region. I ran the project with an instance on AWS Europe (Frankfurt) and one dev branch (both covered by the free plan)
5. Click Create project to create a Neon project with a database.
6. (Optional) Create a dev branch from the Neon console, if you wish your tests to run against a copy of your main db (recommended)<br>

7. The .env file must contain 3 variables:<br>
   1. DATABASE_URL - the connection string that can be obtained from your Neon console
   2. TEST_DATABASE_URL - the connection string for your dev/testing branch (copy the above if you didnt create one)
   3. LOCAL_DB - set to 'true' or 'false', depending on your desire for running against a local in memory db
## Important:
The application is setup to read the .env file and load its contents as env variables in the application. The file _MUST_ be present for the application to work properly!

# Usage:
There is a Makefile that has simple commands for user convenience. Some of them include:
1. Make unit - will run the unit tests of the application, due to time limitations app is not 100% covered on all files
2. Make integration - will run all the tests of the application and provide a basic coverage report
3. Make build-image - builds a docker image for the server. `docker run customer-service` will start the service inside the container
4. Make deploy will create a local kind cluster and install a helm chart with the application into it. Make sure to run `kubectl port-forward svc/customer-service 8080:8080` afterwards and you can call your app from the kind cluster like `http://localhost:8080/customers`
5. The server can also be started manually by running `go run cmd/main.go`. A tool like Postman or cURL can be used to manually validate the endpoints, examples:
`   curl --location 'localhost:8080/customers' \
   --header 'Content-Type: application/json' \
   --data-raw '{
   "first_name": "Joro",
   "middle_name": "Kripto",
   "last_name": "Imperatora",
   "email": "mailera@example.com",
   "phone_number": "1324"
   }'` to add a customer and `curl --location 'http://localhost:8080/customers'` to get the list of customers.

# Improvements:
For Observability we can have and architecture that would leverage fluent-bit (can be installed into our cluster easily) to forward
the pod logs (we should update them to structured) to something like ELK or Splunk. For pods health metrics we could leverage
Prometheus and Grafana. All of those have very good open source operators that can be leveraged. CD pipeline needs to also be implemented, it could look
something like: Push a new helm chart into a repository on each successful commit to master, then have the CD deploy this image into dev and with manual approval to prod.
Make lint shows quite some stuff to be refactored. The architecture generally could be improved with k8s secrets, a bit refactoring of the way to switch dbs, etc.
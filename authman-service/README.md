## Authman Api

## Description

This repository contains the API logic for an Authentication and Authorization service. 

## Deployment
A live deployment of this application can be found @ https://f-c3nm.onrender.com

## Documentation
The documentation of this application can be found @:
<li>Open Api @ https://authman-api.onrender.com/api/docs</li>
<li>Postman @ https://github.com/TosinJs/authman/blob/master/Authman.postman_collection.json </li>

## Run the Application Locally

```bash
# Clone the repository
$ git clone https://github.com/TosinJs/food-api.git

# Install dependencies
$ npm install

# configuration 
# Create .env file in the root folder
$ touch .env

# populate the .env file with your files
$ MONGO_URI = "your mongo URI"
$ REDIS_URL = "Your Redis URL"
$ SMTP_HOST = "mail host"
$ SMTP_USER = "mail user"
$ SMTP_PASS = "mail password"
$ SMTP_PORT = "mail port"
$ SMTP_SERVICE = "mail service"

# watch mode
$ npm run start:dev

# production mode
$ npm run start:prod
```

## The Application

### Funtionality
Users can 
<li>Create account</li>
<li>Login</li>
<li>Email Verification</li>
<li>Password Reset</li>
<li>ROLE Based Access</li>

### Details
<p>On signup users are sent a verification Email, without emailverification, their token roles are limited</p>
<p>Verified users have access to three tokens</p>
	- ID Token
	- Auth Token
	- Refresh Token
<p>The ID Token is used to identify the users, all users have access to this token and it has a life span of 24h</p>
<p>The Auth Token is used for Role Based access. These tokens grant users special access to specific resources. The life span of the auth tokens are about 20 minutes</p>
<p>The Refresh Tokens are used to generate a new set of tokens when the user submits a token that is valid but expired. Refresh tokens have a life span of 3 months. The refresh tokens are stored in a redis database for validation purposes and to enable token revocation.



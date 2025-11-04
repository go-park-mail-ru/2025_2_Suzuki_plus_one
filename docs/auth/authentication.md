# Authentication flow

Minimal authentication flow using JWT access/refresh tokens:

> Note that JWT tokens contain userID and expiration time claims.

## Frontend

If JS has Access token, then User is considered authenticated.

- Access token keeps in JS memory (to obtain hit `/auth/refresh`)
- Refresh token keeps in HttpOnly Secure SameSite=strict cookie (to obtain hit `/auth/signin`)

## Backend

- Access token stored in Redis
- Refresh token stored in Postgres

## Sequence diagram

```mermaid
sequenceDiagram
    note right of Client: Client signs in with credentials

    Client->>Server: POST /auth/signin
    Server<<->>Postgres: Validate credentials
    alt bad credentials
        Server->>Client: 403 Forbidden
    else valid credentials
        Server->>Client: 200 OK + JWT ACCESS Token + JWT REFRESH Token
    end

    note right of Client: Client accesses protected resource using access token

    Client->>Server: GET /protected/resource + Authorization: Bearer <JWT ACCESS Token>
    Server<<->>Redis: Validate access token
    alt access token expired
        Server->>Client: 401 Unauthorized
    else access token valid
        Server->>Client: 200 OK + Protected Resource
    end

    note right of Client: Client refreshes access token when expired
    Client->>Server: GET /auth/refresh + Authorization: Bearer <JWT REFRESH Token>
    Server<<->>Postgres: Validate refresh token
    alt refresh token invalid
        Server->>Client: 403 Forbidden
    else refresh token valid
        Server<<->>Redis: Invalidate old tokens, put new one
        Server->>Client: 200 OK + New JWT ACCESS Token + New JWT REFRESH Token
    end
```

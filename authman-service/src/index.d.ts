interface responseObj {
    accessToken: string,
    refreshToken: string,
    idToken: string,
    tokenType: "Bearer"
}

interface tokenPayload {
    id: string,
    username: string,
    email?: string,
    isVerified?: boolean
}
# Get Google ID Token - Step by Step Guide

## Current Step: Google Cloud Console is Open ✅

You're at the Google Cloud Console. Now follow these steps to get your **Client ID**:

---

## Step 1: Go to APIs & Services > Credentials

1. In left sidebar, click **APIs & Services**
2. Click **Credentials**
3. You should see a page showing your OAuth credentials

---

## Step 2: Get Your Client ID

Look for **OAuth 2.0 Client IDs** section:
- Find the entry with type **Web application**
- Copy the **Client ID** (looks like: `123456789-abcdefghijk.apps.googleusercontent.com`)

---

## Step 3: Use OAuth 2.0 Playground (Easiest Way to Test)

1. Open [OAuth 2.0 Playground](https://developers.google.com/oauthplayground/)
2. Click **gear icon** (⚙️) in top right
3. Check box: **"Use your own OAuth credentials"**
4. Paste your **Client ID** and **Client Secret** from Step 2
5. Click **Close**

---

## Step 4: Authorize

Back on playground main page:
1. On **left side**, expand **Google OAuth2 API v2**
2. Select these scopes:
   - ✅ `https://www.googleapis.com/auth/userinfo.email`
   - ✅ `https://www.googleapis.com/auth/userinfo.profile`
3. Click **AUTHORIZE APIS** button
4. Select your Google account
5. Click **Allow**

---

## Step 5: Exchange Code for Tokens

1. Click **EXCHANGE AUTHORIZATION CODE FOR TOKENS** button
2. In the right panel under **Step 3**, you'll see the response with:
   ```json
   {
     "access_token": "...",
     "expires_in": 3599,
     "refresh_token": "...",
     "scope": "...",
     "token_type": "Bearer",
     "id_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20iLCJhenAiOiI..."
   }
   ```

3. **Copy the `id_token` value** (it's a long JWT string)

---

## Step 6: Test with Your API

Now use this `id_token` to test your login endpoint:

```bash
curl -X POST "http://localhost:8080/api/user/login" \
  -H "Content-Type: application/json" \
  -d '{
    "id_token": "PASTE_YOUR_ID_TOKEN_HERE"
  }'
```

Replace `PASTE_YOUR_ID_TOKEN_HERE` with the actual token from Step 5.

**Expected Response:**
```json
{
  "code": 200,
  "message": "Login successfully",
  "data": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InlvdXJlbWFpbEBnbWFpbC5jb20iLCJleHAiOjE3MzY0MzI4MDB9..."
}
```

---

## Video Alternative (If Steps Above Don't Work)

Follow this video: [Google OAuth 2.0 Setup](https://www.youtube.com/watch?v=496zpmQT0B8)

---

## Quick Reference Table

| Step | Action | Result |
|------|--------|--------|
| 1 | Open Google Cloud Console | ✅ Done |
| 2 | Go to Credentials | Copy **Client ID** |
| 3 | Open OAuth Playground | Ready to authorize |
| 4 | Authorize with scopes | Get authorization code |
| 5 | Exchange for tokens | Get **id_token** |
| 6 | Test with curl | Validate your API works |

---

## Troubleshooting

### "Redirect URI mismatch"
- Make sure in Google Cloud Console under **Credentials > OAuth Client ID**:
  - Add `https://developers.google.com/oauthplayground` to **Authorized redirect URIs**
  - Click **Save**

### "Invalid client_id"
- Copy the exact Client ID from Google Cloud Console
- Remove any spaces or quotes

### "Token is invalid"
- Make sure you copied the full `id_token` value
- It should start with `eyJ...`

---

## File Structure for Reference

Once you get the token, here's where you'll use it in your project:

```
docs/
├── auth_api_examples.md          ← Test curl commands here
├── GOOGLE_OAUTH_SETUP.md         ← Full setup guide
└── GOOGLE_OAUTH_GET_TOKEN.md     ← You are here
```

---

## Next Steps After Getting Token

1. ✅ Test login with `curl` command
2. ✅ Copy the JWT token from response
3. ✅ Use that token for authenticated requests:
   ```bash
   curl -X GET "http://localhost:8080/api/menu" \
     -H "Authorization: Bearer YOUR_JWT_TOKEN"
   ```

---

## Support Files

- **Main setup**: `GOOGLE_OAUTH_SETUP.md`
- **Test examples**: `auth_api_examples.md`
- **Database migrations**: `migrations/005_create_authentication.sql`

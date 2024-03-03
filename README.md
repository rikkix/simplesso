# SimpleSSO
A simple SSO server built for nginx.

## Introduction
There has long been a dilemma about authorizing users on self-hosted websites. A common practice is using a simple username-password model. However, it is always an annoying way to remember all these usernames and passwords, even with the help of password managers. The situation worsens when you have hosted several services, as you must log into every service before using them. It can be even harder to deal with when you want to access these services on a public computer.

This project can solve these problems by setting up an easy-to-deploy authorization layer for each self-hosted service (currently tested with Nginx) and using a request-confirmation mechanism instead of passwords. 

More specifically, this project achieves the following effects:
- When accessing one of your self-hosted services `code.example.com`, it first checks whether you have logged in. If you haven't, the webpage will be redirected to your authorization center `auth.example.com` for further verification.

- On `auth.example.com`, you will be asked for the pre-setted username. A confirmation message will be sent to your Telegram account if you have entered the correct username.

- Click the `Confirm` button on Telegram, and enter the verification code on `auth.example.com`, then you have successfully logged in. The webpage will automatically redirect to the initial page you want to visit.

If you want to visit other self-hosted services, say `file.example.com`, you will be redirected to `auth.example.com` first. Since you have been verified, you will be redirected back to `file.example.com`, and no actions are needed for this time!
# s3-cloudfront-purger

Automatically invalidate your s3 static website cloudfront distribution on updates.

This lambda is triggered by s3 event, just create an event on `putItem` on your `index.html` to trigger this lambda.

What it does:

1. Take the s3 bucket name from the event context
2. Locate the cloudfront distribution that its used
3. Purge the cloudfront distribution

Setup:

1. Create an IAM role with Basic Lambda Execution, Read your bucket(s) and List and Purge Cloudfront distibutions
2. Create a new lambda with the go1.x runtime and the name `s3-cloudfront-purger`
3. Set handler name: `app`
4. Create repository secrets with `updateFunctionCode` IAM permission.
    1. Secret `AWS_ACCESS_KEY_ID`
    2. Secret `AWS_SECRET_ACCESS_KEY`

Deployed by Github Actions on master branch only

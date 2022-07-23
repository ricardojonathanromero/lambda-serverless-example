# api-gateway-example
This artifact is used to expose a lambda function via API Gateway

## Deploy
To deploy this lambda function is necessary configure your aws credentials at first time using the next command:
- `aws configure --profile ${profile_name_to_assign}`

The above code is necessary to create the file credentials located in `$HOME/.aws` folder. Once do that, you could deploy your lambda into your aws account.

To deploy your code, you must execute this command: `serverless deploy`
This command will deploy your lambda into the account specified by default, but, if you have multiple accounts, or if you have not specified one by default, you could set up which use configuring the env variable `AWS_DEFAULT_PROFILE` and passing your account name as value, for example:
- for Windows execute in your powershell `$AWS_DEFAULT_PROFILE=profile_name`
- for Mac/Unix systems `export AWS_DEFAULT_PROFILE=profile_name`

---
However, if you do not want to configure your aws account as env variable, another option is pass the aws account name to the serverless command, e.g.
`serverless deploy --aws-profile profile_name` This command will deploy your lambdas using the aws account specified.

---
Finally, if you have been deployed your lambdas, but you are modified just one, and you would to deploy this change, you can deploy just this change using this command: `serverless deploy --aws-profile profile_name function function_name_to_deploy`
This command will deploy just the specified function.

## APM
I am using DataDog as APM provider, if you need deep and understand how to configure it, please check this [link](https://docs.datadoghq.com/serverless/installation/go/?tab=serverlessframework)

To configure DD as APM, first we need to install the DD plugin. To do that, we exec this command `serverless plugin install --name serverless-plugin-datadog`

Once we have the DD plugin installed, we need to custom our serverless.yaml file adding the next lines:
```yaml
...
custom:
  datadog:
    site: <DATADOG_SITE> # DATADOG_SITE is your DataDog domain
    apiKeySecretArn: <DATADOG_API_KEY_SECRET_ARN> # DATADOG_API_KEY_SECRET_ARN is your API KEY
...
```

I've configured these secrets using a JSON file as you can see within my serverless.yml file, however, I recommend you that these credentials be stored in a credential storage service such as Vault by Hashicorp, or SSM by AWS.

Now, if you are asking about how I'm reading the values from a JSON file, let me refer to this [link](https://www.serverless.com/framework/docs/providers/aws/guide/variables#reference-variables-in-javascript-files) where you can look how to reference values from JSON files.

And, finally, how to refer values from SSM or Vault? You can check this [link](https://www.serverless.com/framework/docs/providers/aws/guide/variables#reference-variables-using-the-ssm-parameter-store)
⚠️ To use this, you need to configure an aws account profile.

Once you configure your datadog credentials, it is time to implement the profiler within your code. To do this, follow this [example](https://docs.datadoghq.com/serverless/installation/go/?tab=serverlessframework#update-your-lambda-function-code) by Datadog where they show how to do that.

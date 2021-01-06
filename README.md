# netlify-teams-webhook

This webhook service can be used to send netlify deployment notifications to
Microsoft Teams. Workflow is

```text
GitHub -> Netlify -> netlify-teams-webhook -> MS Teams
```

When everything is set up, one gets a notification to MS Teams channel, when new
site is published on Netlify, including some metadata like publish time, build
time and most importantly, unique url to the app.

## Installing

Clone repository and run `go build`. Start with command

```bash
PORT=xxxx
TEAMS_WEBHOOK_URL=https://outlook.office.com/webhook/yyyy
./netlify-teams-webhook
```

Here, `TEAMS_WEBHOOK_URL` is achieved by adding incoming webhook to Teams
channel.

## Installing with Docker

```bash
docker run -d --name=netlify-teams-webhook \
    --env TEAMS_WEBHOOK_URL=$TEAMS_WEBHOOK_URL \
    --env PORT=$PORT -p $PORT:$PORT ahojukka5/netlify-teams-webhook
```

## Testing

Test with `curl`, i.e.

```bash
curl -X PORT -H 'Content-Type: application/json' -H 'X-Netlify-Event: deploy_created' --data=@payload.json $HOST:$PORT/deploy_created
```

`payload.json` is the payload that comes from Netlify during triggering the
webhook.

## Other

In principle, program could be extended quite straightforwardly to send "cards"
to MS Teams from other CI/CD pipelines also. Basically one needs only to know
the structure of the payload send by triggering webhook.

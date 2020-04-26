NAME=discord_rainbow_bot

docker stop $NAME
docker rm $NAME

docker run \
-d \
-e DISCORD_TOKEN=$DISCORD_TOKEN \
-e INVITE_URL=$INVITE_URL \
-e DELAY_MS=$DELAY_MS \
--name $NAME \
--restart=on-failure \
"$CI_REGISTRY_IMAGE:$CI_COMMIT_SHORT_SHA"

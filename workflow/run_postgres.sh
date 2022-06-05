CONTAINER_NAME=postgres-rga
PORT=5432
VOLUME_PATH=$HOME/local_codes/remix-golang-authentication/postgres_db
IMAGE_NAME=postgres:latest
POSTGRES_PASSWORD=adminpass

docker run -d \
		--name $CONTAINER_NAME \
		--restart always \
		-p $PORT:$PORT \
		-v $VOLUME_PATH:/var/lib/postgresql/data \
		-e "PORT=$PORT" \
		-e "VOLUME_PATH=$VOLUME_PATH" \
		-e "CONTAINER_NAME=$CONTAINER_NAME" \
		-e "POSTGRES_PASSWORD=$POSTGRES_PASSWORD" \
		$IMAGE_NAME

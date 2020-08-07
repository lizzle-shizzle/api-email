IMAGE='api-email'

echo '>>> Building command...'
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .
echo '...Done'

echo '>>> Building image...'
docker image build -t $IMAGE .

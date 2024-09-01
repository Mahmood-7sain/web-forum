#!bin/bash
echo "Enter the name of the image:"
read image
echo "Enter the name of the container:"
read container
echo "Enter port to listen on:"
read port
echo ""
echo "Initiating build of $image image...."
docker build -f dockerfile -t $image .
echo ""
echo "Image: $image finished building."
echo "Running $container container using $image image on port: $port"
docker container run -p $port:8080 --detach --name $container $image
echo ""
echo "Container running on http://localhost:$port"
echo "To stop container run: docker kill $container"
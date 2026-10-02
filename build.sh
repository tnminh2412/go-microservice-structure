for file_path in $(find "." -type f -name "main.go" ); do
  dir_path=$(dirname $file_path)
  echo $dir_path
  cd $dir_path && go build -o ./bin main.go && cd $GOPATH/src/go-microservice-structure

done;
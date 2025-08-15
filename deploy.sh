#!/bin/bash

LAMBDAS_DIR="$(pwd)/lambdas"
IFS=',' read -ra LAMBDA_ARNS_ARRAY <<< "$LAMBDA_ARNS"
echo "Starting deployment of ${#LAMBDA_ARNS_ARRAY[@]} Lambda functions..."

SUCCESSFUL=()
FAILED=()

for arn in "${LAMBDA_ARNS_ARRAY[@]}"; do
  function_name=$(echo $arn | awk -F: '{print $7}')
  lambda_dir="${LAMBDAS_DIR}/${function_name}"
  
  echo "Deploying $function_name..."
  cd "$lambda_dir"

  echo "Building $function_name..."
  GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap main.go
  
  echo "Packaging $function_name..."
  zip -q function.zip bootstrap main.go
  
  echo "Deploying $function_name to AWS..."
  aws lambda update-function-code --function-name "$function_name" --zip-file fileb://function.zip --region "$AWS_REGION"
  
  if [ $? -eq 0 ]; then
    SUCCESSFUL+=("$function_name")
  else
    FAILED+=("$function_name")
  fi
  
  rm -f bootstrap function.zip
  cd - > /dev/null
done

echo "Successfully deployed: ${SUCCESSFUL[*]}"
if [ ${#FAILED[@]} -gt 0 ]; then
  echo "Failed to deploy: ${FAILED[*]}"
fi
echo "Deployment completed."
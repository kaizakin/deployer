pipeline {
    agent {
        docker {
            image 'golang:1.26-alpine'
            args '--network=devops-secure-mesh'
        }
    }

    environment {
        REGISTRY_TARGET  = 'localhost:5000'
        APPLICATION_NAME = 'go-service'
        APPLICATION_TAG  = "${APPLICATION_NAME}-${env.BUILD_NUMBER}"
    }

    stages {
        stage('Automated Source Verification') {
            steps {
                echo 'Executing static syntax verification and testing sweeps...'
                sh '''
                    apk add --no-cache git
                    go fmt ./...
                    go vet ./...
                    go test -v -coverprofile=coverage.out ./...
                '''
            }
        }

        stage('Rootless Containerization (Kaniko Sandbox)') {
            agent {
                docker {
                    image 'gcr.io/kaniko-project/executor:debug'
                    args '--entrypoint="" --network=devops-secure-mesh'
                }
            }
            steps {
                echo 'Assembling immutable container layers using zero-privilege space...'
                sh """
                    /kaniko/executor \
                        --context=dir://\${WORKSPACE} \
                        --dockerfile=\${WORKSPACE}/Dockerfile \
                        --destination=\${REGISTRY_TARGET}/\${APPLICATION_TAG} \
                        --destination=\${REGISTRY_TARGET}/\${APPLICATION_NAME}:latest \
                        --insecure \
                        --skip-tls-verify
                """
            }
        }

        stage('Zero-Downtime Blue-Green Deployment') {
            agent any
            steps {
                echo 'Orchestrating container lifecycle transition gates...'
                sh """
                    # Fetch the newly validated layer image from the local registry
                    docker pull \${REGISTRY_TARGET}/\${APPLICATION_TAG}

                    # Launch a secondary instance container under a distinct active target ID reference
                    docker run -d \
                        --name \${APPLICATION_NAME}-green \
                        --network devops-secure-mesh \
                        -e APP_PORT=8081 \
                        \${REGISTRY_TARGET}/\${APPLICATION_TAG}

                    # Warmup check: Confirm the green instance is healthy before cutting over
                    sleep 5
                    docker run --rm --network devops-secure-mesh curlimages/curl:latest \
                        -f http://\${APPLICATION_NAME}-green:8081/health

                    # Gracefully stop the older processing node container and shift aliases
                    docker stop \${APPLICATION_NAME}-blue || true
                    docker rm \${APPLICATION_NAME}-blue || true
                    docker rename \${APPLICATION_NAME}-green \${APPLICATION_NAME}-blue
                """
            }
        }
    }
}

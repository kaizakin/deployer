pipeline {
    agent any

    triggers {
        pollSCM('H/5 * * * *')
    }

    environment {
        DOCKERHUB_USER   = 'voidmoon' 
        APPLICATION_NAME = 'go-service'
        APPLICATION_TAG  = "${env.BUILD_NUMBER}"
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
                // IMPORTANT: Create a 'Username with password' credential in Jenkins named 'dockerhub-credentials'
                withCredentials([usernamePassword(credentialsId: 'dockerhub-credentials', passwordVariable: 'DOCKER_PASSWORD', usernameVariable: 'DOCKER_USERNAME')]) {
                    sh """
                        # Create Kaniko Docker config for Docker Hub authentication
                        mkdir -p /kaniko/.docker
                        AUTH=\$(echo -n "\${DOCKER_USERNAME}:\${DOCKER_PASSWORD}" | base64 | tr -d '\\n')
                        echo '{"auths": {"https://index.docker.io/v1/": {"auth": "'\$AUTH'"}}}' > /kaniko/.docker/config.json

                        /kaniko/executor \
                            --context=dir://\${WORKSPACE} \
                            --dockerfile=\${WORKSPACE}/Dockerfile \
                            --destination=\${DOCKERHUB_USER}/\${APPLICATION_NAME}:\${APPLICATION_TAG} \
                            --destination=\${DOCKERHUB_USER}/\${APPLICATION_NAME}:latest
                    """
                }
            }
        }

        stage('Zero-Downtime Blue-Green Deployment') {
            agent any
            steps {
                echo 'Orchestrating container lifecycle transition gates...'
                sh """
                    # Fetch the newly validated layer image from Docker Hub
                    docker pull \${DOCKERHUB_USER}/\${APPLICATION_NAME}:\${APPLICATION_TAG}

                    # Launch a secondary instance container under a distinct active target ID reference
                    docker run -d \
                        --name \${APPLICATION_NAME}-green \
                        --network devops-secure-mesh \
                        -e APP_PORT=8081 \
                        \${DOCKERHUB_USER}/\${APPLICATION_NAME}:\${APPLICATION_TAG}

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

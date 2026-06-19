pipeline {
    agent any

    environment {
        GO_VERSION = '1.22'
        APP_NAME = 'bytestutor-backend'
        DEPLOY_HOST = '64.110.115.33'
        DEPLOY_USER = 'ubuntu'
        SONAR_HOST_URL = 'http://161.118.199.97:9900'
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Unit Tests & Coverage') {
            steps {
                dir('backend') {
                    sh 'go test -v -coverprofile=coverage.out ./...'
                }
            }
        }

        stage('SonarQube Analysis') {
            steps {
                dir('backend') {
                    // Requires SonarQube Scanner configured in Jenkins global tools
                    withSonarQubeEnv('SonarQubeServer') {
                        sh 'sonar-scanner'
                    }
                }
            }
        }

        stage('Build') {
            steps {
                dir('backend') {
                    sh 'go build -o bin/bytestutor-backend main.go'
                }
            }
        }

        stage('Deploy') {
            steps {
                // Use Jenkins SSH credentials to copy and restart services
                sshagent(credentials: ['apps-vm-ssh-key']) {
                    // Copy binary to Apps VM
                    sh "scp -o StrictHostKeyChecking=no backend/bin/bytestutor-backend ${DEPLOY_USER}@${DEPLOY_HOST}:/home/ubuntu/${APP_NAME}"
                    
                    // Restart service on target VM
                    sh "ssh -o StrictHostKeyChecking=no ${DEPLOY_USER}@${DEPLOY_HOST} 'sudo systemctl restart bytestutor || docker restart bytestutor-backend || true'"
                }
            }
        }
    }

    post {
        always {
            cleanWs()
        }
        success {
            echo "CI/CD Pipeline executed successfully!"
        }
        failure {
            echo "Pipeline failed. Check build logs."
        }
    }
}

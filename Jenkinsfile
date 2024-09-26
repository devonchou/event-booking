pipeline {
    agent any

    tools { go '1.22.5' }

    environment {
        registry = "devonchou/cicd-event-booking"
        registryCredential = "dockerhub"
    }

    stages {
        stage('Build') {
            steps {
                sh 'go build -o application .'
            }
        }

        stage('Test') {
            steps {
                sh 'go test -coverprofile=coverage.out ./...'
            }
        }

        stage('Sonar Analysis') {
            environment {
                scannerHome = tool 'sonar4.7'
            }

            steps {
                withSonarQubeEnv('sonar') {
                    sh '''${scannerHome}/bin/sonar-scanner \
                    -Dsonar.projectKey=event-booking-api \
                    -Dsonar.projectName=event-booking-api \
                    -Dsonar.projectVersion=1.0 \
                    -Dsonar.sources=app/ \
                    -Dsonar.exclusions=**/*_test.go \
                    -Dsonar.tests=test/ \
                    -Dsonar.test.inclusions=**/*_test.go \
                    -Dsonar.go.coverage.reportPaths=coverage.out'''
                }
            }
        }

        stage('Quality Gate') {
            steps {
                timeout(time: 1, unit: 'HOURS') {
                    waitForQualityGate abortPipeline: true
                }
            }
        }

        stage('Build Docker Image') {
            steps {
                script {
                    dockerImage = docker.build("${registry}:v${BUILD_NUMBER}")
                }
            }
        }

        stage('Upload Image') {
            steps {
                script {
                    docker.withRegistry('', registryCredential) {
                        dockerImage.push("v$BUILD_NUMBER")
                        dockerImage.push('latest')
                    }
                }
            }
        }

        stage('Remove Unused Docker Image') {
            steps {
                sh "docker rmi ${registry}:v${BUILD_NUMBER}"
            }
        }

        stage('Kubernetes Deploy') {
            agent {label 'KOPS'}

            steps {
                sh "helm upgrade --install --force event-booking-stack helm/event-booking-charts --set apimage=${registry}:v${BUILD_NUMBER} --namespace prod"
            }
        }
    }
}

def COLOR_MAP = [
    'SUCCESS': 'good',
    'FAILURE': 'danger',
]

pipeline {
    agent any

    tools { go '1.22.5' }

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

        stage("Upload Artifact") {
            steps {
                nexusArtifactUploader(
                    nexusVersion: 'nexus3',
                    protocol: 'http',
                    nexusUrl: 'nexus3:8081',
                    groupId: 'QA',
                    version: "${env.BUILD_ID}-${env.BUILD_TIMESTAMP}",
                    repository: 'event-booking-repo',
                    credentialsId: 'nexuslogin',
                    artifacts: [
                        [artifactId: 'event-booking-api',
                        classifier: '',
                        file: './application',
                        type: 'bin']
                    ]
                )
            }
        }
    }

    post {
        always {
            echo 'Slack Notifications.'
            slackSend channel: '#jenkinscicd',
                color: COLOR_MAP[currentBuild.currentResult],
                message: "*${currentBuild.currentResult}:* Job ${env.JOB_NAME} build ${env.BUILD_NUMBER} \n More info at: ${env.BUILD_URL}"
        }
    }
}

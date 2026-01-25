pipeline {
    agent any
    options {
        disableConcurrentBuilds()
    }
    stages {
        stage('Checkout'){
            steps {
                checkout scm
            }
        }
        stage('Test') {
            when {
                branch 'master'
                not { changelog '^skip-tests.*' }
            }
            steps {
                sh """
                    docker run --rm -v \$PWD:/app -w /app golang:1-alpine go test ./...
                    """
            }
        }
        stage('Prep buildx') {
            when { branch 'master' }
            steps {
                script {
                    env.BUILDX_BUILDER = getBuildxBuilder();
                }
            }
        }
        stage('Dockerhub login') {
            when { branch 'master' }
            steps {
                withCredentials([usernamePassword(credentialsId: 'dockerhub', usernameVariable: 'DOCKERHUB_CREDENTIALS_USR', passwordVariable: 'DOCKERHUB_CREDENTIALS_PSW')]) {
                    sh 'docker login -u $DOCKERHUB_CREDENTIALS_USR -p "$DOCKERHUB_CREDENTIALS_PSW"'
                }
            }
        }
        stage('Build Docker Image') {
            when { branch 'master' }
            steps {
                sh """
                    docker buildx build --pull --builder \$BUILDX_BUILDER --platform linux/arm64,linux/amd64 -t nbr23/numerology-go:latest -t nbr23/numerology-go:`git rev-parse --short HEAD` --push .
                    """
            }
        }
        stage('Sync github repos') {
            when { branch 'master' }
            steps {
                syncRemoteBranch('git@github.com:nbr23/numerology-go.git', 'master')
            }
        }
    }
    post {
        always {
            sh 'docker buildx stop $BUILDX_BUILDER || true'
            sh 'docker buildx rm $BUILDX_BUILDER || true'
        }
    }
}

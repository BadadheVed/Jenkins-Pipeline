pipeline {
    agent any

    environment {
        EC2_USER = "ubuntu"
        EC2_HOST = "13.201.168.68"
        EC2_DIR  = "/home/ubuntu/myapp"
        SSH_KEY_ID = "PRIVATE_KEY"
        GO_VERSION = "1.22.1"
    }

    stages {

        stage('Checkout Code') {
            steps {
                checkout scm
            }
        }

        stage('Pre-Install + Prepare Directory on EC2') {
    steps {
        sshagent([SSH_KEY_ID]) {
            sh """
                ssh -o StrictHostKeyChecking=no ${EC2_USER}@${EC2_HOST} '
                    set -e

                    echo "----- Creating project directory if missing -----"
                    mkdir -p ${EC2_DIR} || true

                    echo "----- Installing Go if not installed -----"
                    if ! command -v go >/dev/null 2>&1; then
                        echo "Installing Go..."
                        wget -q https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz
                        sudo rm -rf /usr/local/go
                        sudo tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz
                        echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.profile
                    else
                        echo "Go already installed"
                    fi

                    if ! command -v pm2 >/dev/null 2>&1; then
                        sudo npm install -g pm2 || true
                    fi

                    echo "----- Setup PM2 startup (ignore if already set) -----"
                    sudo -n env PATH=$PATH:/usr/bin pm2 startup systemd -u ubuntu --hp /home/ubuntu || true

                '
            """
        }
    }
}

        stage('Upload Code to EC2') {
            steps {
                sshagent([SSH_KEY_ID]) {
                    sh """
                        echo "----- Uploading code to EC2 -----"
                        rsync -avz --delete -e "ssh -o StrictHostKeyChecking=no" ./ ${EC2_USER}@${EC2_HOST}:${EC2_DIR}
                    """
                }
            }
        }

        stage('Build + Restart App on EC2') {
    steps {
        sshagent([SSH_KEY_ID]) {
            sh """
                ssh -o StrictHostKeyChecking=no ${EC2_USER}@${EC2_HOST} '
                    set -e

                    echo "----- Ensuring Go present in PATH -----"
                    export PATH=\$PATH:/usr/local/go/bin

                    echo "----- cd into project -----"
                    cd ${EC2_DIR}

                    echo "----- Running go mod tidy -----"
                    go mod tidy

                    echo "----- Building Go binary -----"
                    go build -o app main.go

                    echo "----- Restarting PM2 -----"
                    pm2 stop go-app || true
                    pm2 start app --name go-app
                    pm2 save
                '
            """
        }
    }
}
    }

    post {
        success { echo "✔ Deployment Successful!" }
        failure { echo "❌ Deployment Failed!" }
    }
}

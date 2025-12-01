pipeline{
    agent any

    environment {
        EC2_USER = "ubuntu"
        EC2_HOST = "13.201.168.68" 
        SSH_KEY_ID = "PRIVATE_KEY"
        GO_VERSION = "1.22.1"   
    }

    stages{
        stage('Checkout Code') {
            steps {
                checkout scm
            }
        }


        stage('Build on Jenkins') {
            steps {
                sh """
                echo 'Building Go app on Jenkins...'
                go mod tidy
                go build -o app main.go
                """
            }
        }

        stage('Pre-Install Dependencies on EC2') {
            steps {
                sshagent([SSH_KEY_ID]) {
                    sh """
                        ssh -o StrictHostKeyChecking=no ${EC2_USER}@${EC2_HOST} '
                            
                            if ! command -v go >/dev/null 2>&1; then
                                echo "Go not found — installing Go ${GO_VERSION}..."
                                wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz
                                sudo rm -rf /usr/local/go
                                sudo tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz
                                echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.profile
                                source ~/.profile
                            else
                                echo "Go already installed."
                            fi

                           
                            if ! command -v pm2 >/dev/null 2>&1; then
                                echo "PM2 not found — installing PM2..."
                                sudo npm install -g pm2
                            else
                                echo "PM2 already installed."
                            fi

                            # Ensure PM2 restarts on reboot
                            pm2 startup systemd -u ubuntu --hp /home/ubuntu
                        '
                    """
                }
            }
            
        }

    }

    post {
        success { echo "Deployment Successful!" }
        failure { echo "Deployment Failed!" }
    }
}
pipeline {
    agent any
    parameters {
        // Only the main branch gets signed by default, as the number of code signings is limited
        booleanParam(name: 'SIGN_WINDOWS', defaultValue: false, description: 'Sign the Windows MSI packages also for branches other than main')
    }
    environment {
        CIBUILD = "1"
        //ADVINST = "\"C:\\Program Files (x86)\\Caphyon\\Advanced Installer\\bin\\x86\\advinst.exe\""
        ADVINST = "\"C:\\Program Files (x86)\\Caphyon\\Advanced Installer 23.9\\bin\\x86\\AdvancedInstaller.com\""
        //ADVINST = "\"C:\\Program Files (x86)\\Caphyon\\Advanced Installer 18.8.1\\bin\\x86\\AdvancedInstaller.com\""
        // AdvancedInstaller.com is the command line interface of Advanced Installer
        // advinst.exe is for the GUI version
        // https://www.advancedinstaller.com/user-guide/command-line.html
    }
    stages {
        stage('Test') {
            environment {
                CGO_ENABLED = '0'
            }
            parallel {
                stage('windows') {
                    agent {
                        label 'windows'
                    }
                    environment {
                        GOOS = 'windows'
                        BINNAME = 'openitcockpit-agent.exe'
                    }
                    stages {
                        stage('amd64') {
                            environment {
                                GOARCH = 'amd64'
                            }
                            steps {
                                test_windows()
                            }
                        }
                        stage('arm64') {
                            environment {
                                GOARCH = 'arm64'
                            }
                            steps {
                                test_windows()
                            }
                        }
                        stage('386') {
                            environment {
                                GOARCH = '386'
                            }
                            steps {
                                test_windows()
                            }
                        }
                    }
                }
                stage('linux') {
                    environment {
                        GOOS = 'linux'
                        BINNAME = 'openitcockpit-agent'
                    }
                    stages {
                        stage('amd64') {
                            agent {
                                docker {
                                    image 'golang:bookworm'
                                    args "-u root --privileged -v agentgocache:/go"
                                    label 'linux'
                                }
                            }
                            environment {
                                GOARCH = 'amd64'
                            }
                            steps {
                                test()
                            }
                        }
                        stage('386') {
                            agent {
                                docker {
                                    image 'golang:bookworm'
                                    args "-u root --privileged -v agentgocache:/go"
                                    label 'linux'
                                }
                            }
                            environment {
                                GOARCH = '386'
                            }
                            steps {
                                test()
                            }
                        }
                        stage('arm64') {
                            agent {
                                docker {
                                    image 'golang:bookworm'
                                    args "-u root --privileged -v agentgocache:/go"
                                    label 'linux-arm64'
                                }
                            }
                            environment {
                                GOARCH = 'arm64'
                            }
                            steps {
                                test()
                            }
                        }
                    }
                }
                stage('darwin') {

                    environment {
                        GOOS = 'darwin'
                        BINNAME = 'openitcockpit-agent'
                        CGO_ENABLED = '1'
                    }
                    stages {
                        stage('amd64') {
                            agent {
                                label 'macos-intel64'
                            }
                            environment {
                                GOARCH = 'amd64'
                            }
                            steps {
                                test()
                            }
                        }
                        stage('arm64') {
                            agent {
                                label 'macos-arm64'
                            }
                            environment {
                                GOARCH = 'arm64'
                            }
                            steps {
                                test()
                            }
                        }
                    }
                }
            }
        }
        stage('Build') {
            parallel {
                stage('windows') {
                    agent {
                        label 'windows'
                    }
                    environment {
                        GOOS = 'windows'
                        BINNAME = 'openitcockpit-agent.exe'
                    }
                    stages {
                        stage('amd64') {
                            environment {
                                GOARCH = 'amd64'
                            }
                            steps {
                                build_windows_binary()
                            }
                        }
                        stage('arm64') {
                            environment {
                                GOARCH = 'arm64'
                            }
                            steps {
                                build_windows_binary()
                            }
                        }
                        stage('386') {
                            environment {
                                GOARCH = '386'
                            }
                            steps {
                                build_windows_binary()
                            }
                        }
                    }
                }
                stage('linux') {
                    agent {
                        docker {
                            image 'golang:bookworm'
                            args "-u root --privileged -v agentgocache:/go"
                            label 'linux'
                        }
                    }
                    environment {
                        GOOS = 'linux'
                        BINNAME = 'openitcockpit-agent'
                        CGO_ENABLED = '0'
                    }
                    stages {
                        stage('amd64') {
                            environment {
                                GOARCH = 'amd64'
                            }
                            steps {
                                build_binary_linux()
                            }
                        }
                        stage('386') {
                            environment {
                                GOARCH = '386'
                            }
                            steps {
                                build_binary_linux()
                            }
                        }
                        stage('arm') {
                            environment {
                                GOARCH = 'arm'
                            }
                            steps {
                                build_binary_linux()
                            }

                        }
                        stage('arm64') {
                            environment {
                                GOARCH = 'arm64'
                            }
                            steps {
                                build_binary_linux()
                            }
                        }
                        // Build for PowerPC (Little-endian)
                        // Fedora had removed PowerPC (big endian) with Fedora 28 is it is probably mostly dead?
                        // https://fedoraproject.org/wiki/Architectures/PowerPC
                        // Looks like Debian still supports it (at least for now)
                        // https://wiki.debian.org/PPC64
                        // How ever at the moment we ignore this and only build Little-endian versions
                        stage('ppc64le') {
                            environment {
                                GOARCH = 'ppc64le'
                            }
                            steps {
                                build_binary_linux()
                            }
                        }
                    }
                }
                stage('freebsd') {
                    agent {
                        docker {
                            image 'golang:bookworm'
                            args "-u root --privileged -v agentgocache:/go"
                            label 'linux'
                        }
                    }
                    environment {
                        GOOS = 'freebsd'
                        BINNAME = 'openitcockpit-agent'
                        CGO_ENABLED = '0'
                    }
                    stages {
                        stage('amd64') {
                            environment {
                                GOARCH = 'amd64'
                            }
                            steps {
                                build_binary_freebsd()
                            }
                        }
                        stage('386') {
                            environment {
                                GOARCH = '386'
                            }
                            steps {
                                build_binary_freebsd()
                            }
                        }
                        stage('arm64') {
                            environment {
                                GOARCH = 'arm64'
                            }
                            steps {
                                build_binary_freebsd()
                            }
                        }
                    }
                }
                stage('darwin') {
                    environment {
                        GOOS = 'darwin'
                        BINNAME = 'openitcockpit-agent'
                    }
                    stages {
                        stage('amd64') {
                            agent {
                                label 'macos-intel64'
                            }
                            environment {
                                GOARCH = 'amd64'
                            }
                            steps {
                                build_binary_macos()
                            }
                        }
                        stage('arm64') {
                            agent {
                                label 'macos-arm64'
                            }
                            environment {
                                GOARCH = 'arm64'
                            }
                            steps {
                                build_binary_macos()
                            }
                        }
                    }
                }
            }
        }
        stage('Package') {
            environment {
                VERSION = sh(
                    returnStdout: true,
                    script: 'cat VERSION'
                ).trim()
            }
            parallel {
                stage('Linux') {
                    agent {
                        dockerfile {
                            filename 'linux.Dockerfile'
                            dir 'build/docker'
                            label 'linux'
                            args "-u root --privileged"
                        }
                    }
                    environment {
                        GOOS = 'linux'
                        BINNAME = 'openitcockpit-agent'
                    }
                    stages {
                        stage('amd64') {
                            environment {
                                GOARCH = 'amd64'
                                ARCH = 'amd64'
                                DEBARCH = 'amd64'
                                RPMARCH = 'amd64'
                            }
                            steps {
                                package_linux()
                            }
                        }
                        stage('386') {
                            environment {
                                GOARCH = '386'
                                ARCH = 'i386'
                                DEBARCH = 'i386'
                                RPMARCH = 'i386'
                            }
                            steps {
                                package_linux()
                            }
                        }
                        stage('arm64') {
                            environment {
                                GOARCH = 'arm64'
                                ARCH = 'arm64'
                                DEBARCH = 'arm64'
                                RPMARCH = 'aarch64'
                            }
                            steps {
                                package_linux()
                            }
                        }
                        stage('arm') {
                            environment {
                                GOARCH = 'arm'
                                ARCH = 'arm'
                                DEBARCH = 'armhf'
                                RPMARCH = 'armv7hl'
                            }
                            steps {
                                package_linux()
                            }
                        }
                        stage('ppc64le') {
                            environment {
                                GOARCH = 'ppc64le'
                                ARCH = 'ppc64le'
                                DEBARCH = 'ppc64el'
                                RPMARCH = 'ppc64le'
                            }
                            steps {
                                package_linux()
                            }
                        }
                    }
                }
                stage('windows') {
                    agent {
                        label 'windows'
                    }
                    environment {
                        branch = "${env.BRANCH_NAME}"
                        GOOS = 'windows'
                        BINNAME = 'openitcockpit-agent.exe'
                        //SM_CLIENT_CERT_FILE = "C:\\Users\\Administrator\\.digicert-ucpc\\certs\\e7499691-a05e-40d0-a8b5-4620fc5ebe82\\20260204102732-d51901f3-2cc4-4538-92e8-30cc33f94013.p12"
                        SM_CLIENT_CERT_FILE = "C:\\Users\\Administrator\\Certificate_pkcs12.p12"
                        SM_HOST = "https://clientauth.one.digicert.com"
                    }
                    stages {
                        stage('amd64') {
                            environment {
                                GOARCH = 'amd64'
                            }
                            steps {
                                // To get a new SM_CLIENT_CERT_FILE, SM_API_KEY or SM_CLIENT_CERT_PASSWORD
                                // Go to https://one.digicert.com/ click on the User Icon in the top right corner and select "Admin Profile"
                                // Now scroll down and create a new API Key and/or Client Certificate for Code Signing
                                // Also update the credentials in Jenkins with the new values
                                withCredentials([
                                    string(credentialsId: 'SM_API_KEY', variable: 'SM_API_KEY'),
                                    string(credentialsId: 'SM_CLIENT_CERT_PASSWORD', variable: 'SM_CLIENT_CERT_PASSWORD')
                                ]) {
                                    package_windows(branch)
                                }
                            }
                        }
                        stage('arm64') {
                            environment {
                                GOARCH = 'arm64'
                            }
                            steps {
                                withCredentials([
                                    string(credentialsId: 'SM_API_KEY', variable: 'SM_API_KEY'),
                                    string(credentialsId: 'SM_CLIENT_CERT_PASSWORD', variable: 'SM_CLIENT_CERT_PASSWORD')
                                ]) {
                                    package_windows(branch)
                                }
                            }
                        }
                        stage('386') {
                            environment {
                                GOARCH = '386'
                            }
                            steps {
                                withCredentials([
                                    string(credentialsId: 'SM_API_KEY', variable: 'SM_API_KEY'),
                                    string(credentialsId: 'SM_CLIENT_CERT_PASSWORD', variable: 'SM_CLIENT_CERT_PASSWORD')
                                ]) {
                                    package_windows(branch)
                                }
                            }
                        }
                    }
                }
                stage('darwin') {
                    environment {
                        GOOS = 'darwin'
                        BINNAME = 'openitcockpit-agent'
                    }
                    stages {
                        stage('amd64') {
                            agent {
                                label 'macos-intel64'
                            }
                            environment {
                                ARCH = 'amd64'
                                GOARCH = 'amd64'
                            }
                            steps {
                                withCredentials([
                                    string(credentialsId: 'INTEL_MAC_PASSWORD', variable: 'INTEL_MAC_PASSWORD')
                                ]) {
                                    package_darwin_amd64()
                                }
                            }
                        }
                        stage('arm64') {
                            agent {
                                label 'macos-arm64'
                            }
                            environment {
                                ARCH = 'arm64'
                                GOARCH = 'arm64'
                            }
                            steps {
                                withCredentials([
                                    string(credentialsId: 'ARM_MAC_PASSWORD', variable: 'ARM_MAC_PASSWORD')
                                ]) {
                                    package_darwin_arm64()
                                }
                            }
                        }
                    }
                }
                stage('FreeBSD') {
                    agent {
                        dockerfile {
                            filename 'linux.Dockerfile'
                            dir 'build/docker'
                            label 'linux'
                            args "-u root --privileged"
                        }
                    }
                    environment {
                        GOOS = 'freebsd'
                        BINNAME = 'openitcockpit-agent'
                    }
                    stages {
                        stage('amd64') {
                            environment {
                                GOARCH = 'amd64'
                                ARCH = 'amd64'
                            }
                            steps {
                                package_freebsd()
                            }
                        }
                        stage('386') {
                            environment {
                                GOARCH = '386'
                                ARCH = '386'
                            }
                            steps {
                                package_freebsd()
                            }
                        }
                        stage('arm64') {
                            environment {
                                GOARCH = 'arm64'
                                ARCH = 'arm64'
                            }
                            steps {
                                package_freebsd()
                            }
                        }
                    }
                }
            }
        }
        stage('Publish') {
            when {
                branch 'main'
            }
            environment {
                VERSION = sh(
                    returnStdout: true,
                    script: 'cat VERSION'
                ).trim()
            }
            agent {
                label 'linux'
            }
            steps {
                publish_packages()
            }
         }
    }
}

def cleanup_windows() {
    bat 'git.exe clean -f -d -x'
}

def cleanup() {
    sh 'git config --global --add safe.directory "$PWD"'
    sh 'git clean -f -d -x'
}

def test_windows() {
    timeout(time: 10, unit: 'MINUTES') {
        cleanup_windows()

        catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
            if (env.GOARCH == 'arm64') {
                // The arm64 test binaries can not be executed on the amd64 build node, so the code and the tests are only compiled and vetted
                bat "go.exe vet ./..."
            } else {
                bat "go.exe test -timeout=120s ./..."
            }
        }
    }
}

def test() {
    timeout(time: 10, unit: 'MINUTES') {
        cleanup()

        catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
            sh "go test -timeout=120s ./..."
        }
    }
}


def build_windows_binary() {
    timeout(time: 5, unit: 'MINUTES') {
        cleanup_windows()

        catchError(buildResult: null, stageResult: 'FAILURE') {
            bat "mkdir release\\$GOOS\\$GOARCH"
            bat "go.exe build -o release/$GOOS/$GOARCH/$BINNAME main.go"

            // ITC-3498 Contain information about 3rd party licenses in the package
            bat "go.exe install github.com/google/go-licenses/v2@latest"
            bat "%userprofile%\\go\\bin\\go-licenses.exe report . --ignore 'github.com/openITCOCKPIT/openitcockpit-agent-go' > release/$GOOS/$GOARCH/licenses.csv"
            bat "%userprofile%\\go\\bin\\go-licenses.exe check . --ignore \"github.com/openITCOCKPIT/openitcockpit-agent-go\" --allowed_licenses=MIT,ISC,MPL-2.0,BSD-2-Clause,BSD-3-Clause,Apache-2.0"
        }
        archiveArtifacts artifacts: "release/$GOOS/$GOARCH/**", fingerprint: true
        stash name: "release-$GOOS-$GOARCH", includes: "release/$GOOS/$GOARCH/**"
    }
}

def build_binary_linux() {
    timeout(time: 5, unit: 'MINUTES') {
        cleanup()

        catchError(buildResult: null, stageResult: 'FAILURE') {
            sh "mkdir -p release/$GOOS/$GOARCH"
            sh "go build -o release/$GOOS/$GOARCH/$BINNAME main.go"

            // ITC-3498 Contain information about 3rd party licenses in the package
            sh "go install github.com/google/go-licenses/v2@latest"
            sh "/go/bin/go-licenses report . --ignore \"github.com/openITCOCKPIT/openitcockpit-agent-go\" > release/$GOOS/$GOARCH/licenses.csv"
            sh "/go/bin/go-licenses check . --ignore \"github.com/openITCOCKPIT/openitcockpit-agent-go\" --allowed_licenses=MIT,ISC,MPL-2.0,BSD-2-Clause,BSD-3-Clause,Apache-2.0"
        }
        archiveArtifacts artifacts: "release/$GOOS/$GOARCH/**", fingerprint: true
        stash name: "release-$GOOS-$GOARCH", includes: "release/$GOOS/$GOARCH/**"
    }
}

def build_binary_freebsd() {
    timeout(time: 5, unit: 'MINUTES') {
        cleanup()

        catchError(buildResult: null, stageResult: 'FAILURE') {
            sh "mkdir -p release/$GOOS/$GOARCH"
            sh "go build -o release/$GOOS/$GOARCH/$BINNAME main.go"

            // ITC-3498 Contain information about 3rd party licenses in the package
            sh "go install github.com/google/go-licenses/v2@latest"
            sh "/go/bin/go-licenses report . --ignore \"github.com/openITCOCKPIT/openitcockpit-agent-go\" > release/$GOOS/$GOARCH/licenses.csv"
            sh "/go/bin/go-licenses check . --ignore \"github.com/openITCOCKPIT/openitcockpit-agent-go\" --allowed_licenses=MIT,ISC,MPL-2.0,BSD-2-Clause,BSD-3-Clause,Apache-2.0"
        }
        archiveArtifacts artifacts: "release/$GOOS/$GOARCH/**", fingerprint: true
        stash name: "release-$GOOS-$GOARCH", includes: "release/$GOOS/$GOARCH/**"
    }
}


def build_binary_macos() {
    timeout(time: 5, unit: 'MINUTES') {
        cleanup()

        catchError(buildResult: null, stageResult: 'FAILURE') {
            sh "mkdir -p release/$GOOS/$GOARCH"
            sh "go build -o release/$GOOS/$GOARCH/$BINNAME main.go"

            // ITC-3498 Contain information about 3rd party licenses in the package
            sh "go install github.com/google/go-licenses/v2@latest"
            sh "$HOME/go/bin/go-licenses report . --ignore \"github.com/openITCOCKPIT/openitcockpit-agent-go\" > release/$GOOS/$GOARCH/licenses.csv"
            sh "$HOME/go/bin/go-licenses check . --ignore \"github.com/openITCOCKPIT/openitcockpit-agent-go\" --allowed_licenses=MIT,ISC,MPL-2.0,BSD-2-Clause,BSD-3-Clause,Apache-2.0"
        }
        archiveArtifacts artifacts: "release/$GOOS/$GOARCH/**", fingerprint: true
        stash name: "release-$GOOS-$GOARCH", includes: "release/$GOOS/$GOARCH/**"
    }
}

def package_linux() {
    timeout(time: 5, unit: 'MINUTES') {
        cleanup()

        unstash name: "release-$GOOS-$GOARCH"

        sh "mkdir -p package/usr/bin package/etc/openitcockpit-agent/init release/packages/$GOOS package/var/log/openitcockpit-agent package/usr/share/doc/openitcockpit-agent"
        sh 'cp example/config_example.ini package/etc/openitcockpit-agent/config.ini'
        sh 'cp example/customchecks_example.ini package/etc/openitcockpit-agent/customchecks.ini'
        sh 'cp example/prometheus_exporters_example.ini package/etc/openitcockpit-agent/prometheus_exporters.ini'
        sh 'cp example/logfiles_example.ini package/etc/openitcockpit-agent/logfiles.ini'
        sh 'cp example/certificates_example.ini package/etc/openitcockpit-agent/certificates.ini'
        sh 'cp example/files_example.ini package/etc/openitcockpit-agent/files.ini'
        sh 'cp example/http_example.ini package/etc/openitcockpit-agent/http.ini'
        sh 'cp example/redfish_example.ini package/etc/openitcockpit-agent/redfish.ini'
        // The configuration files contain passwords (basic auth, API key of the push mode, BMCs of redfish.ini) and
        // must only be readable by root. Existing files are not changed on updates (config files of the packages).
        sh 'chmod 600 package/etc/openitcockpit-agent/*.ini'
        sh 'cp build/package/openitcockpit-agent.init package/etc/openitcockpit-agent/init/openitcockpit-agent.init'
        sh 'cp build/package/openitcockpit-agent.service package/etc/openitcockpit-agent/init/openitcockpit-agent.service'
        sh "cp release/linux/$GOARCH/$BINNAME package/usr/bin/$BINNAME"
        sh "cp release/linux/$GOARCH/licenses.csv package/usr/share/doc/openitcockpit-agent/licenses.csv"
        sh "chmod +x package/usr/bin/$BINNAME"
        sh "chmod +x package/etc/openitcockpit-agent/init/openitcockpit-agent.init"
        sh """cd release/packages/$GOOS &&
            fpm -s dir -t deb -C ../../../package --name openitcockpit-agent --vendor 'AVENDIS GmbH' \\
            --license 'Apache License Version 2.0' --config-files etc/openitcockpit-agent \\
            --architecture $DEBARCH --maintainer '<d.ziegler@avendis.com>' \\
            --description 'openITCOCKPIT Monitoring Agent and remote plugin executor.' \\
            --url 'https://openitcockpit.io' --before-install ../../../build/package/preinst.sh \\
            --after-install ../../../build/package/postinst.sh --before-remove ../../../build/package/prerm.sh  \\
            --version '$VERSION'"""
        sh """cd release/packages/$GOOS &&
            fpm -s dir -t rpm -C ../../../package --name openitcockpit-agent --vendor 'AVENDIS GmbH' \\
            --license 'Apache License Version 2.0' --config-files etc/openitcockpit-agent \\
            --architecture $RPMARCH --maintainer '<d.ziegler@avendis.com>' \\
            --description 'openITCOCKPIT Monitoring Agent and remote plugin executor.' \\
            --url 'https://openitcockpit.io' --before-install ../../../build/package/preinst.sh \\
            --after-install ../../../build/package/postinst.sh --before-remove ../../../build/package/prerm.sh  \\
            --version '$VERSION'"""
        sh """cd release/packages/$GOOS &&
            fpm -s dir -t pacman -C ../../../package --name openitcockpit-agent --vendor 'AVENDIS GmbH' \\
            --license 'Apache License Version 2.0' --config-files etc/openitcockpit-agent \\
            --architecture $ARCH --maintainer '<d.ziegler@avendis.com>' \\
            --description 'openITCOCKPIT Monitoring Agent and remote plugin executor.' \\
            --url 'https://openitcockpit.io' --before-install ../../../build/package/preinst.sh \\
            --after-install ../../../build/package/postinst.sh --before-remove ../../../build/package/prerm.sh  \\
            --version '$VERSION'"""

        archiveArtifacts artifacts: 'release/packages/**', fingerprint: true
    }
}

def package_windows(branch) {
    timeout(time: 5, unit: 'MINUTES') {
        cleanup_windows()

        unstash name: "release-$GOOS-$GOARCH"

        // Convert Linux new lines to Windows new lines for older Windows Server systems
        bat 'move example\\certificates_example.ini example\\certificates_example_linux.ini'
        bat 'TYPE example\\certificates_example_linux.ini | MORE /P > example\\certificates_example.ini'

        bat 'move example\\config_example.ini example\\config_example_linux.ini'
        bat 'TYPE example\\config_example_linux.ini | MORE /P > example\\config_example.ini'

        bat 'move example\\customchecks_example.ini example\\customchecks_example_linux.ini'
        bat 'TYPE example\\customchecks_example_linux.ini | MORE /P > example\\customchecks_example.ini'

        bat 'move example\\files_example.ini example\\files_example_linux.ini'
        bat 'TYPE example\\files_example_linux.ini | MORE /P > example\\files_example.ini'

        bat 'move example\\http_example.ini example\\http_example_linux.ini'
        bat 'TYPE example\\http_example_linux.ini | MORE /P > example\\http_example.ini'

        bat 'move example\\logfiles_example.ini example\\logfiles_example_linux.ini'
        bat 'TYPE example\\logfiles_example_linux.ini | MORE /P > example\\logfiles_example.ini'

        bat 'move example\\perfcounters_example.ini example\\perfcounters_example_linux.ini'
        bat 'TYPE example\\perfcounters_example_linux.ini | MORE /P > example\\perfcounters_example.ini'

        bat 'move example\\prometheus_exporters_example.ini example\\prometheus_exporters_linux.ini'
        bat 'TYPE example\\prometheus_exporters_linux.ini | MORE /P > example\\prometheus_exporters_example.ini'

        bat 'move example\\redfish_example.ini example\\redfish_example_linux.ini'
        bat 'TYPE example\\redfish_example_linux.ini | MORE /P > example\\redfish_example.ini'

        powershell "& \"C:\\Program Files\\DigiCert\\DigiCert One Signing Manager Tools\\smctl.exe\" healthcheck"

        //powershell "& $ADVINST /loadpathvars \"build\\msi\\PathVariables_Jenkins.apf\""
        powershell "& $ADVINST /edit \"build\\msi\\openitcockpit-agent-${GOARCH}.aip\" \\NewPathVariable -name AGENT_SOURCE -value \"${WORKSPACE}\" -valuetype Folder -global"
        powershell "& $ADVINST /edit \"build\\msi\\openitcockpit-agent-${GOARCH}.aip\" \\SetVersion \"$VERSION\""

        if (env.BRANCH_NAME != 'main' && !params.SIGN_WINDOWS) {
            // Disable code-signing for development builds branches, unless the build parameter SIGN_WINDOWS is set
            powershell "& $ADVINST /edit \"build\\msi\\openitcockpit-agent-${GOARCH}.aip\" \\ResetSig"
        }

        powershell "& $ADVINST /build \"build\\msi\\openitcockpit-agent-${GOARCH}.aip\""
        archiveArtifacts artifacts: 'release/packages/**', fingerprint: true
    }
}

def package_darwin_amd64() {
    timeout(time: 5, unit: 'MINUTES') {
        cleanup()

        unstash name: "release-$GOOS-$GOARCH"

        sh "mkdir -p package/Applications/openitcockpit-agent package_osx_uninstaller release/packages/$GOOS"
        sh "cp release/$GOOS/$GOARCH/$BINNAME package/Applications/openitcockpit-agent/"
        sh "cp release/$GOOS/$GOARCH/licenses.csv package/Applications/openitcockpit-agent/"
        sh "cp example/config_example.ini package/Applications/openitcockpit-agent/config.ini"
        sh "cp example/customchecks_example.ini package/Applications/openitcockpit-agent/customchecks.ini"
        sh "cp example/prometheus_exporters_example.ini package/Applications/openitcockpit-agent/prometheus_exporters.ini"
        sh "cp build/package/io.openitcockpit.agent.plist package/Applications/openitcockpit-agent/io.openitcockpit.agent.plist"
        sh "chmod +x package/Applications/openitcockpit-agent/$BINNAME"

        // Unlock macOS Keychain to access the code signing certificates
        sh """security unlock-keychain -p '$INTEL_MAC_PASSWORD' login.keychain"""
        sh """security set-keychain-settings -t 3600 -u login.keychain"""

        // Sign the binary (Developer ID Application)
        sh "codesign --force --options runtime --timestamp --sign 'EEFC9653F25EDBC2397EE55B33EF2E426A2DA394' ./package/Applications/openitcockpit-agent/$BINNAME"

        sh """/usr/local/bin/packagesbuild --package-version "${VERSION}" --reference-folder . build/macos/openITCOCKPIT\\ Monitoring\\ Agent/openITCOCKPIT\\ Monitoring\\ Agent.pkgproj"""
        sh """mv -f openitcockpit-agent-darwin-amd64.pkg release/packages/${GOOS}/openitcockpit-agent-${VERSION}-darwin-${GOARCH}_unsigned.pkg"""

        // Sign the pkg-installer (Developer ID Installer)
        sh """productsign --sign "83DD3DEAE7391DEC3F92567E83A4F3ABAEC8A421" ./release/packages/${GOOS}/openitcockpit-agent-${VERSION}-darwin-${GOARCH}_unsigned.pkg ./release/packages/${GOOS}/openitcockpit-agent-${VERSION}-darwin-${GOARCH}.pkg"""

        // Sent to Apple for notarization
        if (env.BRANCH_NAME == 'main') {
            // Enable notarization for main branch
            sh """xcrun notarytool submit ./release/packages/${GOOS}/openitcockpit-agent-${VERSION}-darwin-${GOARCH}.pkg --keychain-profile "agent-notary-profile" --wait"""
            // Attach ticket for offline verification on other machines
            sh """xcrun stapler staple ./release/packages/${GOOS}/openitcockpit-agent-${VERSION}-darwin-${GOARCH}.pkg"""
        }

        sh """rm -f ./release/packages/${GOOS}/openitcockpit-agent-${VERSION}-darwin-${GOARCH}_unsigned.pkg"""


        sh """cd release/packages/$GOOS &&
            fpm -s dir -t osxpkg -C ../../../package_osx_uninstaller --name openitcockpit-agent-uninstaller --vendor "AVENDIS GmbH" \\
            --license "Apache License Version 2.0" --config-files Applications/openitcockpit-agent \\
            --maintainer "<d.ziegler@avendis.com>" \\
            --description "Uninstaller of openITCOCKPIT Monitoring Agent and remote plugin executor." --url "https://openitcockpit.io" \\
            --before-install ../../../build/package/macos/prerm.sh --version '$VERSION' --osxpkg-payload-free &&
            mv openitcockpit-agent-uninstaller-${VERSION}.pkg openitcockpit-agent-uninstaller-${VERSION}-darwin-all.pkg"""

        archiveArtifacts artifacts: 'release/packages/**', fingerprint: true

        /*sh """cd release/packages/$GOOS &&
            fpm -s dir -t osxpkg -C ../../../package --name openitcockpit-agent --vendor 'it-novum GmbH' \\
            --license "Apache License Version 2.0" --config-files Applications/openitcockpit-agent \\
            --architecture $ARCH --maintainer "<daniel.ziegler@it-novum.com>" \\
            --description "openITCOCKPIT Monitoring Agent and remote plugin executor." \\
            --url "https://openitcockpit.io" --before-install ../../../build/package/preinst.sh \\
            --after-install ../../../build/package/postinst.sh --version '$VERSION' &&
            mv openitcockpit-agent-${VERSION}.pkg openitcockpit-agent-${VERSION}-darwin-${GOARCH}.pkg"""
        sh """cd release/packages/$GOOS &&
            fpm -s dir -t osxpkg -C ../../../package_osx_uninstaller --name openitcockpit-agent-uninstaller --vendor "it-novum GmbH" \\
            --license "Apache License Version 2.0" --config-files Applications/openitcockpit-agent \\
            --architecture $ARCH --maintainer "<daniel.ziegler@it-novum.com>" \\
            --description "openITCOCKPIT Monitoring Agent and remote plugin executor." --url "https://openitcockpit.io" \\
            --before-install ../../../build/package/prerm.sh --version '$VERSION' --osxpkg-payload-free &&
            mv openitcockpit-agent-uninstaller-${VERSION}.pkg openitcockpit-agent-uninstaller-${VERSION}-darwin-${GOARCH}.pkg"""
        archiveArtifacts artifacts: 'release/packages/**', fingerprint: true*/
    }
}

def package_darwin_arm64() {
    timeout(time: 5, unit: 'MINUTES') {
        cleanup()

        unstash name: "release-$GOOS-$GOARCH"

        sh "mkdir -p package/Applications/openitcockpit-agent package_osx_uninstaller release/packages/$GOOS"
        sh "cp release/$GOOS/$GOARCH/$BINNAME package/Applications/openitcockpit-agent/"
        sh "cp release/$GOOS/$GOARCH/licenses.csv package/Applications/openitcockpit-agent/"
        sh "cp example/config_example.ini package/Applications/openitcockpit-agent/config.ini"
        sh "cp example/customchecks_example.ini package/Applications/openitcockpit-agent/customchecks.ini"
        sh "cp example/prometheus_exporters_example.ini package/Applications/openitcockpit-agent/prometheus_exporters.ini"
        sh "cp build/package/io.openitcockpit.agent.plist package/Applications/openitcockpit-agent/io.openitcockpit.agent.plist"
        sh "chmod +x package/Applications/openitcockpit-agent/$BINNAME"

        // Unlock macOS Keychain to access the code signing certificates
        sh """security unlock-keychain -p "$ARM_MAC_PASSWORD" login.keychain"""
        sh """security set-keychain-settings -t 3600 -u login.keychain"""

        // Sign the binary (Developer ID Application)
        sh "codesign --force --options runtime --timestamp --sign 'EEFC9653F25EDBC2397EE55B33EF2E426A2DA394' ./package/Applications/openitcockpit-agent/$BINNAME"

        sh """/usr/local/bin/packagesbuild --package-version "${VERSION}" --reference-folder . build/macos/openITCOCKPIT\\ Monitoring\\ Agent\\ arm64/openITCOCKPIT\\ Monitoring\\ Agent.pkgproj"""
        sh """mv -f openitcockpit-agent-darwin-arm64.pkg release/packages/${GOOS}/openitcockpit-agent-${VERSION}-darwin-${GOARCH}_unsigned.pkg"""

        // Sign the pkg-installer (Developer ID Installer)
        sh """productsign --sign "83DD3DEAE7391DEC3F92567E83A4F3ABAEC8A421" ./release/packages/${GOOS}/openitcockpit-agent-${VERSION}-darwin-${GOARCH}_unsigned.pkg ./release/packages/${GOOS}/openitcockpit-agent-${VERSION}-darwin-${GOARCH}.pkg"""

        // Sent to Apple for notarization
        if (env.BRANCH_NAME == 'main') {
            // Enable notarization for main branch
            sh """xcrun notarytool submit ./release/packages/${GOOS}/openitcockpit-agent-${VERSION}-darwin-${GOARCH}.pkg --keychain-profile "agent-notary-profile" --wait"""
            // Attach ticket for offline verification on other machines
            sh """xcrun stapler staple ./release/packages/${GOOS}/openitcockpit-agent-${VERSION}-darwin-${GOARCH}.pkg"""
        }

        sh """rm -f ./release/packages/${GOOS}/openitcockpit-agent-${VERSION}-darwin-${GOARCH}_unsigned.pkg"""

        archiveArtifacts artifacts: 'release/packages/**', fingerprint: true
    }
}

def package_freebsd() {
    timeout(time: 5, unit: 'MINUTES') {
        cleanup()

        unstash name: "release-$GOOS-$GOARCH"

        // For now we will just publish the binary without packaging it
        sh "mkdir -p release/packages/freebsd"
        sh "cp release/$GOOS/$GOARCH/$BINNAME release/packages/$GOOS/$BINNAME-${VERSION}.${GOOS}-${ARCH}.bin"

        archiveArtifacts artifacts: 'release/packages/**', fingerprint: true
    }
}

def publish_packages() {
    timeout(time: 5, unit: 'MINUTES') {

        dir('publish'){
        /* get all packages */
            unarchive mapping: ['release/packages/' : '.']

            sh """mkdir packages"""
            sh """mv -f release/packages/darwin/* packages/"""
            sh """mv -f release/packages/linux/* packages/"""
            sh """mv -f release/packages/windows/* packages/"""
            sh """mv -f release/packages/freebsd/* packages/"""

            sh 'ssh -o StrictHostKeyChecking=no -i /var/lib/jenkins/.ssh/id_rsa oitc@srvitnweb05.static.itsm.love "mkdir -p /var/www/openitcockpit.io/files/openitcockpit-agent-3.x"'
            sh 'rsync -avz -e "ssh -o StrictHostKeyChecking=no -i /var/lib/jenkins/.ssh/id_rsa" --delete --progress packages/* oitc@srvitnweb05.static.itsm.love:/var/www/openitcockpit.io/files/openitcockpit-agent-3.x/'

            /* Remove old packages */
            sh '/var/lib/jenkins/openITCOCKPIT-build/aptly.sh repo remove openitcockpit-agent-stable openitcockpit-agent'
            sh '/var/lib/jenkins/openITCOCKPIT-build/aptly.sh publish drop deb filesystem:openitcockpit-agent:deb/stable'

            /* Add packages to own apt repository */
            sh '/var/lib/jenkins/openITCOCKPIT-build/aptly.sh repo add -force-replace openitcockpit-agent-stable packages/openitcockpit-agent_*.deb;'

            /* Publish apt repository */
            sh '/var/lib/jenkins/openITCOCKPIT-build/aptly.sh publish repo -distribution deb -architectures amd64,i386,arm64,armhf,ppc64el -passphrase-file /opt/repository/aptly/pw -batch openitcockpit-agent-stable filesystem:openitcockpit-agent:deb/stable'
            /* apt02 = https://packages.openitcockpit.io/ */
            /* apt03 = https://packages5.openitcockpit.io/ */
            sh "rsync -rv --delete-after /opt/repository/aptly/openitcockpit-agent/deb/stable/ www-data@srvoitcapt02.static.itsm.love:/var/www/html/openitcockpit-agent/deb/stable/"
            sh "rsync -rv --delete-after /opt/repository/aptly/openitcockpit-agent/deb/stable/ www-data@srvoitcapt03.static.itsm.love:/var/www/html/openitcockpit-agent/deb/stable/"

            /* sign rpm packages */
            sh """rpmsign --define "_gpg_name it-novum GmbH <support@itsm.it-novum.com>" --define "__gpg_sign_cmd %{__gpg} gpg --no-verbose --no-armor --batch --pinentry-mode loopback --passphrase-file /opt/repository/aptly/pw %{?_gpg_digest_algo:--digest-algo %{_gpg_digest_algo}} --no-secmem-warning -u '%{_gpg_name}' -sbo %{__signature_filename} %{__plaintext_filename}" --addsign packages/*.rpm"""

            /* create yum repository */
            sh """mkdir -p rpm/stable"""
            sh """cp packages/*.rpm rpm/stable/"""
            /*
            createrepo is the old Python createrepo which does not exists on Ubuntu Focal.
            Ubuntu Focal has no createrepo_c.
            For this reason we use a Docker image with a recent version of createrepo_c which works fine
            https://github.com/openITCOCKPIT/createrepo_c-docker
            */
            /*sh """createrepo rpm/stable"""*/
            sh """docker run --rm --user 111:116 -v '${env.WORKSPACE}/publish/rpm/stable':/rpm openitcockpit/createrepo_c"""
            /*sh """chown jenkins:jenkins ${env.WORKSPACE} -R"""*/

            /* Sign Repository meta data with GPG for SLES systems */
            /* this will create the repomd.xml.asc file */
            sh """
            gpg --batch --yes --pinentry-mode loopback \
                --passphrase-file /opt/repository/aptly/pw \
                --local-user "it-novum GmbH <support@itsm.it-novum.com>" \
                --detach-sign --armor \
                ${env.WORKSPACE}/publish/rpm/stable/repodata/repomd.xml
            """

            /* Publish yum repository */
            /* apt02 = https://packages.openitcockpit.io/ */
            /* apt03 = https://packages5.openitcockpit.io/ */
            sh "rsync -rv --delete-after rpm/stable/ www-data@srvoitcapt02.static.itsm.love:/var/www/html/openitcockpit-agent/rpm/stable/"
            sh "rsync -rv --delete-after rpm/stable/ www-data@srvoitcapt03.static.itsm.love:/var/www/html/openitcockpit-agent/rpm/stable/"
        }

    }
}

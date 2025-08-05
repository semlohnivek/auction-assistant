// internal/sftp.go
package internal

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"

	"bidzauction/config"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func UploadImagesToSFTP(localDir, remoteDir string, filenames []string) error {
	cfg := config.Current.SFTP

	sshConfig := &ssh.ClientConfig{
		User: cfg.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(cfg.Pass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	conn, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to SFTP server: %w", err)
	}
	defer conn.Close()

	sftpClient, err := sftp.NewClient(conn)
	if err != nil {
		return fmt.Errorf("failed to create SFTP client: %w", err)
	}
	defer sftpClient.Close()

	// Ensure remote directory exists
	if err := sftpClient.MkdirAll(remoteDir); err != nil {
		return fmt.Errorf("failed to ensure remote directory: %w", err)
	}

	for _, name := range filenames {
		localPath := filepath.Join(localDir, name)
		remotePath := path.Join(remoteDir, name)

		srcFile, err := os.Open(localPath)
		if err != nil {
			fmt.Printf("Skipping %s: cannot open (%v)\n", name, err)
			continue
		}
		dstFile, err := sftpClient.Create(remotePath)
		if err != nil {
			fmt.Printf("Skipping %s: cannot create remote file (%v)\n", name, err)
			srcFile.Close()
			continue
		}
		_, err = io.Copy(dstFile, srcFile)
		srcFile.Close()
		dstFile.Close()
		if err != nil {
			fmt.Printf("Error copying %s: %v\n", name, err)
			continue
		}
		fmt.Printf("Uploaded: %s\n", name)
	}

	return nil
}

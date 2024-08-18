package test

import (
	"fmt"
	"io/ioutil"
	"flag"
	"testing"
	"strings"

	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"github.com/go-sql-driver/mysql"

	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"golang.org/x/crypto/ssh"

)

var folder = flag.String("folder", "", "Folder ID in Yandex.Cloud")
var sshKeyPath = flag.String("ssh-key-pass", "", "Private ssh key for access to virtual machines")

func TestEndToEndDeploymentScenario(t *testing.T) {
    fixtureFolder := "../"

    test_structure.RunTestStage(t, "setup", func() {
		terraformOptions := &terraform.Options{
			TerraformDir: fixtureFolder,

			Vars: map[string]interface{}{
			"yc_folder":    *folder,
		    },
	    }

		test_structure.SaveTerraformOptions(t, fixtureFolder, terraformOptions)

		terraform.InitAndApply(t, terraformOptions)
	})

	test_structure.RunTestStage(t, "validate", func() {
	    fmt.Println("Run some tests...")

	    terraformOptions := test_structure.LoadTerraformOptions(t, fixtureFolder)

 	dbhosts := terraform.Output(t, terraformOptions, "database_host_fqdn")

        // test load balancer ip existing
	    loadbalancerIPAddress := terraform.Output(t, terraformOptions, "load_balancer_public_ip")

	    if loadbalancerIPAddress == "" {
			t.Fatal("Cannot retrieve the public IP address value for the load balancer.")
		}

		// test ssh connect
		vmLinuxPublicIPAddress := terraform.Output(t, terraformOptions, "vm_linux_public_ip_address")

		res1 := strings.Split(string(vmLinuxPublicIPAddress), " ")
		res2 := strings.Split(string(dbhosts), " ")


		key, err := ioutil.ReadFile(*sshKeyPath)
		if err != nil {
			t.Fatalf("Unable to read private key: %v", err)
		}

		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			t.Fatalf("Unable to parse private key: %v", err)
		}

		sshConfig := &ssh.ClientConfig{
			User: "ubuntu",
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(signer),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}

		//1 host

		sshConnection, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", res1[0] + "]"), sshConfig)
		if err != nil {
			t.Fatalf("Cannot establish SSH connection to vm-linux public IP address: %v", err)
		}

		defer sshConnection.Close()

		sshSession, err := sshConnection.NewSession()
		if err != nil {
			t.Fatalf("Cannot create SSH session to vm-linux public IP address: %v", err)
		}

		defer sshSession.Close()

		err = sshSession.Run(fmt.Sprintf("ping -c 1 8.8.8.8"))
		if err != nil {
			t.Fatalf("Cannot ping 8.8.8.8: %v", err)
		}


		//DB Test

		host := res2[0] + "]"
		user := terraform.Output(t, terraformOptions, "db_user")
		dbname := terraform.Output(t, terraformOptions, "db_name")
		password := terraform.Output(t, terraformOptions, "db_pass")

		const (
			port     = 3306
		)


		rootCertPool := x509.NewCertPool()
  		pem, err := ioutil.ReadFile("./.mysql/root.crt")
  		if err != nil {
		  t.Fatalf("No cert for DB: %v", err)
  		}
  		if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
    		  t.Fatalf("Failed to append PEM: %v", err)
  		}

  		mysql.RegisterTLSConfig("custom", &tls.Config{
    		  RootCAs: rootCertPool,
  		})

  		mysqlInfo := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?tls=custom",
    		  user, password, host, port, dbname)
  		  conn, err := sql.Open("mysql", mysqlInfo)
  		  if err != nil {
    		    t.Fatalf("Error connection for DB: %v", err)
  		  }

  		defer conn.Close()

  		q, err := conn.Query("SELECT version()")
  		if err != nil {
    		  t.Fatalf("No run query: %v", err)
  		}


		var result string

  		for q.Next() {
		  q.Scan(&result)
		  fmt.Println(result)
		}

		fmt.Println("Version DB: %v \n", q)



	})

	test_structure.RunTestStage(t, "teardown", func() {
		terraformOptions := test_structure.LoadTerraformOptions(t, fixtureFolder)
		terraform.Destroy(t, terraformOptions)
	})
}

package vault

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.wdf.sap.corp/ICN-ML/aicore/system-services/platform/pkg/exec"
	"github.wdf.sap.corp/ICN-ML/aicore/system-services/platform/pkg/log"
)

type config struct {
	aiVaultAddress string
	appRoleID      string
	secretID       string
	secretPath     string
	vaultToken     string
	jenkinsEnv     bool
}

type awsConfig struct {
	data map[string]map[string]string
}

var conf *config

// Client is public interface for interaction with the vault package
type Client interface {
	GetSecret() map[string]string
	CreateOrUpdateSecret(key iamtypes.AccessKey)
	DeleteSecret()
}

// New is the entry function for the vault package. It creates a new Vault client
// to interact with.
// Input:
//  - appRoleID: appRole credential defined in the jenkins pipeline
//  - secretID: secretId credential defined in
func New(
	appRoleID string,
	secretID string,
	secretPath string,
	jenkinsEnv bool,
) Client {
	conf = &config{
		aiVaultAddress: "https://vault.ml.only.sap",
		appRoleID:      appRoleID,
		secretID:       secretID,
		secretPath:     secretPath,
		jenkinsEnv:     jenkinsEnv,
	}

	if jenkinsEnv {
		conf.vaultToken = conf.getVaultToken()
	} else {
		os.Setenv("VAULT_ADDR", conf.aiVaultAddress)
		token, err := getLocalVaultToken()
		if err != nil {
			log.Sugar.Error("please check if you are logged into vault")
			return nil
		}
		log.Sugar.Debug(token)
		conf.vaultToken = token
	}
	os.Setenv("VAULT_TOKEN", conf.vaultToken)
	log.Sugar.Debug("vault client successfully created")
	return conf
}

func getLocalVaultToken() (string, error) {
	ctx := context.Background()
	out, err := exec.New(ctx).Command("vault", "token", "lookup").Output()

	if err != nil {
		return "", err
	}
	vals := strings.Split(strings.ReplaceAll(string(out), "\r\n", "\n"), "\n")

	for _, val := range vals {
		keyVal := strings.Fields(val)
		if keyVal[0] == "id" {
			return keyVal[1], nil
		}
	}

	return "", fmt.Errorf("no client token found while lookup")
}

func (conf *config) getVaultToken() string {
	postBody, _ := json.Marshal(map[string]string{
		"role_id":   conf.appRoleID,
		"secret_id": conf.secretID,
	})

	resBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(fmt.Sprintf("%s/v1/auth/approle/login", conf.aiVaultAddress), "application/json", resBody)
	if err != nil {
		log.Sugar.Error(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Sugar.Error(err)
	}
	m := make(map[string]interface{})
	err = json.Unmarshal(body, &m)

	if err != nil {
		log.Sugar.Error(err)
	}
	vaultToken := m["auth"].(map[string]interface{})["client_token"]
	return vaultToken.(string)
}

func (conf *config) GetSecret() map[string]string {
	log.Sugar.Debug("Getting secrets")
	if conf.jenkinsEnv {

		url := fmt.Sprintf("%s/v1/secret/data/%s", conf.aiVaultAddress, conf.secretPath)

		client := &http.Client{}
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("X-Vault-Token", conf.vaultToken)
		res, err := client.Do(req)
		if err != nil {
			log.Sugar.Error(err)
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Sugar.Error(err)
		}
		g := make(map[string]interface{})
		err = json.Unmarshal(body, &g)
		if err != nil {
			log.Sugar.Error(err)
		}
		return g["data"].(map[string]interface{})["data"].(map[string]string)
	} else {
		homedir, _ := os.UserHomeDir()
		awsCreds, err := readAwsConfig(fmt.Sprintf("%s\\.aws\\credentials", homedir))
		if err != nil {
			log.Sugar.Error(err)
			return nil
		}
		creds := make(map[string]string)
		creds["access-key"] = awsCreds.data["default"]["aws_access_key_id"]
		creds["secret-key"] = awsCreds.data["default"]["aws_secret_access_key"]
		return creds
	}
}

func readAwsConfig(fname string) (c *awsConfig, err error) {
	var f *os.File

	if f, err = os.Open(fname); err != nil {
		return nil, err
	}
	c = new(awsConfig)
	c.data = make(map[string]map[string]string)
	c.addSection("default")
	if err = c.reader(f); err != nil {
		return nil, err
	}
	if err = f.Close(); err != nil {
		return nil, err
	}

	return c, nil

}

func (c *awsConfig) addSection(section string) bool {
	section = strings.ToLower(section)

	if _, ok := c.data[section]; ok {
		return false
	}
	c.data[section] = make(map[string]string)

	return true
}

func (c *awsConfig) addKey(section string, key string, value string) bool {
	c.addSection(section)

	section = strings.ToLower(section)
	key = strings.ToLower(key)

	_, ok := c.data[section][key]

	c.data[section][key] = value

	return !ok
}

func (c *awsConfig) reader(f io.Reader) (err error) {
	buf := bufio.NewReader(f)

	var section, key string
	section = "default"

	for {
		line, bufErr := buf.ReadString('\n')
		line = strings.TrimSpace(line)

		if bufErr != nil {
			if bufErr != io.EOF {
				return err
			}

			if len(line) == 0 {
				break
			}
		}

		switch {
		case len(line) == 0:
			continue
		case line[0] == '#':
			continue
		case line[0] == ';':
			continue
		case line[0] == '[' && line[len(line)-1] == ']':
			key = ""
			section = strings.TrimSpace(line[1 : len(line)-1])
			c.addSection(section)
		default:
			i := strings.IndexAny(line, "=")
			key = strings.TrimSpace(line[0:i])
			value := strings.TrimSpace(line[i+1:])
			c.addKey(section, key, value)
		}

		if bufErr == io.EOF {
			break
		}
	}

	return nil
}

func (conf *config) CreateOrUpdateSecret(key iamtypes.AccessKey) {
	url := fmt.Sprintf("%s/v1/secret/data/%s", conf.aiVaultAddress, conf.secretPath)

	reqBody := make(map[string]interface{})
	reqBody["data"] = key
	reqBodyBytes, _ := json.Marshal(reqBody)

	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(reqBodyBytes))
	req.Header.Set("X-Vault-Token", conf.vaultToken)
	res, err := client.Do(req)
	if err != nil {
		log.Sugar.Fatal(err)
	}
	defer res.Body.Close()
}

func (conf *config) DeleteSecret() {
	url := fmt.Sprintf("%s/v1/secret/metadata/%s", conf.aiVaultAddress, conf.secretPath)

	client := &http.Client{}
	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Set("X-Vault-Token", conf.vaultToken)
	res, err := client.Do(req)
	if err != nil {
		log.Sugar.Fatal(err)
	}
	defer res.Body.Close()
}

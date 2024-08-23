/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// installCmd represents the install command
const (
	DZN_COMMON_FUNCTIONS = "dzn_commonFunctions"
	DZN_GEAR             = "dzn_gear"
	DZN_DYNAI            = "dzn_dynai"
	DZN_TSFRAMEWORK      = "dzn_tSFramework"

	HTTPS_PREFIX           = "https://"
	GITHUB_DOWNLOAD_SUFFIX = "archive/refs/heads/%s.zip"
	GITHUB_DEFAULT_BRANCH  = "master"
)

var (
	cfgFile  string
	defaults = map[string]string{
		DZN_COMMON_FUNCTIONS: "github.com/10Dozen/dzn_commonFunctions",
		DZN_GEAR:             "github.com/10Dozen/dzn_gear",
		DZN_DYNAI:            "github.com/10Dozen/dzn_dynai",
		DZN_TSFRAMEWORK:      "github.com/10Dozen/dzn_tSFramework",
	}

	installCmd = &cobra.Command{
		Use:   "install",
		Short: "Installs the tSFramework to given directory",
		Long: `Downloads, unpack and install required components of the tSFramewrok:
- dzn_commonFunctions
- dzn_gear
- dzn_dynai
- dzn_tSFramework
	
By default installs current version from main branch, but exact repo may be specified by flags.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("install called")
			install()

			// downloadRepo(viper.GetString("dzn_commonFunctions"))
		},
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	installCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tSFInstaller.yaml)")

	installCmd.PersistentFlags().String("dir", "", "output directory (where to install, required)")
	installCmd.PersistentFlags().StringP(DZN_COMMON_FUNCTIONS, "c", "", "repo path to dzn_commonFunctions")
	installCmd.PersistentFlags().StringP(DZN_GEAR, "g", "", "repo path to dzn_gear")
	installCmd.PersistentFlags().StringP(DZN_DYNAI, "d", "", "repo path to dzn_dynai")
	installCmd.PersistentFlags().StringP(DZN_TSFRAMEWORK, "f", "", "repo path to dzn_tSFramework")

	viper.BindPFlag(
		"dir", installCmd.PersistentFlags().Lookup("dir"))
	installCmd.MarkPersistentFlagRequired("dir")

	rootCmd.AddCommand(installCmd)
}

func initConfig() {
	viper.SetConfigFile(".tSFInstaller.yaml")

	if cfgFile != "" {
		// Use config passed as parameter
		viper.SetConfigFile(cfgFile)
	}

	// Defaults - allow configs to contain only diff
	viper.SetDefault(DZN_COMMON_FUNCTIONS, defaults[DZN_COMMON_FUNCTIONS])
	viper.SetDefault(DZN_GEAR, defaults[DZN_GEAR])
	viper.SetDefault(DZN_DYNAI, defaults[DZN_DYNAI])
	viper.SetDefault(DZN_TSFRAMEWORK, defaults[DZN_TSFRAMEWORK])

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else if cfgFile != "" {
		// Config was provided but not found
		panic(err)
	} else {
		// In case no config was provided and default config not yet exists
		viper.WriteConfig()
		println("Default config generated!")
	}
}

// Functions
func install() {
	outDir := viper.GetString("dir")
	copyFolder(
		downloadRepo(DZN_COMMON_FUNCTIONS, viper.GetString(DZN_COMMON_FUNCTIONS)),
		outDir)
	copyFolder(
		downloadRepo(DZN_GEAR, viper.GetString(DZN_GEAR)),
		outDir)
	copyFolder(
		downloadRepo(DZN_DYNAI, viper.GetString(DZN_DYNAI)),
		outDir)
	copyFolder(
		downloadRepo(DZN_TSFRAMEWORK, viper.GetString(DZN_TSFRAMEWORK)),
		outDir)
}

func downloadRepo(name, url string) string {
	fqUrl := getDownloadUrl(name, url)
	fmt.Println(fqUrl)

	err := downloadFile(name, fqUrl)
	if err != nil {
		panic(err)
	}
	dir, err := os.MkdirTemp("", "temp")
	if err != nil {
		panic(err)
	}

	err = unzipFile(name, dir)
	if err != nil {
		panic(err)
	}

	os.Remove(name)

	return dir
}

func getDownloadUrl(name, url string) string {
	// Pre-format url
	fqUrl := completeUrl(url)
	if fqUrl == completeUrl(defaults[name]) {
		fqUrl += fmt.Sprintf(GITHUB_DOWNLOAD_SUFFIX, GITHUB_DEFAULT_BRANCH)
		fmt.Println("Default repo:", fqUrl)
		return fqUrl
	}

	// Non-default repo path
	urlParts := strings.Split(fqUrl, "/")
	branchName := urlParts[6]
	return strings.Join(urlParts[0:5], "/") + "/" + fmt.Sprintf(GITHUB_DOWNLOAD_SUFFIX, branchName)
}

func completeUrl(url string) string {
	fqUrl := url
	if !strings.HasPrefix(url, HTTPS_PREFIX) {
		fqUrl = HTTPS_PREFIX + fqUrl
	}
	lastChar := string(fqUrl[len(fqUrl)-1])
	if lastChar != "/" {
		fqUrl += "/"
	}
	return fqUrl
}

func downloadFile(filepath, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func unzipFile(source, dest string) error {
	read, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer read.Close()

	for _, file := range read.File {
		if file.Mode().IsDir() {
			continue
		}
		open, err := file.Open()
		if err != nil {
			return err
		}
		name := path.Join(dest, file.Name)
		os.MkdirAll(path.Dir(name), os.ModeDir)
		create, err := os.Create(name)
		if err != nil {
			return err
		}
		defer create.Close()
		create.ReadFrom(open)
	}
	return nil
}

func copyFolder(source, dest string) {
	f, err := os.Open(source)
	if err != nil {
		log.Fatal(err)
	}
	fileInfo, err := f.ReadDir(-1)
	dirToCopy := path.Join(source, fileInfo[0].Name())
	fmt.Println(dirToCopy)

	err = os.CopyFS(dest, os.DirFS(dirToCopy))
	if err != nil {
		log.Fatal(err)
	}
}

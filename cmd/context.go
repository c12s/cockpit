package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/c12s/cockpit/cmd/helper"
	"github.com/c12s/cockpit/cmd/model"
	"github.com/c12s/cockpit/cmd/model/request"
	"github.com/spf13/cobra"
	"os"
	"os/user"
	"path/filepath"
	"time"
)

var ContextCmd = &cobra.Command{
	Use:   "context",
	Short: "Init empty CLI context environment, to interact with region/s cluster/s node/s job/s",
	Long:  "change all data inside regions, clusters and nodes",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please provide some of avalible commands or type help for help")
	},
}

func createDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	} else {
		return errors.New("Context already exists.")
	}
	return nil
}

func setContext(path, address, version, user, token string) error {
	filename := filepath.Join(path, "context.yml")
	file, err := os.Create(filename)
	defer file.Close()
	if err != nil {
		return err
	}

	var v = "v1"
	if len(version) > 0 {
		v = version
	}

	c := &model.CContext{
		Context: &model.Content{
			Version:   v,
			Address:   address,
			Namespace: "default",
			User:      user,
			Token:     token,
		},
	}

	nerr, data := model.Marshall(c)
	if nerr != nil {
		return nerr
	}
	fmt.Fprintf(file, data)

	return nil
}

func updateContext(path, address, version, namespace, user, token string) error {
	filename := filepath.Join(path, "context.yml")
	file, err := os.Create(filename)
	defer file.Close()
	if err != nil {
		return err
	}

	var v = "v1"
	if len(version) > 0 {
		v = version
	}

	c := &model.CContext{
		Context: &model.Content{
			Version:   v,
			Address:   address,
			Namespace: namespace,
			User:      user,
			Token:     token,
		},
	}

	nerr, data := model.Marshall(c)
	if nerr != nil {
		return nerr
	}
	fmt.Fprintf(file, data)

	return nil
}

func initContext(path, address, version string) error {
	filename := filepath.Join(path, "context.yml")
	file, err := os.Create(filename)
	defer file.Close()
	if err != nil {
		return err
	}

	var v = "v1"
	if len(version) > 0 {
		v = version
	}

	c := &model.CContext{
		Context: &model.Content{
			Version:   v,
			Address:   address,
			Namespace: "default",
			User:      "",
			Token:     "",
		},
	}

	nerr, data := model.Marshall(c)
	if nerr != nil {
		return nerr
	}
	fmt.Fprintf(file, data)

	return nil
}

func getContextPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	contextPath := filepath.Join(usr.HomeDir, ".constellations")
	return contextPath, err
}

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Init empty CLI context environment, to interact with region/s cluster/s node/s job/s",
	Long:  "change all data inside regions, clusters and nodes",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		a := cmd.Flag("address").Value.String()
		v := cmd.Flag("version").Value.String()
		contextPath, err := getContextPath()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Empty context initialized in %s.\n", contextPath)
		fmt.Println("Run 'cockpit context register -f your_file.yml' to register")
		fmt.Println("Run 'cockpit context login -u your_username -p your_password' to login")
		err = createDir(contextPath)
		if err != nil {
			fmt.Println(err)
			return
		}

		initContext(contextPath, a, v)
	},
}

func doLogin(username, password string) {
	crd := &request.Credentials{
		Username: username,
		Password: password,
	}

	err, ctx := helper.GetContext()
	if err != nil {
		fmt.Println(err)
		return
	}

	callPath := helper.FormCall("auth", "login", ctx, map[string]string{})
	data, err := json.Marshal(crd)
	if err != nil {
		fmt.Println(err)
		return
	}

	err, resp, token := helper.PostCallExtractToken(10*time.Second, callPath, string(data))
	if err != nil {
		fmt.Println(err)
		return
	}

	contextPath, err := getContextPath()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(resp)

	setContext(contextPath, ctx.Address(), ctx.Version(), username, token)
	fmt.Println("Context set up")
}

var LoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login user on inited context environment, to interact with region/s cluster/s node/s job/s",
	Long:  "change all data inside regions, clusters and nodes",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		u := cmd.Flag("username").Value.String()
		p := cmd.Flag("password").Value.String()

		doLogin(u, p)
	},
}

var LogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout user on inited context environment, to interact with region/s cluster/s node/s job/s",
	Long:  "change all data inside regions, clusters and nodes",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please provide some of avalible commands or type help for help")
	},
}

var RegisterCmd = &cobra.Command{
	Use:   "register",
	Short: "Register new user to the system",
	Long:  "Create new user to the platform and create default namespace for him",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		file := cmd.Flag("file").Value.String()
		if _, err := os.Stat(file); err == nil {
			f, err := mutateNFile(file)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			err, ctx := helper.GetContext()
			if err != nil {
				fmt.Println(err)
				return
			}

			err, data := users(f)
			if err != nil {
				fmt.Println(err)
				return
			}

			h := map[string]string{
				"Content-Type": "application/json; charset=UTF-8",
			}

			callPath := helper.FormCall("auth", "register", ctx, map[string]string{})
			err, resp := helper.Post(10*time.Second, callPath, data, h)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(resp)
		}
	},
}

func drop(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}

	err = os.Remove(dir)
	if err != nil {
		return err
	}

	return nil
}

var DropCmd = &cobra.Command{
	Use:   "drop",
	Short: "Drop inited context environment.",
	Long:  "change all data inside regions, clusters and nodes",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		usr, err := user.Current()
		if err != nil {
			fmt.Println(err)
		}
		contextPath := filepath.Join(usr.HomeDir, ".constellations")
		err = drop(contextPath)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Current context dropped!")
	},
}

var SwitchCmd = &cobra.Command{
	Use:   "switch",
	Short: "Switch namespace or context used",
	Long:  "Switch used context or namespace that is currently in use",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		ns := cmd.Flag("namespace").Value.String()
		c, _ := cmd.Flags().GetBool("context")

		path, err := getContextPath()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		err, ctx := helper.GetContext()
		if err != nil {
			fmt.Println(err)
			return
		}

		if c {
			fmt.Println("Switch context not implemented yet")
		} else {
			fmt.Println("Switch namespace, new one is", ns)
			err = updateContext(path, ctx.Address(), ctx.Version(), ns, ctx.User(), ctx.Token())
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}
	},
}

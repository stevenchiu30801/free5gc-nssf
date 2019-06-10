/*
 * NSSF NS Selection
 *
 * NSSF Network Slice Selection Service
 */

package factory

import (
    "fmt"
    "io/ioutil"

    "gopkg.in/yaml.v2"

    flog "../flog"
)

var NssfConfig Config

func checkErr(err error) {
    if err != nil {
        err = fmt.Errorf("[Configuration] %s", err.Error())
        flog.System.Fatal(err)
    }
}

// TODO: Support configuration update from REST api
func InitConfigFactory(f string) {
    content, err := ioutil.ReadFile(f)
    checkErr(err)

    err = yaml.Unmarshal([]byte(content), &NssfConfig)
    checkErr(err)

    err = NssfConfig.CheckIntegrity()
    checkErr(err)

    // d, err := yaml.Marshal(NssfConfig)
    // checkErr(err)
    // fmt.Printf(string(d))
}

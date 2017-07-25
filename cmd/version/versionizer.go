// author versionizer
package version
import "runtime"
import "fmt"

func VersionString() string {
	o := fmt.Sprintf("Built with %s at 2017-07-24 21:25:32.082278056 -0500 CDT at git hash 98e388944f55490db3fd12a8bacc230ab6157141 refs/heads/master\n", runtime.Version())
	return o
}

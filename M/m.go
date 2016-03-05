package MOctopus

import "sync"

var (
	versions map[string]node;
	versionsMutex sync.RWMutex;
)

type node interface {
	Id () string	//id
	Type() string	//fork or version
	Next(map[string]string) string; //map is user info(in k/v), return id
}

type judgment string

func judge(j judgment, m map[string]string) bool {

}

type fork struct {
	id string
	judgment string
	yes string;	//id if judgment
	no string;	//id
}

func (f *fork) Id() {
	return f.id;
}

func (f *fork) Next(m map[string]string) string {
	if judge(f.judgment, m) {
		return f.yes
	}
	return f.no
}

func (f *fork) Type() string {
	return "fork"
}

type version struct {
	id string
	next string
	VersionName string
}

func (v *version) Id() {
	return v.id;
}

func (v *version) Next(m map[string]string) string {
	return v.next;
}

func (v *version) Type() string {
	return "version"
}
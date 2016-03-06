package MOctopus

import "sync"

var (
	versions map[string]node
	versionsWithName map[string]*version
	lstModify int64
	versionsMutex sync.RWMutex
)

func init() {
	versions = map[string]node{}
	versionsWithName = map[string]*version{}
}

type node interface {
	Id () string	//id
	Type() string	//fork or version
	Next(map[string]string) string; //map is user info(in k/v), return id
}

type judgment string

func NewJudgment(j string) (judgment, error) {

}

func (j judgment)Judge(m map[string]string) bool {

}

type fork struct {
	id string
	judgment judgment
	yes string;	//id if judgment
	no string;	//id
}

func (f *fork) Id() {
	return f.id;
}

func (f *fork) Next(m map[string]string) string {
	if f.judgment.Judge(m) {
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
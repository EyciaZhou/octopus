package MOctopus

import (
	"github.com/nu7hatch/gouuid"
	"errors"
	"github.com/Sirupsen/logrus"
)

type operate struct {
}

var Operate operate

type KVS map[string]string

/*
只在内部调用，防止死锁，不内置lock
调用前应lock
 */
func (o *operate) getWithVersion_Recursion(n node, kvs *KVS) (*version) {
	if n == nil {
		return nil
	}

	id := n.Next(kvs)
	ans := o.getWithVersion_Recursion(versions[id], kvs)

	//两种情况
	// 1. 返回的是null， 即后面没有节点或者后面只有分叉
	//   1.1 如果自己也是分叉 返回null
	//   1.2 如果自己是version 返回自己
	// 2. 返回version

	if ans == nil {
		if n.Type() == "fork" {
			return nil
		} else {
			return n.(*version)
		}
	}
	return ans
}

//update fork
func (o *operate) newForkWithId(Id string, JudgmentString string, yesId string, noId string) {
	j, err := NewJudgment(JudgmentString)
	if err != nil {
		return err
	}

	if _, ok := versionsWithName[yesId]; yesId != "" && !ok {
		return nil, errors.New("不存在的目标节点:符合条件")
	}

	if _, ok := versionsWithName[noId]; noId != "" && !ok {
		return nil, errors.New("不存在的目标节点:不符合条件")
	}

	return &fork{Id, j, yesId, noId}, nil
}

func (o *operate) newFork(JudgmentString string, yesId string, noId string) (*fork, error) {
	uuid, err := uuid.NewV4()
	if err != nil {
		logrus.Error("无法生成uuid"+err.Error())
		return nil, errors.New("服务器错误， 无法生成uuid")
	}

	return o.newForkWithId(uuid, JudgmentString, yesId, noId)
}

func (o *operate) newVersionWithId(Id string, nextId string, VersionName string) (*version, error) {
	if _, ok := versionsWithName[nextId]; nextId != "" && !ok {
		return nil, errors.New("不存在的目标节点:符合条件")
	}
	return &version{Id, nextId, VersionName}
}

func (o *operate) newVersion(nextId string, VersionName string) (*version, error) {
	uuid, err := uuid.NewV4()
	if err != nil {
		logrus.Error("无法生成uuid"+err.Error())
		return errors.New("服务器错误， 无法生成uuid")
	}

	return o.newVersionWithId(uuid, nextId, VersionName)
}


/*
根据版本名称和 用户信息， 返回用户应更新到的版本
如果不用更新什么的 返回null
 */
func (o *operate) GetWithVersion(VersionName string, kvs *KVS) *version {
	versionsMutex.RLock()
	defer versionsMutex.RUnlock()

	if _, ok := versionsWithName[VersionName]; !ok {
		return nil
	}

	return o.getWithVersion_Recursion(versionsWithName[VersionName], kvs)
}

func (o *operate) AddFork(JudgmentString string, yesId string, noId string) (error) {
	versionsMutex.RLock()
	defer versionsMutex.RUnlock()
}

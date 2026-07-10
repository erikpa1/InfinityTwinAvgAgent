package interfaces

type UidProvider interface {
	GetUid() string
	SetUid(uid string)
}

type IdProvider interface {
	GetId() int64
	SetId(uid int64)
}

type UidMap map[string]interface{}

func (self UidMap) GetUid() string {
	tmp, _ := self["uid"].(string)
	return tmp
}

func (self UidMap) SetUid(new_uid string) {
	//Tuna som musel vymazat pointer z UidMapy boh vie preco
	self["uid"] = new_uid
}

type TypeProvider interface {
	GetType() string
}

type StepConsumer interface {
	Step()
}

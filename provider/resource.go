package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ResourceData struct {
	*schema.ResourceData
	Schema map[string]*schema.Schema
}

func (rd *ResourceData) Get(key string) interface{} {
	var err error
	vt, ok := rd.Schema[key]
	if !ok {
		return nil
	}

	v := rd.ResourceData.Get(key)
	switch vt.Type {
	case schema.TypeSet:
		if vt.DefaultFunc != nil && v.(*schema.Set).Len() == 0 {
			v, err = vt.DefaultFunc()
			if err != nil {
				panic(err)
			}
			// NOTE: this returns Set to mitigate error
			// Attribute must be a list
			// also this makes default value consistent
			// with user-defined values for sets
			return schema.NewSet(HashAny, v.([]interface{}))
		}
	case schema.TypeList:
		if vt.DefaultFunc != nil && len(v.([]interface{})) == 0 {
			v, err = vt.DefaultFunc()
			if err != nil {
				panic(err)
			}
			return v
		}
	case schema.TypeMap:
		if vt.DefaultFunc != nil && len(v.(map[string]interface{})) == 0 {
			v, err = vt.DefaultFunc()
			if err != nil {
				panic(err)
			}
			return v
		}
	}
	return v
}

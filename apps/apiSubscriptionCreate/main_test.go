package main

import (
	"fmt"
	"testing"
)

func Test_checkEvent(t *testing.T) {
	var event string
	e := &Event{}
	//
	event = `{"Name":"foo", "Age":20}`
	fmt.Println(checkEvent([]byte(event), e))
	//
	event = `{"custom_id":"1"}`
	fmt.Println(checkEvent([]byte(event), e))
	//
	event = `{"custom_id":"1", "plan_id":"123"}`
	fmt.Println(checkEvent([]byte(event), e))
	//
	event = `{"custom_id":"1", "plan_id":"123", "give_name":"developer"}`
	fmt.Println(checkEvent([]byte(event), e))
	//
	event = `{"custom_id":"1", "plan_id":"123", "give_name":"developer", "surname":"astra"}`
	fmt.Println(checkEvent([]byte(event), e))
	//
	event = `{"custom_id":"1", "plan_id":"123", "give_name":"developer", "surname":"astra", "email_address":"astra@email.com"}`
	fmt.Println(checkEvent([]byte(event), e))
	fmt.Printf("%#v\n", *e.Surname)
}

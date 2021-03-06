package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
        _ "github.com/hashicorp/terraform/terraform"
	"log"
	"fmt"
	_ "errors"
	"encoding/json"
)

func resourceKibanaSavedObject() *schema.Resource {
	return &schema.Resource{
		Create: resourceElasticSavedObjectCreate,
		Read: resourceElasticSavedObjectRead,
		Update: resourceElasticSavedObjectUpdate,
		Delete: resourceElasticSavedObjectDelete,
		Schema: map[string]*schema.Schema{
			"saved_object_type": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{ "index-pattern", "visualization", "search", "timelion-sheet",}, false),
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"attributes": &schema.Schema{
				Type:             schema.TypeString,
				Required:         true,
			},
		},
	}
}

func resourceElasticSavedObjectCreate(d *schema.ResourceData, meta interface{}) error {
	url := meta.(*ElasticInfo).kibanaUrl
	_ = d.Get("name").(string)

	attributes := d.Get("attributes").(string)
	saved_object_type := d.Get("saved_object_type").(string)

	url = fmt.Sprintf("%v/api/saved_objects/%v", url, saved_object_type)

	var savedObjectHeader SavedObjectHeader
	json.Unmarshal([]byte(attributes), &savedObjectHeader.Attributes)
	body, err := json.Marshal(&savedObjectHeader)
	if err != nil {
		return err
	}

	respBody, err := postKibRequest(d, meta, url, string(body))
	if err != nil {
		return err
	}

	var savedObject SavedObjectHeader
	json.Unmarshal(*respBody, &savedObject)	
	log.Printf("Raw Body: %s", respBody)
	log.Printf("ID: %s", savedObject.Id)
	log.Printf("UpdatedAt: %s", savedObject.UpdatedAt)
	log.Printf("Version: %v", savedObject.Version)

	d.SetId(savedObject.Id)
	d.Set("version", savedObject.Version)

	return err
}

func resourceElasticSavedObjectRead(d *schema.ResourceData, meta interface{}) error {
	url := meta.(*ElasticInfo).kibanaUrl
	id := d.Id()
	saved_object_type := d.Get("saved_object_type").(string)

	url = fmt.Sprintf("%v/api/saved_objects/%v/%v", url, saved_object_type, id)
	respBody, err := getKibRequest(d, meta, url)
	if err != nil {
       	    return err
	}

	var savedObject SavedObjectHeader
	json.Unmarshal(*respBody, &savedObject)
	return nil
}

func resourceElasticSavedObjectUpdate(d *schema.ResourceData, meta interface{}) error {
	url := meta.(*ElasticInfo).kibanaUrl
	id := d.Id()
	saved_object_type := d.Get("saved_object_type").(string)

	attributes := d.Get("attributes").(string)
	var savedObjectHeader SavedObjectHeader
	json.Unmarshal([]byte(attributes), &savedObjectHeader.Attributes)
	body, err := json.Marshal(&savedObjectHeader)
	if err != nil {
		return err
	}

	url = fmt.Sprintf("%v/api/saved_objects/%v/%v", url, saved_object_type, id)
	respBody, err := putKibRequest(d, meta, url, string(body))
	if err != nil {
		return err
	}

	var savedObject SavedObjectHeader
	json.Unmarshal(*respBody, &savedObject)	
	log.Printf("Raw Body: %s", respBody)
	log.Printf("ID: %s", savedObject.Id)
	log.Printf("UpdatedAt: %s", savedObject.UpdatedAt)
	log.Printf("Version: %v", savedObject.Version)
	d.Set("version", savedObject.Version)
	return nil
}

func resourceElasticSavedObjectDelete(d *schema.ResourceData, meta interface{}) error {
	url := meta.(*ElasticInfo).kibanaUrl
	id := d.Id()
	saved_object_type := d.Get("saved_object_type").(string)

	url = fmt.Sprintf("%v/api/saved_objects/%v/%v", url, saved_object_type, id)
	_, err := deleteKibRequest(d, meta, url)	
	if err != nil {
       		return err    
	}
	d.Set("version", nil)
	d.SetId("")
	return nil
}



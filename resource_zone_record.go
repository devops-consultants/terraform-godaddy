package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

type zoneRecordResource struct {
	Customer string
	Domain   string
	Records  []*DomainRecord
}

func newZoneRecordResource(d *schema.ResourceData) (zoneRecordResource, error) {
	var err error
	r := zoneRecordResource{}

	if attr, ok := d.GetOk("customer"); ok {
		r.Customer = attr.(string)
	}

	if attr, ok := d.GetOk("domain"); ok {
		r.Domain = attr.(string)
	}

	if attr, ok := d.GetOk("record"); ok {
		records := attr.(*schema.Set).List()
		r.Records = make([]*DomainRecord, len(records))

		for i, rec := range records {
			data := rec.(map[string]interface{})
			r.Records[i], err = NewDomainRecord(
				data["name"].(string),
				data["type"].(string),
				data["data"].(string),
				data["ttl"].(int))

			if err != nil {
				return r, err
			}
		}
	}

	return r, err
}

func resourceZoneRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceZoneRecordUpdate,
		Read:   resourceZoneRecordRead,
		Update: resourceZoneRecordUpdate,
		Delete: resourceZoneRecordDelete,

		Schema: map[string]*schema.Schema{
			// Required
			"domain": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			// Optional
			"customer": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"record": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"data": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"ttl": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Default:  DefaultTTL,
						},
					},
				},
			},
		},
	}
}

func resourceZoneRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*GoDaddyClient)
	customer := d.Get("customer").(string)
	domain := d.Get("domain").(string)
	rectype := "A"
	recname := ""

	if attr, ok := d.GetOk("record"); ok {
		records := attr.(*schema.Set).List()
		// r.Records = make([]*DomainRecord, len(records))

		for i, rec := range records {
			data := rec.(map[string]interface{})
			rectype := data["type"].(string)
			recname := data["name"].(string)
		}
	}

	log.Println("Fetching", domain, "record for ", rectype, recname, "...")
	records, err := client.GetZoneRecords(customer, domain, rectype, recname)
	if err != nil {
		return fmt.Errorf("couldn't find zone record (%s %s %s): %s", domain, rectype, recname, err.Error())
	}

	return populateResourceDataFromResponse(records, d)
}

func resourceZoneRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	// client := meta.(*GoDaddyClient)
	// r, err := newZoneRecordResource(d)
	// if err != nil {
	// 	return err
	// }

	// if err = populateDomainInfo(client, &r, d); err != nil {
	// 	return err
	// }

	// log.Println("Updating", r.Domain, "domain records...")
	// r.converge()
	// return client.UpdateZoneRecords(r.Customer, r.Domain, r.Records)
	return ""
}

func resourceZoneRecordDelete(d *schema.ResourceData, meta interface{}) error {
	// client := meta.(*GoDaddyClient)
	// customer := d.Get("customer").(string)
	// domain := d.Get("domain").(string)
	// rectype := "A"
	// recname := ""

	// if attr, ok := d.GetOk("record"); ok {
	// 	records := attr.(*schema.Set).List()
	// 	// r.Records = make([]*DomainRecord, len(records))

	// 	for i, rec := range records {
	// 		data := rec.(map[string]interface{})
	// 		rectype := data["type"].(string)
	// 		recname := data["name"].(string)
	// 	}
	// }

	// log.Println("Deleting", domain, "record for ", rectype, recname, "...")
	// records, err := client.GetZoneRecords(customer, domain, rectype, recname)
	// if err != nil {
	// 	return fmt.Errorf("couldn't find zone record (%s %s %s): %s", domain, rectype, recname, err.Error())
	// }

	// return populateResourceDataFromResponse(records, d)
	return ""
}

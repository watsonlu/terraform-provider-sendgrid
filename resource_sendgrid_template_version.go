package sendgrid

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mitchellh/go-homedir"
	"github.com/romanlaguta/terraform-client-sendgrid"
)

func resourceSendgridTemplateVersion() *schema.Resource {
	return &schema.Resource{
		Create: resourceSendgridTemplateVersionCreate,
		Read:   resourceSendgridTemplateVersionRead,
		Update: resourceSendgridTemplateVersionUpdate,
		Delete: resourceSendgridTemplateVersionDelete,
		Exists: resourceSendgridTemplateVersionExists,
		Importer: &schema.ResourceImporter{
			State: resourceSendgridTemplateVersionImport,
		},

		Schema: map[string]*schema.Schema{
			"template_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"subject": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"html_content_file": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"plain_content_file": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"html_content": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"plain_content": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"active": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func buildTemplateVersionStruct(d *schema.ResourceData) (*sendgrid_client.TemplateVersion, error) {
	htmlContent, err := loadFileContent(d.Get("html_content_file").(string))
	if err != nil {
		return nil, err
	}

	plainContent, err := loadFileContent(d.Get("plain_content_file").(string))
	if err != nil {
		return nil, err
	}

	active := 0
	if d.Get("active").(bool) {
		active = 1
	}

	m := sendgrid_client.TemplateVersion{
		TemplateId:   d.Get("template_id").(string),
		Name:         d.Get("name").(string),
		Subject:      d.Get("subject").(string),
		HtmlContent:  string(htmlContent),
		PlainContent: string(plainContent),
		Active:       active,
	}

	return &m, nil
}

func resourceSendgridTemplateVersionExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	// Exists - This is called to verify a resource still exists. It is called prior to Read,
	// and lowers the burden of Read to be able to assume the resource exists.
	client := meta.(*sendgrid_client.Client)

	fmt.Println("Exist template_version")
	if _, err := client.GetTemplateVersion(d.Get("template_id").(string), d.Id()); err != nil {
		if strings.Contains(err.Error(), "404 Not Found") {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func resourceSendgridTemplateVersionCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*sendgrid_client.Client)

	m, err := buildTemplateVersionStruct(d)
	if err != nil {
		return fmt.Errorf("error updating template_version: %s", err.Error())
	}
	fmt.Println("Create template_version1")
	m, err = client.CreateTemplateVersion(m)
	if err != nil {
		return fmt.Errorf("error updating template_version: %s", err.Error())
	}
	fmt.Println("Create template_version2")
	d.SetId(m.Id)
	d.Set("html_content", m.HtmlContent)
	d.Set("plain_content", m.PlainContent)

	return nil
}

func resourceSendgridTemplateVersionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*sendgrid_client.Client)

	fmt.Println("Read template_version")
	m, err := client.GetTemplateVersion(d.Get("template_id").(string), d.Id())
	if err != nil {
		return err
	}
	fmt.Println("[DEBUG] TemplateVersion: %v", m)
	d.Set("name", m.Name)
	d.Set("template_id", m.TemplateId)
	d.Set("subject", m.Subject)
	d.Set("html_content", m.HtmlContent)
	d.Set("plain_content", m.PlainContent)
	d.Set("active", m.Active)

	return nil
}

func resourceSendgridTemplateVersionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*sendgrid_client.Client)

	m, err := buildTemplateVersionStruct(d)
	if err != nil {
		return err
	}
	fmt.Println("Update template_version")
	if err := client.UpdateTemplateVersion(d.Id(), m); err != nil {
		return fmt.Errorf("error updating TemplateVersion: %s", err.Error())
	}

	return resourceSendgridTemplateVersionRead(d, meta)
}

func resourceSendgridTemplateVersionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*sendgrid_client.Client)

	fmt.Println("Delete template_version")
	if err := client.DeleteTemplateVersion(d.Get("template_id").(string), d.Id()); err != nil {
		return err
	}

	return nil
}

func resourceSendgridTemplateVersionImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	fmt.Println("Import template_version")
	if err := resourceSendgridTemplateVersionRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

// loadFileContent returns contents of a file in a given path
func loadFileContent(v string) ([]byte, error) {
	filename, err := homedir.Expand(v)
	if err != nil {
		fmt.Println("File %s can't be expand. %s", v, err)
		return nil, err
	}
	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("File %s can't be read. %s", filename, err)
		return nil, err
	}
	return fileContent, nil
}
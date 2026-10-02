package global_conf

import (
	"context"
	"fmt"
	"go-pkg/helper/context_helper"
	"go-pkg/helper/error_helper"
	"go-pkg/log/base"
	m "go-pkg/mongo"
	service_pkg "go-pkg/service"

	"go.mongodb.org/mongo-driver/bson"
)

const templateCollection = "template"

type Template struct {
	ID              string                 `bson:"_id"`
	Html            string                 `bson:"html"`
	HtmlByPlatform  map[string]string      `bson:"htmlByPlatform"`
	Params          map[string]interface{} `bson:"params"`
	HtmlParams      map[string]string      `bson:"htmlParams"`
	Domains         []string               `bson:"domains"`
	HasClientParams bool                   `bson:"hasClientParams"`
}

func (service *Service) GetTemplate(ctx context.Context, domain string) (*Template, error) {
	err := service.subscribe(templateCollection, func(ctx context.Context, data *m.WatchEvent) {
		service.removeTemplate()
	})
	if err != nil {
		context_helper.Logger(ctx).Error(func() string {
			return fmt.Sprintf("subscribe template error %v", err)
		})
	}

	if template, ok := service.getTemplate(domain); ok {
		return template.Clone(), nil
	}

	//don't lock here so that warm up requests not waiting for each others
	templates, err := service.getTemplates(ctx)
	if err != nil {
		return nil, err
	}

	service.assignTemplate(templates)

	template, ok := service.getTemplate(domain)
	if !ok {
		return nil, error_helper.New(error_helper.NotFound, service_pkg.Domain, service_pkg.Name, "not found template", nil, "")
	}
	return template.Clone(), nil
}

func (service *Service) getTemplate(domain string) (*Template, bool) {
	service.templateLock.RLock()
	defer service.templateLock.RUnlock()

	template, found := service.templates[domain]
	return template, found
}

func (service *Service) getTemplates(ctx context.Context) (map[string]*Template, error) {
	db, err := m.Client()
	if err != nil {
		context_helper.Logger(ctx).Error(func() string {
			return fmt.Sprintf("get global db error %v", err)
		})
		return nil, err
	}

	cursor, err := db.Database("app_global").Collection(templateCollection).Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var templates []*Template
	err = cursor.All(ctx, &templates)
	if err != nil {
		context_helper.Logger(ctx).Error(func() string {
			return "fetch template error"
		}, base.Error(err))
		return nil, err
	}

	mTemplate := make(map[string]*Template, len(templates))
	for _, template := range templates {
		for _, domain := range template.Domains {
			mTemplate[domain] = template
		}
	}

	return mTemplate, nil
}

func (service *Service) assignTemplate(templates map[string]*Template) {
	service.templateLock.Lock()
	defer service.templateLock.Unlock()

	service.templates = templates
}

func (service *Service) removeTemplate() {
	service.templateLock.Lock()
	defer service.templateLock.Unlock()

	service.templates = make(map[string]*Template)

}

func (t *Template) Clone() *Template {
	return &Template{
		ID:              t.ID,
		Html:            t.Html,
		HtmlByPlatform:  t.HtmlByPlatform,
		Params:          t.Params,
		HtmlParams:      t.HtmlParams,
		Domains:         t.Domains,
		HasClientParams: t.HasClientParams,
	}
}

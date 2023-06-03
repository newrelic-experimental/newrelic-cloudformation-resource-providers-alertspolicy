package resource

import (
   "fmt"
   "github.com/newrelic/newrelic-cloudformation-resource-providers-common/model"
   log "github.com/sirupsen/logrus"
)

//
// Generic, should be able to leave these as-is
//

type Payload struct {
   model  *Model
   models []interface{}
}

func NewPayload(m *Model) *Payload {
   return &Payload{
      model:  m,
      models: make([]interface{}, 0),
   }
}

func (p *Payload) GetResourceModel() interface{} {
   return p.model
}

func (p *Payload) GetResourceModels() []interface{} {
   log.Debugf("GetResourceModels: returning %+v", p.models)
   return p.models
}

func (p *Payload) AppendToResourceModels(m model.Model) {
   p.models = append(p.models, m.GetResourceModel())
}
func (p *Payload) GetTags() map[string]string {
   return nil
}

func (p *Payload) HasTags() bool {
   return false
}

//
// These are API specific, must be configured per API
//

var typeName = "NewRelic::Observability::AlertsPolicy"

func (p *Payload) NewModelFromGuid(g interface{}) (m model.Model) {
   s := fmt.Sprintf("%s", g)
   return NewPayload(&Model{Id: &s})
}

func (p *Payload) SetIdentifier(g *string) {
   p.model.Id = g
   log.Debugf("SetIdentifier: %s", *p.model.Id)
}

func (p *Payload) GetTagIdentifier() *string {
   return p.model.Id
}

func (p *Payload) GetIdentifier() *string {
   return p.model.Id
}

func (p *Payload) GetIdentifierKey(a model.Action) string {
   return "id"
}

var emptyString = "  "

func (p *Payload) GetGraphQLFragment() *string {
   return &emptyString
}

func (p *Payload) GetVariables() map[string]string {
   vars := make(map[string]string)
   if p.model.Variables != nil {
      for k, v := range p.model.Variables {
         vars[k] = v
      }
   }

   if p.model.Id != nil {
      vars["ID"] = *p.model.Id
   }

   if p.model.IncidentPreference != nil {
      vars["INCIDENTPREFERENCE"] = *p.model.IncidentPreference
   }

   if p.model.Name != nil {
      vars["NAME"] = *p.model.Name
   }

   lqf := ""
   if p.model.ListQueryFilter != nil {
      lqf = *p.model.ListQueryFilter
   }
   vars["LISTQUERYFILTER"] = lqf

   return vars
}

func (p *Payload) GetErrorKey() string {
   return ""
}

func (p *Payload) GetCreateMutation() string {
   return `
mutation {
  alertsPolicyCreate(accountId: {{{ACCOUNTID}}}, policy: {incidentPreference: {{{INCIDENTPREFERENCE}}}, name: "{{{NAME}}}"}) {
    id
  }
}
`
}

func (p *Payload) GetDeleteMutation() string {
   return `
mutation {
  alertsPolicyDelete(accountId: {{{ACCOUNTID}}}, id: "{{{ID}}}") {
    id
  }
}
`
}

func (p *Payload) GetUpdateMutation() string {
   return `
mutation {
  alertsPolicyUpdate(accountId: {{{ACCOUNTID}}}, id: "{{{ID}}}" policy: {incidentPreference: {{{INCIDENTPREFERENCE}}}, name: "{{{NAME}}}"}) {
    id
  }
}
`
}

func (p *Payload) GetReadQuery() string {
   return `
{
  actor {
    account(id: {{{ACCOUNTID}}}) {
      alerts {
        policy(id: "{{{ID}}}") {
          id
          name
        }
      }
    }
  }
}
`
}

func (p *Payload) GetListQuery() string {
   return `
{
  actor {
    account(id: {{{ACCOUNTID}}}) {
      alerts {
        policiesSearch(cursor: "{{{NEXTCURSOR}}}") {
          nextCursor
          policies {
            id
            name
            incidentPreference
          }
        }
      }
    }
  }
}
`
}

func (p *Payload) GetListQueryNextCursor() string {
   return `
{
  actor {
    account(id: {{{ACCOUNTID}}}) {
      alerts {
        policiesSearch(cursor: "{{{NEXTCURSOR}}}") {
          nextCursor
          policies {
            id
            name
            incidentPreference
          }
        }
      }
    }
  }
}
`
}

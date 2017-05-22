package utils

import (
	"errors"
	"time"

	"github.com/DarkMetrix/monitor/server/src/protocol"

	//log "github.com/cihub/seelog"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoDBUtils struct {
	session  *mgo.Session  //Mongo db session
	database *mgo.Database //Mongo db database
}

//Init mongodb
func (db *MongoDBUtils) Init(address string, database string) error {
	var err error

	if len(address) == 0 || len(database) == 0 {
		return errors.New("Param error, address or database empty")
	}

	db.session, err = mgo.DialWithTimeout(address, time.Second*5)

	if err != nil {
		return err
	}

	db.session.SetSyncTimeout(time.Second * 5)
	db.session.SetSocketTimeout(time.Second * 5)

	db.database = db.session.DB(database)

	return nil
}

//******************Node*****************//
//Add node
func (db *MongoDBUtils) AddNode(node protocol.Node) error {
	collection := db.database.C("nodes")

	if len(node.NodeIP) == 0 {
		return errors.New("Param Invalid! node ip is empty!")
	}

	_, err := collection.Upsert(bson.M{"node_ip": node.NodeIP}, node)

	if err != nil {
		db.session.Refresh()
		return err
	}

	return nil
}

//Get node
func (db *MongoDBUtils) GetNode(ip string) ([]protocol.Node, error) {
	collection := db.database.C("nodes")

	var nodes []protocol.Node

	if ip == "all" {
		err := collection.Find(bson.M{}).All(&nodes)

		if err != nil {
			return nil, err
		}

	} else {
		err := collection.Find(bson.M{"node_ip": ip}).All(&nodes)

		if err != nil {
			return nil, err
		}
	}

	return nodes, nil
}

//Del node
func (db *MongoDBUtils) DelNode(ip string) error {
	collection := db.database.C("nodes")

	if ip == "all" {
		_, err := collection.RemoveAll(bson.M{})

		if err != nil {
			return err
		}
	} else {
		err := collection.Remove(bson.M{"node_ip": ip})

		if err != nil {
			return err
		}
	}

	return nil
}

//******************Application instance*****************//
//Add application
func (db *MongoDBUtils) AddApplicationInstance(instance protocol.ApplicationInstance) error {
	collection := db.database.C("applications")

	if len(instance.Name) == 0 {
		return errors.New("Param Invalid! node ip is empty!")
	}

	_, err := collection.Upsert(bson.M{"name": instance.Name}, instance)

	if err != nil {
		db.session.Refresh()
		return err
	}

	return nil
}

//Get application instance
func (db *MongoDBUtils) GetApplicationInstance(instance string) ([]protocol.ApplicationInstance, error) {
	collection := db.database.C("applications")

	var instances []protocol.ApplicationInstance

	if instance == "all" {
		err := collection.Find(bson.M{}).All(&instances)

		if err != nil {
			return nil, err
		}

	} else {
		err := collection.Find(bson.M{"name": instance}).All(&instances)

		if err != nil {
			return nil, err
		}
	}

	return instances, nil
}

//Del application instance
func (db *MongoDBUtils) DelApplicationInstance(instance string) error {
	collection := db.database.C("applications")

	if instance == "all" {
		_, err := collection.RemoveAll(bson.M{})

		if err != nil {
			return err
		}
	} else {
		err := collection.Remove(bson.M{"name": instance})

		if err != nil {
			return err
		}
	}

	return nil
}

//******************Application instance & node mapping*****************//
//Add application instance & node mapping
func (db *MongoDBUtils) AddApplicationInstanceNodeMapping(mapping protocol.ApplicationInstanceNodeMapping) error {
	collection := db.database.C("mapping")

	if len(mapping.NodeIP) == 0 || len(mapping.Instance) == 0 {
		return errors.New("Param Invalid! ip or instance is empty")
	}

	mapping.Key = mapping.Instance + "@" + mapping.NodeIP

	_, err := collection.Upsert(bson.M{"key": mapping.Key}, mapping)

	if err != nil {
		db.session.Refresh()
		return err
	}

	return nil
}

//Get application instance by ip
func (db *MongoDBUtils) GetApplicationInstanceNodeMappingByIP(ip string) ([]protocol.ApplicationInstanceNodeMapping, error) {
	collection := db.database.C("mapping")

	if len(ip) == 0 {
		return nil, errors.New("Param Invalid! ip is empty")
	}

	var mapping []protocol.ApplicationInstanceNodeMapping

	err := collection.Find(bson.M{"node_ip": ip}).All(&mapping)

	if err != nil {
		db.session.Refresh()
		return nil, err
	}

	return mapping, nil
}

//Get application instance by instance
func (db *MongoDBUtils) GetApplicationInstanceNodeMappingByInstance(instance string) ([]protocol.ApplicationInstanceNodeMapping, error) {
	collection := db.database.C("mapping")

	if len(instance) == 0 {
		return nil, errors.New("Param Invalid! instance is empty")
	}

	var mapping []protocol.ApplicationInstanceNodeMapping

	err := collection.Find(bson.M{"instance": instance}).All(&mapping)

	if err != nil {
		db.session.Refresh()
		return nil, err
	}

	return mapping, nil
}

//******************View*****************//
//Add project
func (db *MongoDBUtils) AddProject(project protocol.Project) error {
	collection := db.database.C("project")

	if len(project.Name) == 0 || project.Name == "all" {
		return errors.New("Param Invalid! name is empty or is 'all'")
	}

	_, err := collection.Upsert(bson.M{"project": project.Name}, project)

	if err != nil {
		db.session.Refresh()
		return err
	}

	return nil
}

//Get project
func (db *MongoDBUtils) GetProject(project string) ([]protocol.Project, error) {
	collection := db.database.C("project")

	if len(project) == 0 {
		return nil, errors.New("Param Invalid! project name is empty")
	}

	var projects []protocol.Project

	if project == "all" {
		err := collection.Find(bson.M{}).All(&projects)

		if err != nil {
			return nil, err
		}
	} else {
		err := collection.Find(bson.M{"project": project}).All(&projects)

		if err != nil {
			return nil, err
		}
	}

	return projects, nil
}

//Set project
func (db *MongoDBUtils) SetProject(project protocol.Project) error {
	collection := db.database.C("project")

	if len(project.Name) == 0 || project.Name == "all" {
		return errors.New("Param Invalid! name is empty or is 'all'")
	}

	err := collection.Update(bson.M{"project": project.Name}, project)

	if err != nil {
		db.session.Refresh()
		return err
	}

	return nil
}

//Del project
func (db *MongoDBUtils) DelProject(project string) error {
	collection := db.database.C("project")

	if project == "all" {
		_, err := collection.RemoveAll(bson.M{})

		if err != nil {
			return err
		}
	} else {
		err := collection.Remove(bson.M{"project": project})

		if err != nil {
			return err
		}
	}

	return nil
}

//Add service
func (db *MongoDBUtils) AddService(project string, service protocol.Service) error {
	collection := db.database.C("service")

	if len(project) == 0 || len(service.Name) == 0 || project == "all" || service.Name == "all" {
		return errors.New("Param Invalid! project name or service name is empty or is 'all'")
	}

	service.Project = project

	projects, err := db.GetProject(project)

	if err != nil {
		return err
	}

	if len(projects) == 0 {
		return errors.New("project not found!")
	}

	_, err = collection.Upsert(bson.M{"project": project, "service":service.Name}, service)

	if err != nil {
		db.session.Refresh()
		return err
	}

	return nil
}

//Get service
func (db *MongoDBUtils) GetService(project string, service string) ([]protocol.Service, error) {
	collection := db.database.C("service")

	if len(project) == 0 || len(service) == 0 {
		return nil, errors.New("Param Invalid! project name or service name is empty")
	}

	var services []protocol.Service

	if project == "all" {
		if service == "all" {
			err := collection.Find(bson.M{}).All(&services)

			if err != nil {
				return nil, err
			}
		} else {
			err := collection.Find(bson.M{"service": service}).All(&services)

			if err != nil {
				return nil, err
			}
		}
	} else {
		if service == "all" {
			err := collection.Find(bson.M{"project": project}).All(&services)

			if err != nil {
				return nil, err
			}
		} else {
			err := collection.Find(bson.M{"project": project, "service": service}).All(&services)

			if err != nil {
				return nil, err
			}
		}
	}

	return services, nil
}

//Set service
func (db *MongoDBUtils) SetService(project string, service protocol.Service) error {
	collection := db.database.C("service")

	if len(project) == 0 || len(service.Name) == 0 || project == "all" || service.Name == "all" {
		return errors.New("Param Invalid! project name or service name is empty or is all")
	}

	service.Project = project

	err := collection.Update(bson.M{"project": project, "service":service.Name}, service)

	if err != nil {
		db.session.Refresh()
		return err
	}

	return nil
}

//Del service
func (db *MongoDBUtils) DelService(project string, service string) error {
	collection := db.database.C("service")

	if project == "all" {
		if service == "all" {
			_, err := collection.RemoveAll(bson.M{})

			if err != nil {
				return err
			}
		} else {
			_, err := collection.RemoveAll(bson.M{"service": service})

			if err != nil {
				return err
			}
		}
	} else {
		if service == "all" {
			_, err := collection.RemoveAll(bson.M{"project": project})

			if err != nil {
				return err
			}
		} else {
			_, err := collection.RemoveAll(bson.M{"project": project, "service": service})

			if err != nil {
				return err
			}
		}
	}

	return nil
}

//Add module
func (db *MongoDBUtils) AddModule(project string, service string, module protocol.Module) error {
	collection := db.database.C("module")

	if len(project) == 0 || len(service) == 0 || len(module.Name) == 0 || project == "all" || service == "all" || module.Name == "all" {
		return errors.New("Param Invalid! project name or service name or module name is empty or is 'all'")
	}

	module.Project = project
	module.Service = service

	services, err := db.GetService(project, service)

	if err != nil {
		return err
	}

	if len(services) == 0 {
		return errors.New("project or service not found!")
	}

	_, err = collection.Upsert(bson.M{"project": project, "service": service, "module":module.Name}, module)

	if err != nil {
		db.session.Refresh()
		return err
	}

	return nil
}

//Get module
func (db *MongoDBUtils) GetModule(project string, service string, module string) ([]protocol.Module, error) {
	collection := db.database.C("module")

	if len(project) == 0 || len(service) == 0 || len(module) == 0 {
		return nil, errors.New("Param Invalid! project name or service name or module name is empty")
	}

	var modules []protocol.Module

	if project == "all" {
		if service == "all" {
			if module == "all" {
				err := collection.Find(bson.M{}).All(&modules)

				if err != nil {
					return nil, err
				}
			} else {
				err := collection.Find(bson.M{"module": module}).All(&modules)

				if err != nil {
					return nil, err
				}
			}
		} else {
			if module == "all" {
				err := collection.Find(bson.M{"service": service}).All(&modules)

				if err != nil {
					return nil, err
				}
			} else {
				err := collection.Find(bson.M{"service": service, "module": module}).All(&modules)

				if err != nil {
					return nil, err
				}
			}
		}
	} else {
		if service == "all" {
			if module == "all" {
				err := collection.Find(bson.M{"project": project}).All(&modules)

				if err != nil {
					return nil, err
				}
			} else {
				err := collection.Find(bson.M{"project": service, "module": module}).All(&modules)

				if err != nil {
					return nil, err
				}
			}
		} else {
			if module == "all" {
				err := collection.Find(bson.M{"project": project, "service": service}).All(&modules)

				if err != nil {
					return nil, err
				}
			} else {
				err := collection.Find(bson.M{"project": project, "service": service, "module": module}).All(&modules)

				if err != nil {
					return nil, err
				}
			}
		}
	}

	return modules, nil
}

//Set module
func (db *MongoDBUtils) SetModule(project string, service string, module protocol.Module) error {
	collection := db.database.C("module")

	if len(project) == 0 || len(service) == 0 || len(module.Name) == 0 || project == "all" || service == "all" || module.Name == "all" {
		return errors.New("Param Invalid! project name or service name or module name is empty or is 'all'")
	}

	module.Project = project
	module.Service = service

	err := collection.Update(bson.M{"project": project, "service": service, "module":module.Name}, module)

	if err != nil {
		db.session.Refresh()
		return err
	}

	return nil
}

//Del module
func (db *MongoDBUtils) DelModule(project string, service string, module string) error {
	collection := db.database.C("module")

	if project == "all" {
		if service == "all" {
			if module == "all" {
				_, err := collection.RemoveAll(bson.M{})

				if err != nil {
					return err
				}
			} else {
				_, err := collection.RemoveAll(bson.M{"module": module})

				if err != nil {
					return err
				}
			}
		} else {
			if module == "all" {
				_, err := collection.RemoveAll(bson.M{"service": service})

				if err != nil {
					return err
				}
			} else {
				_, err := collection.RemoveAll(bson.M{"service": service, "module": module})

				if err != nil {
					return err
				}
			}
		}
	} else {
		if service == "all" {
			if module == "all" {
				_, err := collection.RemoveAll(bson.M{"project": project})

				if err != nil {
					return err
				}
			} else {
				_, err := collection.RemoveAll(bson.M{"project": project, "module": module})

				if err != nil {
					return err
				}
			}
		} else {
			if module == "all" {
				_, err := collection.RemoveAll(bson.M{"project": project, "service": service})

				if err != nil {
					return err
				}
			} else {
				_, err := collection.RemoveAll(bson.M{"project": project, "service": service, "module": module})

				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

//Add application instance to module
func (db *MongoDBUtils) AddApplicationInstancesToModule(project string, service string, module string, instances []string) error {
	collection := db.database.C("module")

	if len(project) == 0 || len(service) == 0 || len(module) == 0 || project == "all" || service == "all" || module == "all" {
		return errors.New("Param Invalid! project name or service name or module name is empty or is 'all'")
	}

	if len(instances) == 0 {
		return errors.New("Param Invalid! instances empty")
	}

	modules, err := db.GetModule(project, service, module)

	if err != nil {
		return err
	}

	if len(modules) == 0 {
		return errors.New("project or service or module not found!")
	}

	err = collection.Update(bson.M{"project": project, "service": service, "module":module}, bson.M{"$addToSet":bson.M{"instances":bson.M{"$each":instances}}})

	if err != nil {
		db.session.Refresh()
		return err
	}

	return nil
}

//Del application instance from module
func (db *MongoDBUtils) DelApplicationInstancesFromModule(project string, service string, module string, instances []string) error {
	collection := db.database.C("module")

	if len(project) == 0 || len(service) == 0 || len(module) == 0 || project == "all" || service == "all" || module == "all" {
		return errors.New("Param Invalid! project name or service name or module name is empty or is 'all'")
	}

	if len(instances) == 0 {
		return errors.New("Param Invalid! instances empty")
	}

	modules, err := db.GetModule(project, service, module)

	if err != nil {
		return err
	}

	if len(modules) == 0 {
		return errors.New("project or service or module not found!")
	}

	err = collection.Update(bson.M{"project": project, "service": service, "module":module}, bson.M{"$pull":bson.M{"instances":bson.M{"$in":instances}}})

	if err != nil {
		db.session.Refresh()
		return err
	}

	return nil
}

//******************User Setting*****************//

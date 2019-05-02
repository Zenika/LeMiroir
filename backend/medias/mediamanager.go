package medias

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/Zenika/MARCEL/backend/clients"
	"github.com/Zenika/MARCEL/backend/commons"
	"github.com/Zenika/MARCEL/backend/containers"
	"github.com/Zenika/MARCEL/backend/plugins"
)

type Manager struct {
	configPath     string
	configFileName string
	configFullpath string
	Config         *Configuration

	pluginManager  *plugins.Manager
	clientsService *clients.Service
}

func NewManager(pluginManager *plugins.Manager, clientsService *clients.Service, configPath, configFilename string) *Manager {
	manager := new(Manager)

	manager.configPath = configPath
	manager.configFileName = configFilename

	manager.pluginManager = pluginManager
	manager.clientsService = clientsService

	manager.configFullpath = filepath.Join(configPath, configFilename)
	manager.Config = NewConfiguration()

	return manager
}

// LoadMedias loads medias configuration from DB and stor it in memory
func (m *Manager) LoadFromDB() {
	log.Debugln("Start Loading Medias from DB.")

	commons.LoadFromDB(m)

	for _, media := range m.Config.Medias {
		if media.IsActive {
			m.Activate(&media)
		}
	}

	log.Debugln("Medias configurations is loaded...")
}

func (m *Manager) GetConfiguration() *Configuration {
	log.Debugln("Getting global medias config")

	return m.Config
}

func (m *Manager) GetConfig() interface{} {
	return m.Config
}

func (m *Manager) GetAll() []Media {
	log.Debugln("Getting all medias")

	return m.Config.Medias
}

// GetMedia Return the media with this id
func (m *Manager) Get(idMedia int) (*Media, error) {

	log.Debugln("Getting media with id: ", idMedia)
	for _, media := range m.Config.Medias {
		if idMedia == media.ID {
			return &media, nil
		}
	}

	return nil, errors.New("No Media found with ID " + strconv.Itoa(idMedia))
}

// CreateMedia Create a new Media, save it into memory and commit
func (m *Manager) CreateEmpty(owner string) *Media {

	log.Debugln("Creating media")

	newMedia := NewMedia()
	newMedia.ID = m.GetNextID()
	newMedia.Name = "Media " + strconv.Itoa(newMedia.ID)
	newMedia.Owner = owner

	//save it into the MediasConfiguration
	m.SaveIntoDB(newMedia)
	m.Commit()

	return newMedia
}

// RemoveMedia RemoveFromDB media from memory
func (m *Manager) RemoveFromDB(media *Media) {
	log.Debugln("Removing media")
	i := m.getPosition(media)

	if i >= 0 {
		m.Config.Medias = append(m.Config.Medias[:i], m.Config.Medias[i+1:]...)
	}
}

// SaveIntoDB saves media information in memory.
func (m *Manager) SaveIntoDB(media *Media) {
	log.Debugln("Saving media")
	m.RemoveFromDB(media)
	m.Config.Medias = append(m.Config.Medias, *media)
}

// Commit SaveIntoDB all medias in DB.
// Here DB is a JSON file
func (m *Manager) Commit() error {
	return commons.Commit(m)
}

// CreateSaveFileIfNotExist check if the save file for medias exists and create it if not.
func (m *Manager) CreateSaveFileIfNotExist(filePath string, fileName string) {

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Infoln("Data directory did not exist. Create it.")
		os.Mkdir(filePath, 0755)
	}

	var fullPath string = filepath.Join(filePath, fileName)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {

		f, err := os.Create(fullPath)
		commons.Check(err)

		log.Infoln("Medias configuration file created at %v", fullPath)

		f.Close()

		//commit a first time to ensure the configuration has been saved
		m.Commit()
	}
}

func (m *Manager) Activate(media *Media) error {
	errorMessages := ""

	for _, mp := range media.Plugins {

		plugin, err := m.pluginManager.Get(mp.EltName)
		if err != nil {
			//plugin does not exist (anymore ?) in the catalog. Obviously, it should never append.
			log.Errorln(err.Error())
			//Don't return an error now, we need to activate the other plugins
			errorMessages += err.Error() + "\n"
		}

		// duplicate plugin files into "medias/{idMedia}/{plugins_EltName}/{idInstance}"
		mpPath := m.GetPluginDirectory(media, mp.EltName, mp.InstanceId)
		err = m.copyNewInstanceOfPlugin(media, &mp, mpPath)
		if err != nil {
			log.Errorln(err.Error())
			//Don't return an error now, we need to activate the other plugins
			errorMessages += err.Error() + "\n"
		}

		if mp.BackEnd != nil {
			retour, err := containers.InstallImage(filepath.Join(mpPath, "back", plugin.Backend.Dockerimage))
			if err != nil {
				//Don't return an error now, we need to activate the other plugins
				log.Errorln(err.Error())
				errorMessages += err.Error() + "\n"
			}

			imageName := strings.TrimSpace(strings.TrimPrefix(retour, "Loaded image: "))
			externalPort := m.GetPortNumberForPlugin()

			dockerContainerId, err := containers.StartContainer(imageName, plugin.Backend.Port, externalPort, mp.BackEnd.Props, mpPath)
			if err != nil {
				//Don't return an error now, we need to activate the other plugins
				log.Errorln(err.Error())
				errorMessages += err.Error() + "\n"
			} else {
				mp.BackEnd.Port = externalPort
				mp.BackEnd.DockerImageName = imageName
				mp.BackEnd.DockerContainerId = strings.TrimSpace(dockerContainerId)
			}
		}
	}

	media.IsActive = true

	m.SaveIntoDB(media)

	if errorMessages != "" {
		return errors.New(errorMessages)
	}

	return nil
}

func (m *Manager) Deactivate(media *Media) error {

	errorMessages := ""
	//stop all backends instances and free ports number
	for _, mp := range media.Plugins {
		if mp.BackEnd != nil {

			err := containers.StopContainer(mp.BackEnd.DockerContainerId)
			if err != nil {
				errorMessages += err.Error() + "\n"
			} else {
				m.FreePortNumberForPlugin(mp.BackEnd.Port)
			}
		}
	}

	media.IsActive = false

	m.SaveIntoDB(media)

	if errorMessages != "" {
		return errors.New(errorMessages)
	}

	return nil
}

func (m *Manager) Delete(media *Media) error {

	m.Deactivate(media)

	m.RemoveFromDB(media)
	m.Commit()

	//remove plugins files
	err := os.RemoveAll(filepath.Join("medias", strconv.Itoa(media.ID)))
	if err != nil {
		return err
	}

	return nil
}

// GetMediaPosition Return position of a media in the list
func (m *Manager) getPosition(media *Media) int {
	for p, m := range m.Config.Medias {
		if m.ID == media.ID {
			return p
		}
	}
	return -1
}

func (m *Manager) GetPortNumberForPlugin() int {

	// 1 : try to pop a port number from the pool
	if len(m.Config.PortsPool) > 0 {
		p := m.Config.PortsPool[0]
		//remove the first number
		m.Config.PortsPool = m.Config.PortsPool[1:]
		return p
	}

	// 2 : if pool is empty, just increment the counter
	p := m.Config.NextFreePortNumber
	m.Config.NextFreePortNumber += 1

	return p
}

func (m *Manager) FreePortNumberForPlugin(portNumber int) {
	for i, v := range m.Config.PortsPool {
		if portNumber == v {
			m.Config.PortsPool = append(m.Config.PortsPool[:i], m.Config.PortsPool[i+1:]...)
		}
	}
}

func (m *Manager) GetNextID() int {
	m.Config.LastID = m.Config.LastID + 1
	return m.Config.LastID
}

func (m *Manager) copyNewInstanceOfPlugin(media *Media, mp *MediaPlugin, path string) error {
	//Copy onlyd frontend and backend dirs since there the only relevant files
	err := commons.CopyDir(filepath.Join("plugins", mp.EltName, "frontend"), filepath.Join(path, "frontend"))
	if _, err := os.Stat(filepath.Join("plugins", mp.EltName, "backend")); !os.IsNotExist(err) {
		err = commons.CopyDir(filepath.Join("plugins", mp.EltName, "backend"), filepath.Join(path, "backend"))
	}

	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) GetPluginDirectory(media *Media, eltName string, instanceId string) string {
	return filepath.Join("medias", strconv.Itoa(media.ID), eltName, instanceId)
}

func (m *Manager) GetSaveFilePath() (string, string, string) {
	return m.configFullpath, m.configPath, m.configFileName
}

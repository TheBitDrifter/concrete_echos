package persistence

import (
	"encoding/json"

	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"
)

type (
	SceneName      string
	EntityEnum     int
	PersistenceID  int
	OptionalSaveID int // If multiple save locations in scene use ID to differ which was used
)

var State = state{
	Scenes:           make(map[SceneName]PersistedEntityTypes),
	FirstPlay:        false,
	FtCheckpointsMap: map[FastTravelCheckpointName]FastTravelCheckpoint{},
	FtCheckpoints:    []FastTravelCheckpoint{},
}

type state struct {
	FirstPlay        bool
	FtCheckpointsMap map[FastTravelCheckpointName]FastTravelCheckpoint
	FtCheckpoints    []FastTravelCheckpoint

	PlayerSingleton    *warehouse.SerializedEntity
	LastScene          SceneName
	LastOptionalSaveID int

	Scenes map[SceneName]PersistedEntityTypes
}

type jsonState struct {
	ScenePersistedEntities ScenePersistedEntities                            `json:"scenePersistedEntities"`
	PlayerSingleton        *warehouse.SerializedEntity                       `json:"playerSingleton"`
	LastScene              SceneName                                         `json:"lastScene"`
	LastOptionalSaveID     int                                               `json:"lastOptionalSaveID"`
	FtCheckpointsMap       map[FastTravelCheckpointName]FastTravelCheckpoint `json:"ftCheckpointsMap"`
	FtCheckpoints          []FastTravelCheckpoint                            `json:"ftCheckpoints"`
}

type ScenePersistedEntities map[SceneName]PersistedEntityTypes

func (ps ScenePersistedEntities) MarshalJSON() ([]byte, error) {
	preparedData, err := warehouse.PrepareForJSONMarshal(map[SceneName]PersistedEntityTypes(ps))
	if err != nil {
		return nil, err
	}
	return json.Marshal(preparedData)
}

type PersistedEntityTypes struct {
	Entities map[EntityEnum]Entities `json:"persistedEntityTypes"`
}
type Entities struct {
	Items map[PersistenceID]warehouse.SerializedEntity `json:"entities"`
}

type FastTravelCheckpointName string

type FastTravelCheckpoint struct {
	FastTravelCheckpointName FastTravelCheckpointName
	SceneName                SceneName
	DropOff                  spatial.Position
}

func (s state) Get(sceneName string, entityType EntityEnum, id PersistenceID) (warehouse.SerializedEntity, bool) {
	sceneState, ok := s.Scenes[SceneName(sceneName)]
	if !ok {
		return warehouse.SerializedEntity{}, false
	}

	entityTypes, ok := sceneState.Entities[entityType]
	if !ok {
		return warehouse.SerializedEntity{}, false
	}

	en, ok := entityTypes.Items[id]
	return en, ok
}

func (s *state) Set(sceneName string, entityType EntityEnum, id PersistenceID, se warehouse.SerializedEntity) {
	sn := SceneName(sceneName)

	sceneState, ok := s.Scenes[sn]
	if !ok {
		sceneState = PersistedEntityTypes{
			Entities: make(map[EntityEnum]Entities),
		}
	}

	etCollection, ok := sceneState.Entities[entityType]
	if !ok {
		etCollection = Entities{
			Items: make(map[PersistenceID]warehouse.SerializedEntity),
		}
	}

	etCollection.Items[id] = se

	if sceneState.Entities == nil {
		sceneState.Entities = map[EntityEnum]Entities{}
	}

	sceneState.Entities[entityType] = etCollection
	s.Scenes[sn] = sceneState
}

func (s state) MarshalJSON() ([]byte, error) {
	js := jsonState{
		ScenePersistedEntities: s.Scenes,
		PlayerSingleton:        s.PlayerSingleton,
		LastScene:              s.LastScene,
		LastOptionalSaveID:     s.LastOptionalSaveID,
		FtCheckpointsMap:       s.FtCheckpointsMap,
		FtCheckpoints:          s.FtCheckpoints,
	}

	return json.Marshal(js)
}

func (s *state) UnmarshalJSON(data []byte) error {
	var js jsonState
	if err := json.Unmarshal(data, &js); err != nil {
		return err
	}

	s.Scenes = js.ScenePersistedEntities
	s.PlayerSingleton = js.PlayerSingleton
	s.LastScene = js.LastScene
	s.LastOptionalSaveID = js.LastOptionalSaveID
	s.FtCheckpointsMap = js.FtCheckpointsMap
	s.FtCheckpoints = js.FtCheckpoints

	if s.Scenes == nil {
		s.Scenes = make(map[SceneName]PersistedEntityTypes)
	}

	if s.FtCheckpointsMap == nil {
		s.FtCheckpointsMap = map[FastTravelCheckpointName]FastTravelCheckpoint{}
	}

	return nil
}

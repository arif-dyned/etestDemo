package provider

import (
	"github.com/DynEd/etest/app/driver"
	"time"

	questionHttp "github.com/DynEd/etest/question/delivery/http"
	questionRepository "github.com/DynEd/etest/question/repository"
	questionUsecase "github.com/DynEd/etest/question/usecase"

	optionHttp "github.com/DynEd/etest/option/delivery/http"
	optionRepository "github.com/DynEd/etest/option/repository"
	optionUsecase "github.com/DynEd/etest/option/usecase"

	levelHttp "github.com/DynEd/etest/level/delivery/http"
	levelRepository "github.com/DynEd/etest/level/repository"
	levelUsecase "github.com/DynEd/etest/level/usecase"

	typesHttp "github.com/DynEd/etest/types/delivery/http"
	typesRepository "github.com/DynEd/etest/types/repository"
	typesUsecase "github.com/DynEd/etest/types/usecase"

	tagsHttp "github.com/DynEd/etest/tags/delivery/http"
	tagsRepository "github.com/DynEd/etest/tags/repository"
	tagsUsecase "github.com/DynEd/etest/tags/usecase"

	groupHttp "github.com/DynEd/etest/group/delivery/http"
	groupRepository "github.com/DynEd/etest/group/repository"
	groupUsecase "github.com/DynEd/etest/group/usecase"

	questionSlotHttp "github.com/DynEd/etest/question_slot/delivery/http"
	questionSlotRepository "github.com/DynEd/etest/question_slot/repository"
	questionSlotUsecase "github.com/DynEd/etest/question_slot/usecase"

	"github.com/gorilla/mux"
)

// RegisterQuestionService register and initialize question services
// The service includes delivery, repository and usecase
func RegisterProviderService(r *mux.Router) {
	timeout := time.Second * 5

	//connection to Database
	connection, _ := driver.ConnectToMysql()

	//Question
	questionHttp.NewQuestionHandler(r, questionUsecase.NewQuestionUsecase(questionRepository.NewMysqlQuestionRepository(connection), timeout))

	//Option
	optionHttp.NewOptionHandler(r, optionUsecase.NewOptionUsecase(optionRepository.NewMysqlOptionRepository(connection), timeout))

	//Level
	levelHttp.NewLevelHandler(r, levelUsecase.NewLevelUsecase(levelRepository.NewMysqlLevelRepository(connection), timeout))
	
	//Types
	typesHttp.NewTypesHandler(r, typesUsecase.NewTypesUsecase(typesRepository.NewMysqlTypesRepository(connection), timeout))
	
	//Tags
	tagsHttp.NewTagsHandler(r, tagsUsecase.NewTagsUsecase(tagsRepository.NewMysqlTagsRepository(connection), timeout))
	
	//Groups
	groupHttp.NewGroupsHandler(r, groupUsecase.NewGroupsUsecase(groupRepository.NewMysqlGroupsRepository(connection), timeout))
	
	//Question Slots
	questionSlotHttp.NewQuestionSlotHandler(r, questionSlotUsecase.NewQuestionSlotUsecase(questionSlotRepository.NewMysqlQuestionSlotRepository(connection), timeout))
}
